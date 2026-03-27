package cgi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gomt/core/iam"
	"gomt/core/model"
	"gomt/core/proto/priv"
	"io"
	"mime/multipart"
	"net/url"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

const (
	TIMEOUT_SHORT    = time.Second * 10
	TIMEOUT_LONG     = time.Second * 30
	TIMEOUT_LONGLONG = time.Second * 60
)

type CGIResponse struct {
	Status int
	Body   io.ReadCloser
}

type CGIRequestHandler interface {
	ServeRequest(method string, script string, query url.Values, body io.Reader, contentType string, timeout time.Duration) (*CGIResponse, error)
}

func NewRemoteCGIHandler(schema string, deviceTypeName string, version string, defs model.ParameterDefines, base string) *CGIHandler {
	s := &CGIHandler{
		deviceTypeName: deviceTypeName,
		defs:           defs,
	}
	s.log = logrus.WithFields(logrus.Fields{"cgi": "remote", "deviceType": deviceTypeName, "base": base})
	s.h = NewRemoteCGIBaseHandler(base)
	return s
}

func NewLocalCGIHandler(schema string, deviceTypeName string, version string, defs model.ParameterDefines, dir string) *CGIHandler {
	s := &CGIHandler{
		deviceTypeName: deviceTypeName,
		defs:           defs,
	}
	s.log = logrus.WithFields(logrus.Fields{"cgi": "local", "deviceType": deviceTypeName})
	s.h = NewLocalCGIBaseHandler(dir)
	return s
}

type CGIHandler struct {
	sync.Mutex
	deviceTypeName string
	defs           model.ParameterDefines
	h              CGIRequestHandler
	log            *logrus.Entry
	lastLoginState int
}

func (s *CGIHandler) ServeRequest(
	method string,
	script string,
	query url.Values,
	body io.Reader,
	contentType string,
	timeout time.Duration,
) (*CGIResponse, error) {
	return s.h.ServeRequest(method, script, query, body, contentType, timeout)
}

func (s *CGIHandler) ServeGetRequest(script string, query url.Values, timeout time.Duration) (*CGIResponse, error) {
	return s.ServeRequest("GET", script, query, nil, "", timeout)
}

func (s *CGIHandler) ServePostJsonRequest(script string, query url.Values, data any, timeout time.Duration) (*CGIResponse, error) {
	if data == nil {
		return s.ServeRequest("POST", script, query, nil, "application/json;charset=UTF-8", timeout)
	}
	w := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(data); err != nil {
		return nil, errors.Wrap(err, "marshal data")
	}
	jsonData := w.Bytes()
	body := bytes.NewReader(jsonData)
	return s.ServeRequest("POST", script, query, body, "application/json;charset=UTF-8", timeout)
}

func (s *CGIHandler) ServePostFile(script string, fieldName string, filename string, f io.Reader, timeout time.Duration) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(fieldName, filename)
	if err != nil {
		return err
	}
	if _, err := io.Copy(part, f); err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}
	resp, err := s.ServeRequest("POST", script, nil, body, writer.FormDataContentType(), timeout)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if parts := strings.Split(string(respBody), ","); len(parts) == 2 && len(parts[0]) > 0 {
		code, err := strconv.ParseInt(parts[0], 10, 32)
		if err != nil {
			return err
		}
		switch code {
		case 1:
			return nil
		case 2:
			resp.Body.Close()
			return errors.New("file is invalid")
		default:
			resp.Body.Close()
			return errors.New("file is too large")
		}
	}
	return nil
}

func (s *CGIHandler) ServeLogin(username string, password string) (string, error) {
	data := map[string]any{
		"data": fmt.Sprintf("authenticatelogin=login&username=%v&password=%v&date=%v",
			username, password, time.Now().UnixMilli()),
	}
	resp, err := s.ServePostJsonRequest("index.cgi", nil, data, TIMEOUT_LONG)
	if err != nil {
		return "", errors.Wrap(err, "serve post json request")
	}
	defer resp.Body.Close()

	result, _ := io.ReadAll(resp.Body)
	return string(result), nil
}

func (s *CGIHandler) ServeLogout() error {
	q := url.Values{}
	q.Add("checklogin", "logout")
	resp, err := s.ServeGetRequest("index.cgi", q, TIMEOUT_SHORT)
	if err != nil {
		return errors.Wrap(err, "serve get request")
	}
	defer resp.Body.Close()
	return nil
}

func (s *CGIHandler) ServeGetUsers() ([]*iam.User, error) {
	s.lockLogin()
	defer s.unlockLogin()

	data := map[string]any{
		"data": map[string]any{
			"data": "getuserinfo=1",
		},
	}
	resp, err := s.ServePostJsonRequest("index.cgi", nil, data, TIMEOUT_SHORT)
	if err != nil {
		return nil, errors.Wrap(err, "serve post json request")
	}
	defer resp.Body.Close()
	respData, _ := io.ReadAll(resp.Body)
	result := parseUserList(respData)
	if len(result) == 0 {
		result = []*iam.User{
			iam.MakeUser("admin", "admin", nil),
		}
	}
	return result, nil
}

func (s *CGIHandler) ServeCreateUser(name string, password string) error {
	s.lockLogin()
	defer s.unlockLogin()

	data := map[string]string{
		"data": fmt.Sprintf("usrManage=add&user=%v&password=%v&permission=2&date=%v",
			name, password, fmt.Sprintf("%v", time.Now().UnixMilli())),
	}
	resp, err := s.ServePostJsonRequest("index.cgi", nil, data, TIMEOUT_SHORT)
	if err != nil {
		return errors.Wrap(err, "serve post json request")
	}
	defer resp.Body.Close()

	respData, _ := io.ReadAll(resp.Body)
	if string(respData) != "success" {
		return errors.Errorf(string(respData))
	}
	return nil
}

func (s *CGIHandler) ServeDeleteUser(username string) error {
	s.lockLogin()
	defer s.unlockLogin()

	data := map[string]any{
		"data": fmt.Sprintf("usrManage=del&user=%v&password=&date=%v", username, fmt.Sprintf("%v", time.Now().UnixMilli())),
	}
	resp, err := s.ServePostJsonRequest("index.cgi", nil, data, TIMEOUT_SHORT)
	if err != nil {
		return errors.Wrap(err, "serve post json request")
	}
	defer resp.Body.Close()

	respData, _ := io.ReadAll(resp.Body)
	if string(respData) != "success" {
		return errors.Errorf(string(respData))
	}
	return nil
}

func (s *CGIHandler) ServeSetUserPassword(name string, password string) error {
	s.lockLogin()
	defer s.unlockLogin()

	data := map[string]any{
		"data": fmt.Sprintf("usrManage=modify&user=%v&password=%v&permission=2&date=%v",
			name, password, fmt.Sprintf("%v", time.Now().UnixMilli())),
	}
	resp, err := s.ServePostJsonRequest("index.cgi", nil, data, TIMEOUT_SHORT)
	if err != nil {
		return errors.Wrap(err, "serve post json request")
	}
	defer resp.Body.Close()

	respData, _ := io.ReadAll(resp.Body)
	if string(respData) != "success" {
		return errors.Errorf(string(respData))
	}
	return nil
}

func (s *CGIHandler) ServeQueryDevices(schema string, query string, cb func(total uint8, index uint8, info *model.DeviceInfo)) ([]*model.DeviceInfo, error) {
	var wg sync.WaitGroup
	out := []*model.DeviceInfo{}
	total, index, info, err := s.doServeCgiQueryDevices(schema, 1, 1, query)
	if err != nil || info == nil {
		return out, errors.Wrap(err, "query one device")
	}
	out = append(out, info)
	if cb != nil {
		wg.Add(1)
		go func(total uint8, index uint8, info *model.DeviceInfo) {
			defer wg.Done()
			cb(total, index, info)
		}(total, index, info)
	}

	for i := index + 1; i <= total; i++ {
		_, _, info, err := s.doServeCgiQueryDevices(schema, total, i, query)
		if err != nil || info == nil {
			s.log.Error(errors.Wrap(err, "query one device"))
			continue
		}
		out = append(out, info)
		if cb != nil {
			wg.Add(1)
			go func(total uint8, index uint8, info *model.DeviceInfo) {
				defer wg.Done()
				cb(total, index, info)
			}(total, i, info)
		}
	}

	return out, nil
}

func (s *CGIHandler) doServeCgiQueryDevices(schema string, total uint8, index uint8, query string) (uint8, uint8, *model.DeviceInfo, error) {
	cmd, caches, err := MarshalQueryDeviceRequestData(s.defs, total, index, query)
	if err != nil || cmd == "" {
		return total, index, nil, errors.Wrap(err, "marshal request data")
	}
	data := map[string]string{
		"data": cmd,
	}
	resp, err := s.ServePostJsonRequest("index.cgi", nil, data, TIMEOUT_SHORT)
	if err != nil {
		return total, index, nil, errors.Wrap(err, "post request")
	}
	respData, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	_total, _index, info, err := UnmarshalQueryDevicesResponseContent(schema, s.defs, caches, respData)
	if err != nil || info == nil {
		return total, index, nil, errors.Wrap(err, "unmarshal response data")
	}
	return _total, _index, info, nil
}

func (s *CGIHandler) ServeGetParameterValues(inputValues []priv.Parameter) ([]priv.Parameter, error) {
	values := []priv.Parameter{}
	s.dumpParamterValues("get parameter request", inputValues)

	cmds, caches, err := MarshalGetParameterValuesRequestData(s.defs, inputValues, true)
	if err != nil {
		return values, errors.Wrap(err, "marshal request data")
	}
	for _, cmd := range cmds {
		data := map[string]string{
			"data": cmd,
		}
		s.log.Tracef("%v", cmd)

		resp, err := s.ServePostJsonRequest("index.cgi", nil, data, TIMEOUT_SHORT)
		if err != nil {
			return values, errors.Wrap(err, "post request")
		}
		respData, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		values2, err := UmarshalIndexResponseContent(s.defs, caches, respData)
		if err != nil {
			return values, errors.Wrapf(err, "unmarshal response data, %v", string(respData))
		}
		values = append(values, values2...)
	}
	s.dumpParamterValues("get parameter response", values)
	return values, nil
}

func (s *CGIHandler) ServeSetParameterValues(inputValues []priv.Parameter) ([]priv.Parameter, error) {
	values := []priv.Parameter{}
	s.dumpParamterValues("set parameter request", inputValues)

	cmds, caches, err := MarshalSetParameterValuesRequestData(s.defs, inputValues, true)
	if err != nil {
		return values, errors.Wrap(err, "marshal request data")
	}
	for _, cmd := range cmds {
		data := map[string]string{
			"data": cmd,
		}
		s.log.Tracef("%v", cmd)

		resp, err := s.ServePostJsonRequest("index.cgi", nil, data, TIMEOUT_SHORT)
		if err != nil {
			return values, errors.Wrap(err, "post request")
		}
		respData, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		values2, err := UmarshalIndexResponseContent(s.defs, caches, respData)
		if err != nil {
			return values, errors.Wrap(err, "unmarshal response data")
		}
		values = append(values, values2...)
	}
	s.dumpParamterValues("set parameter response", values)
	return values, nil
}

func (s *CGIHandler) dumpParamterValues(prefix string, values []priv.Parameter) {
	args := []string{}
	for _, value := range values {
		if value.Code != "" {
			args = append(args, fmt.Sprintf("%v=%v(%v)", value.ID, value.Value, value.Code))
		} else {
			args = append(args, fmt.Sprintf("%v=%v", value.ID, value.Value))
		}
	}
	s.log.Tracef("%v, %v", prefix, strings.Join(args, ", "))
}

func (s *CGIHandler) ServeStartUpgrade(filename string, force bool, byArm bool) (*UpgradeResponseData, error) {
	s.lockLogin()
	defer s.unlockLogin()

	q := url.Values{}
	if byArm {
		q.Add("upgraupdatefilebyarm", filename)
	} else {
		q.Add("upgraupdatefile", filename)
	}
	if force {
		q.Add("sForce", "1")
	}
	resp, err := s.ServeGetRequest("upload.cgi", q, TIMEOUT_LONGLONG)
	if err != nil {
		return nil, errors.Wrap(err, "serve cgi get request")
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read response content")
	}

	upgradeResp, err := parseUpgradeResult(data)
	if err != nil {
		return nil, errors.Wrap(err, "parse response content")
	}
	return &upgradeResp, nil
}

func (s *CGIHandler) ServeSetUpgradeReboot() (bool, error) {
	s.lockLogin()
	defer s.unlockLogin()

	q := url.Values{}
	q.Add("reboot", "1")
	resp, err := s.ServeGetRequest("upload.cgi", q, TIMEOUT_LONGLONG)
	if err != nil {
		return false, errors.Wrap(err, "serve cgi get request")
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, errors.Wrap(err, "read response content")
	}

	return string(data) == "1", nil
}

type FileInfo struct {
	FileName string
	FileSize int64
	ModTime  int64
}

func (s *CGIHandler) ServeListFiles(dir string) ([]FileInfo, error) {
	s.lockLogin()
	defer s.unlockLogin()

	result := []FileInfo{}

	q := url.Values{}
	q.Add("date", fmt.Sprintf("%v", time.Now().UnixMilli()))
	q.Add("refresh", filepath.ToSlash(filepath.Clean(dir)))
	resp, err := s.ServeGetRequest("upload.cgi", q, TIMEOUT_SHORT)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	doc.Find("ul>li[class=li_normal]").Each(func(i int, selection *goquery.Selection) {
		info := FileInfo{}
		selection.Find(".filenamewithicon").Each(func(i int, v *goquery.Selection) {
			if text := v.Text(); text != "" {
				info.FileName = text
			}
		})
		if info.FileName == "" {
			return
		}
		selection.Find(".fileByte").Each(func(i int, v *goquery.Selection) {
			if text := v.Text(); text != "" {
				re := regexp.MustCompile(`\s*(\d+)\s*(B|KB|MB|GB)\s*`)
				if subs := re.FindStringSubmatch(strings.ToUpper(text)); len(subs) == 3 {
					if subs[2] == "GB" {
						info.FileSize = cast.ToInt64(subs[1]) * 1024 * 1024 * 1024
					} else if subs[2] == "MB" {
						info.FileSize = cast.ToInt64(subs[1]) * 1024 * 1024
					} else if subs[2] == "KB" {
						info.FileSize = cast.ToInt64(subs[1]) * 1024
					} else if subs[2] == "B" {
						info.FileSize = cast.ToInt64(subs[1])
					}
				}
			}
		})
		selection.Find(".fileTime").Each(func(i int, v *goquery.Selection) {
			if text := v.Text(); text != "" {
				re := regexp.MustCompile(`-(\d{1})-`)
				re2 := regexp.MustCompile(`-(\d{1}) `)
				text = re.ReplaceAllString(text, "-0$1-")
				text = re2.ReplaceAllString(text, "-0$1 ")
				if t, err := time.ParseInLocation("2006-01-02 15:04", text, time.UTC); err == nil {
					info.ModTime = t.Unix()
				}
			}
		})
		result = append(result, info)
	})

	return result, nil
}

func (s *CGIHandler) ServeGetFile(dir string, filename string) (*bytes.Reader, error) {
	s.lockLogin()
	defer s.unlockLogin()

	q := url.Values{}
	q.Add("download", filepath.ToSlash(filepath.Clean(dir)))
	q.Add("filename", filename)
	q.Add("encodename", filename)
	resp, err := s.ServeGetRequest("download.cgi", q, TIMEOUT_LONG)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
	f := bytes.NewReader(data)
	return f, nil
}

func (s *CGIHandler) ServeSaveFile(dir string, filename string, fieldName string, f io.Reader) error {
	s.lockLogin()
	defer s.unlockLogin()

	q := url.Values{}
	q.Add("date", fmt.Sprintf("%v", time.Now().UnixMilli()))
	q.Add("saveuploadforder", filepath.ToSlash(filepath.Clean(dir)))
	resp, err := s.ServeGetRequest("upload.cgi", q, TIMEOUT_LONG)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := s.ServePostFile("uploadfile.cgi", fieldName, filename, f, TIMEOUT_LONGLONG); err != nil {
		return err
	}

	return nil
}

func (s *CGIHandler) ServeGetUpgradeFilePacketInfo(dir string, filename string) error {
	s.lockLogin()
	defer s.unlockLogin()

	q := url.Values{}
	q.Add("date", fmt.Sprintf("%v", time.Now().UnixMilli()))
	q.Add("viewInfo", filepath.ToSlash(filepath.Clean(filepath.Join(dir, filename))))
	resp, err := s.ServeGetRequest("upload.cgi", q, TIMEOUT_LONG)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	result, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(result), "Succeed") {
		return errors.New(string(result))
	}

	return nil
}

func (s *CGIHandler) ServeRemoveFile(dir string, fname string) error {
	s.lockLogin()
	defer s.unlockLogin()

	filename := filepath.Join(dir, fname)
	filename = filepath.ToSlash(filepath.Clean(filename))
	q := url.Values{}
	q.Add("date", fmt.Sprintf("%v", time.Now().UnixMilli()))
	q.Add("delfile", filename)
	resp, err := s.ServeGetRequest("upload.cgi", q, TIMEOUT_LONG)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (s *CGIHandler) ServeDeleteKeyAndLogs() error {
	s.lockLogin()
	defer s.unlockLogin()

	resp, err := s.ServeGetRequest("upload.cgi?DeleteKeyandLogs", nil, TIMEOUT_LONG)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (s *CGIHandler) serveGetLoginState() (int, error) {
	values, err := s.ServeGetParameterValues([]priv.Parameter{
		{ID: "TB4.P0B31"},
		{ID: "TB4.P0B33"},
		{ID: "TB4.P0B34"},
	})
	if err != nil {
		return 0, errors.Wrap(err, "get parameter values")
	}
	loginState := 0
	for _, value := range values {
		if value.ID == "TB4.P0B31" {
			loginState = cast.ToInt(value.Value)
		}
	}
	return loginState, nil
}

func (s *CGIHandler) serveSetLoginState(state int) error {
	username := ""
	if state != 0 {
		username = "admin"
	}
	values, err := s.ServeSetParameterValues([]priv.Parameter{
		{ID: "TB4.P0B31", Value: fmt.Sprintf("%d", state)},
		{ID: "TB4.P0B33", Value: username},
		{ID: "TB4.P0B34", Value: fmt.Sprintf("%d", state)},
	})
	if err != nil {
		return errors.Wrap(err, "set parameter values")
	}
	for _, value := range values {
		if value.Code != "00" {
			return errors.Wrapf(err, "set login state parameters response with error code %v", value.Code)
		}
	}
	return nil
}

func (s *CGIHandler) lockLogin() {
	s.Lock()
	s.log.Trace("lock login")
	state, err := s.serveGetLoginState()
	if err != nil {
		s.log.Error(errors.Wrap(err, "get login state"))
	}
	s.lastLoginState = state
	if state == 0 {
		if err := s.serveSetLoginState(1); err != nil {
			s.log.Error(errors.Wrap(err, "set login state"))
		}
	}
}

func (s *CGIHandler) unlockLogin() {
	defer s.Unlock()
	s.log.Trace("unlock login")

	// if s.lastLoginState == 0 {
	// 	if err := s.serveSetLoginState(0); err != nil {
	// 		s.log.Error(errors.Wrap(err, "set login state"))
	// 	}
	// 	// s.serveLogout()
	// }
}

// func (s *CGIHandler) serveLogout() error {
// 	q := url.Values{}
// 	q.Add("checklogin", "logout")
// 	resp, err := s.ServeGetRequest("index.cgi", q, TIMEOUT_SHORT)
// 	if err != nil {
// 		return errors.Wrap(err, "serve get request")
// 	}
// 	defer resp.Body.Close()
// 	return nil
// }
