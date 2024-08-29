package netx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Httpx struct {
	config HttpRequestConfig
	resp   *http.Response
}

func NewHttpx(config HttpRequestConfig) *Httpx {
	if config.Timeout == 0 {
		config.Timeout = time.Second * 60
	}
	if config.MaxIdleConn == 0 {
		config.MaxIdleConn = 10
	}
	if config.IdleConnTimeout == 0 {
		config.IdleConnTimeout = time.Second * 60
	}
	return &Httpx{
		config: config,
	}
}

type HttpRequestConfig struct {
	Ctx             context.Context // 上下文
	Url             string          // 请求地址
	Params          url.Values      // 请求参数
	Header          http.Header     // 请求头
	Cookie          string          // Cookie
	Ua              string          // UserAgent
	Referer         string          // Referer
	ContentType     string          // Content-Type
	Proxy           string          // 代理, 如 http://127.0.0.1:7890、socks5://127.0.0.1:7891、sock5h://127.0.0.1:7892
	Timeout         time.Duration   // 超时
	MaxIdleConn     int             // 最大空闲连接
	IdleConnTimeout time.Duration   // 空闲连接超时
}

func (h *Httpx) Get() (*http.Response, error) {
	return h.Request(http.MethodGet, nil, nil)
}

func (h *Httpx) GetAndDecode(output any) (*http.Response, error) {
	resp, err := h.Get()
	if err != nil {
		return resp, err
	}
	return resp, Decode(resp, output)
}

func (h *Httpx) Post(formData url.Values, body any) (*http.Response, error) {
	return h.Request(http.MethodPost, formData, body)
}

func (h *Httpx) PostAndDecode(formData url.Values, body any, output any) (*http.Response, error) {
	resp, err := h.Post(formData, body)
	if err != nil {
		return nil, err
	}
	return resp, Decode(resp, output)
}

func (h *Httpx) newRequest(method string, formData url.Values, body any) (req *http.Request, err error) {
	u, err := url.Parse(h.config.Url)
	if err != nil {
		return nil, fmt.Errorf("url 处理错误: %v", err)
	}
	query := u.Query()
	for k, vs := range h.config.Params {
		for _, v := range vs {
			query.Add(k, v)
		}
	}
	u.RawQuery = query.Encode()
	var bodyReader io.Reader
	if len(formData) > 0 {
		bodyReader = bytes.NewBufferString(formData.Encode())
	} else if body != nil {
		if b, ok := body.(io.Reader); ok {
			bodyReader = b
		} else if b, ok := body.([]byte); ok {
			bodyReader = bytes.NewBuffer(b)
		} else if b, ok := body.(string); ok {
			bodyReader = bytes.NewBufferString(b)
		} else {
			buf := new(bytes.Buffer)
			b, err := json.Marshal(body)
			if err != nil {
				return nil, err
			}
			buf = bytes.NewBuffer(b)
			bodyReader = buf
		}
	}
	if h.config.Ctx != nil {
		req, err = http.NewRequestWithContext(h.config.Ctx, method, u.String(), bodyReader)
	} else {
		req, err = http.NewRequest(method, u.String(), bodyReader)
	}
	if err != nil {
		return nil, fmt.Errorf("创建请求异常: %v", err)
	}
	if h.config.ContentType != "" {
		req.Header.Set("Content-Type", h.config.ContentType)
	} else if len(formData) > 0 {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, vs := range h.config.Header {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}
	if h.config.Ua != "" {
		req.Header.Set("User-Agent", h.config.Ua)
	}
	if h.config.Cookie != "" {
		req.Header.Set("Cookie", h.config.Cookie)
	}
	if h.config.Referer != "" {
		req.Header.Set("Referer", h.config.Referer)
	}
	return req, err
}

func (h *Httpx) Request(method string, formData url.Values, body any) (*http.Response, error) {
	req, err := h.newRequest(method, formData, body)
	if err != nil {
		return nil, err
	}
	tp, err := NewHttpTransport(h.config.Proxy)
	if err != nil {
		return nil, err
	}
	if tp == nil {
		tp = &http.Transport{}
	}
	tp.MaxIdleConns = h.config.MaxIdleConn
	tp.IdleConnTimeout = h.config.IdleConnTimeout
	client := &http.Client{
		Timeout:   h.config.Timeout,
		Transport: tp,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求异常: %v", err)
	}
	return resp, nil
}
