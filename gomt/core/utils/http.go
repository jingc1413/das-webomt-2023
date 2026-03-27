package utils

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

var defaultTimeout = 20 * time.Second

func HttpRequestWithTimeout(method string, u string, headers map[string]string, body io.Reader, timeout time.Duration) (*http.Response, error) {
	if timeout <= time.Second {
		timeout = time.Second
	}
	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
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
	return client.Do(req)
}

/*
	func HttpRequestWithTimeout(method string, u string, headers map[string]string, body io.Reader, timeout time.Duration) (*http.Response, error) {
		if timeout <= time.Second {
			timeout = time.Second
		}

		// 手动建立 TLS 连接
		conn, err := customTLSConnect("9.7.3.200:443", timeout)
		if err != nil {
			return nil, fmt.Errorf("TLS 连接失败: %v", err)
		}
		defer conn.Close()

		// 创建 HTTP 请求
		req, err := http.NewRequest(method, u, body)
		if err != nil {
			return nil, err
		}

		for k, v := range headers {
			req.Header.Set(k, v)
		}

		// 使用自定义连接发送请求
		transport := &http.Transport{
			DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return conn, nil
			},
			DisableCompression: true,
			ForceAttemptHTTP2:  false,
		}

		client := &http.Client{
			Transport: transport,
			Timeout:   timeout,
		}

		rsp, err := client.Do(req)
		//conn.Close()
		//body1, _ := io.ReadAll(rsp.Body)
		//fmt.Println(string(body1))
		return rsp, err
	}
*/
func HttpRequestWithTimeout2(method string, u string, headers map[string]string, body io.Reader, timeout time.Duration) ([]byte, error) {
	if timeout <= time.Second {
		timeout = time.Second
	}

	// 手动建立 TLS 连接
	conn, err := customTLSConnect("9.7.3.200:443", timeout)
	if err != nil {
		return nil, fmt.Errorf("TLS 连接失败: %v", err)
	}
	defer conn.Close()

	// 创建 HTTP 请求
	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// 使用自定义连接发送请求
	transport := &http.Transport{
		DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return conn, nil
		},
		DisableCompression: true,
		ForceAttemptHTTP2:  false,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}

	rsp, err := client.Do(req)
	//conn.Close()
	body1, err := io.ReadAll(rsp.Body)
	return body1, err
}

func customTLSConnect(addr string, timeout time.Duration) (net.Conn, error) {
	// 先建立 TCP 连接
	dialer := &net.Dialer{Timeout: timeout}
	tcpConn, err := dialer.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	// 配置 TLS
	config := &tls.Config{
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS10,
		MaxVersion:         tls.VersionTLS13,
		CipherSuites: []uint16{
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
		//ServerName: "9.7.3.200", // 设置 SNI
	}

	// 升级为 TLS 连接
	tlsConn := tls.Client(tcpConn, config)

	// 执行 TLS 握手
	if err := tlsConn.Handshake(); err != nil {
		tcpConn.Close()
		return nil, fmt.Errorf("TLS 握手失败: %v", err)
	}

	// 打印连接信息
	state := tlsConn.ConnectionState()
	fmt.Printf("TLS 连接成功 - 版本: %v, 密码套件: %v\n", state.Version, state.CipherSuite)

	return tlsConn, nil
}

// 自定义 close
type readCloser struct {
	io.ReadCloser
	conn net.Conn
}

func (r *readCloser) Close() error {
	r.conn.Close()
	return r.ReadCloser.Close()
}

func HttpGetWithTimeout(u string, timeout time.Duration) (*http.Response, error) {
	return HttpRequestWithTimeout("GET", u, nil, nil, timeout)
}

func HttpGetWithTimeout2(u string, timeout time.Duration) ([]byte, error) {
	return HttpRequestWithTimeout2("GET", u, nil, nil, timeout)
}

func HttpGet(u string) (*http.Response, error) {
	return HttpGetWithTimeout(u, defaultTimeout)
}

func HttpGet2(u string) ([]byte, error) {
	return HttpGetWithTimeout2(u, defaultTimeout)
}

func HttpPostJson(u string, data any) (*http.Response, error) {
	w := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(data); err != nil {
		return nil, errors.Wrap(err, "marshal data")
	}
	jsonData := w.Bytes()
	body := bytes.NewReader(jsonData)
	headers := map[string]string{
		"Content-Type": "application/json;charset=UTF-8",
	}
	return HttpRequestWithTimeout("POST", u, headers, body, defaultTimeout)
}

func HttpPostFile(u string, fieldname string, filename string, f io.Reader) (*http.Response, error) {
	if fieldname == "" {
		fieldname = "file"
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(fieldname, filename)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(part, f); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	headers := map[string]string{
		"Content-Type": writer.FormDataContentType(),
	}
	return HttpRequestWithTimeout("POST", u, headers, body, defaultTimeout)
}
