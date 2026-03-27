package cgi

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type RemoteCGIRequestHandler struct {
	sync.Mutex
	base string
}

func NewRemoteCGIBaseHandler(base string) CGIRequestHandler {
	h := RemoteCGIRequestHandler{
		base: base,
	}
	return &h
}

func (h *RemoteCGIRequestHandler) ServeRequest(
	method string,
	script string,
	query url.Values,
	body io.Reader,
	contentType string,
	timeout time.Duration,
) (*CGIResponse, error) {
	rawURL := fmt.Sprintf("%v/%v", h.base, script)
	if len(query) > 0 {
		if queryString := query.Encode(); queryString != "" {
			rawURL += "?" + queryString
		}
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				CipherSuites: []uint16{
					tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
					tls.TLS_RSA_WITH_AES_128_CBC_SHA,
					tls.TLS_RSA_WITH_AES_128_CBC_SHA256,
					tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_RSA_WITH_AES_256_CBC_SHA,
					tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				},
				MinVersion: tls.VersionTLS12,
				MaxVersion: tls.VersionTLS12,
			},
			DisableCompression: true,
		},
		Timeout: timeout,
	}
	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	h.Lock()
	defer h.Unlock()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return &CGIResponse{
		Status: resp.StatusCode,
		Body:   resp.Body,
	}, nil
}
