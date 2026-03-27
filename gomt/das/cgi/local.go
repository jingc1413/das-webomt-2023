package cgi

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/cgi"
	"net/url"
	"path/filepath"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type LocalCGIRequestHandler struct {
	sync.Mutex
	dir string
	// defs model.ParameterDefines
}

func NewLocalCGIBaseHandler(dir string) CGIRequestHandler {
	h := LocalCGIRequestHandler{
		dir: dir,
	}
	// h.defs = model.GetParameterDefinesByDeviceTypeName(deviceTypeName)
	return &h
}

// func (h *LocalCGIRequestHandler) ServePostJsonRequest(script string, query url.Values, data any) (*CGIResponse, error) {
// 	if data == nil {
// 		return h.ServeRequest("POST", script, query, nil, "application/json;charset=UTF-8")
// 	}
// 	w := bytes.NewBuffer([]byte{})
// 	encoder := json.NewEncoder(w)
// 	encoder.SetEscapeHTML(false)
// 	if err := encoder.Encode(data); err != nil {
// 		return nil, errors.Wrap(err, "marshal data")
// 	}
// 	jsonData := w.Bytes()
// 	body := bytes.NewReader(jsonData)
// 	return h.ServeRequest("POST", script, query, body, "application/json;charset=UTF-8")
// }

func (h *LocalCGIRequestHandler) ServeRequest(
	method string,
	script string,
	query url.Values,
	body io.Reader,
	contentType string,
	timeout time.Duration,
) (*CGIResponse, error) {
	path := fmt.Sprintf("/cgi-bin/%v", script)
	if len(query) > 0 {
		if queryString := query.Encode(); queryString != "" {
			path += "?" + queryString
		}
	}
	u, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	handler := new(cgi.Handler)
	handler.Path = filepath.FromSlash(req.URL.Path[len("/cgi-bin/"):])
	handler.Dir = h.dir

	h.Lock()
	defer h.Unlock()

	rw := newResponseWriter()
	handler.ServeHTTP(rw, req)
	if rw.status != 200 {
		return nil, errors.New("response with error status")
	}

	resp := CGIResponse{
		Status: rw.status,
		Body:   io.NopCloser(rw.body),
	}
	return &resp, nil
}

func newResponseWriter() *responseWriter {
	return &responseWriter{
		header: make(http.Header),
		status: http.StatusOK,
		body:   &bytes.Buffer{},
	}
}

type responseWriter struct {
	header http.Header
	status int
	body   *bytes.Buffer
}

func (w *responseWriter) Header() http.Header {
	return w.header
}

func (w *responseWriter) Write(b []byte) (int, error) {
	return w.body.Write(b)
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
}
