package agent

import (
	"gomt/das/cgi"
	"io"
	"path"
	"strings"

	"github.com/pkg/errors"
)

// const (
// 	UPGRADE_MODE_ARM_NORMAL = "arm"
// 	UPGRADE_MODE_ARM_FORCE  = "arm:force"
// 	UPGRADE_MODE_CGI_NORMAL = "cgi"
// 	UPGRADE_MODE_CGI_FORCE  = "cgi:force"
// )

// const (
// 	UPGRADE_STATE_NONE      = 0
// 	UPGRADE_STATE_UPGRADING = 1
// 	UPGRADE_STATE_FINISH    = 2
// )

// type UpgradeLog struct {
// 	Status string
// 	Text   string
// }

// type UpgradeStatusData struct {
// 	Filename      string
// 	Mode          string
// 	Force         bool
// 	State         int
// 	Progress      int
// 	ResultStatus  string
// 	ResultMessage string

// 	TotalNum   int
// 	TimeoutNum int
// 	FailedNum  int
// 	SuccessNum int
// 	Logs       []*UpgradeLog

// 	CgiResp *cgi.UpgradeResponseData
// }

// type UpgradeStatus struct {
// 	sync.Mutex
// 	data *UpgradeStatusData
// }

// func (m *UpgradeStatus) GetData() UpgradeStatusData {
// 	m.Lock()
// 	defer m.Unlock()
// 	data := m.data
// 	return *data
// }

// func (m *UpgradeStatus) IsUpgrading() bool {
// 	m.Lock()
// 	defer m.Unlock()

// 	return m.data.State == UPGRADE_STATE_UPGRADING
// }

// func (m *UpgradeStatus) SetUpgrading(filename string, mode string, force bool) bool {
// 	m.Lock()
// 	defer m.Unlock()

// 	if m.data.State == UPGRADE_STATE_UPGRADING {
// 		return false
// 	}
// 	m.data.State = UPGRADE_STATE_UPGRADING
// 	m.data.Filename = filename
// 	m.data.Mode = mode
// 	m.data.Force = force

// 	m.data.Progress = 0
// 	m.data.TotalNum = -1
// 	m.data.FailedNum = -1
// 	m.data.SuccessNum = -1
// 	m.data.TimeoutNum = -1
// 	m.data.ResultStatus = ""
// 	m.data.ResultMessage = ""
// 	m.data.Logs = []*UpgradeLog{}
// 	m.data.CgiResp = nil

// 	return true
// }

// func (m *UpgradeStatus) SetProgress(progress int, log *UpgradeLog) {
// 	m.Lock()
// 	defer m.Unlock()

// 	m.data.Progress = progress
// 	if log != nil {
// 		m.data.Logs = append(m.data.Logs, log)
// 	}
// }

// func (m *UpgradeStatus) SetResult(status string, msg string) {
// 	m.Lock()
// 	defer m.Unlock()

// 	m.data.ResultStatus = status
// 	m.data.ResultMessage = msg
// 	if m.data.ResultStatus == "success" {
// 		m.data.Progress = 100
// 	}
// 	m.data.State = UPGRADE_STATE_FINISH
// }

func (s *DasDeviceAgent) StartUpgrade(filename string, force bool, byArm bool) (*cgi.UpgradeResponseData, error) {
	// byArm := false
	// if s.info != nil && s.info.ProductType == "AU" {
	// 	byArm = true
	// }

	ftype := s.fileTypes.Get("UpgradeFile")
	if ftype == nil {
		return nil, errors.New("invalid file type for upgrade")
	}
	filename2 := path.Join(ftype.Dir, filename)
	// if utils.ExistsFile(filename2) == false {
	// 	return nil, errors.Errorf("cant find upgrade file %v", filename)
	// }

	resp, err := s.ServeCgiStartUpgrade(filename2, force, byArm)
	if err != nil {
		return nil, errors.Wrap(err, "start upgrade")
	}
	return resp, nil
}

type FirmwareInfo struct {
	Name string
	CRC  string
}

func (s *DasDeviceAgent) GetFirmwareInfos() ([]FirmwareInfo, error) {
	_, f, err := s.GetFile("Config", "FilenameCrcInfo.txt")
	if err != nil {
		return nil, errors.Wrap(err, "get info file")
	}
	defer f.Close()

	out := []FirmwareInfo{}

	buf, err := io.ReadAll(f)
	if err != nil {
		return nil, errors.Wrap(err, "read info file")
	}
	parts := strings.Split(string(buf), ",")
	for _, part := range parts {
		if args := strings.Split(part, "/"); len(args) == 2 {
			out = append(out, FirmwareInfo{
				Name: args[0],
				CRC:  args[1],
			})
		}
	}
	return out, nil
}

func (s *DasDeviceAgent) DeleteFirmwareInfo(name string) error {
	param, err := s.SetParameterValue("TB4.P0ABC", name)
	if err != nil {
		return errors.Wrap(err, "set parameter value")
	}
	if param.Code != "00" {
		return errors.New("response with error code")
	}
	return nil
}
