package netx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/heibizi/go-media-kit/core/utils/mapstructurex"
	"golang.org/x/net/proxy"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Httpx struct {
	params HttpRequestParams
}

func NewHttpx(params HttpRequestParams) *Httpx {
	if params.Method == "" {
		params.Method = http.MethodGet
	}
	if params.Timeout == 0 {
		params.Timeout = time.Second * 60
	}
	if params.MaxIdleConn == 0 {
		params.MaxIdleConn = 10
	}
	if params.IdleConnTimeout == 0 {
		params.IdleConnTimeout = time.Second * 60
	}
	return &Httpx{params: params}
}

type HttpRequestParams struct {
	Ctx             context.Context // 上下文
	Method          string          // 请求方法
	Url             string          // 请求地址
	Params          url.Values      // 请求参数
	FormData        url.Values      // 请求表单
	Body            any             // 请求体, io.Reader、[]byte、string、map、struct
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

func (h *Httpx) NewRequest() (req *http.Request, err error) {
	params := h.params
	u, err := url.Parse(params.Url)
	if err != nil {
		return nil, fmt.Errorf("url 处理错误: %v", err)
	}
	query := u.Query()
	for k, vs := range params.Params {
		for _, v := range vs {
			query.Add(k, v)
		}
	}
	u.RawQuery = query.Encode()
	var bodyReader io.Reader
	body := params.Body
	if len(params.FormData) > 0 {
		bodyReader = bytes.NewBufferString(params.FormData.Encode())
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
	if params.Ctx != nil {
		req, err = http.NewRequestWithContext(params.Ctx, params.Method, u.String(), bodyReader)
	} else {
		req, err = http.NewRequest(params.Method, u.String(), bodyReader)
	}
	if err != nil {
		return nil, fmt.Errorf("创建请求异常: %v", err)
	}
	if params.ContentType != "" {
		req.Header.Set("Content-Type", params.ContentType)
	} else if len(params.FormData) > 0 {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, vs := range params.Header {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}
	if params.Ua != "" {
		req.Header.Set("User-Agent", params.Ua)
	}
	if params.Cookie != "" {
		req.Header.Set("Cookie", params.Cookie)
	}
	if params.Referer != "" {
		req.Header.Set("Referer", params.Referer)
	}
	return req, err
}

func (h *Httpx) NewClient() (*http.Client, error) {
	tp, err := NewHttpTransport(h.params.Proxy)
	if err != nil {
		return nil, err
	}
	if tp == nil {
		tp = &http.Transport{}
	}
	tp.MaxIdleConns = h.params.MaxIdleConn
	tp.IdleConnTimeout = h.params.IdleConnTimeout
	client := &http.Client{
		Timeout:   h.params.Timeout,
		Transport: tp,
	}
	return client, nil
}

func (h *Httpx) Request() (*http.Response, error) {
	req, err := h.NewRequest()
	if err != nil {
		return nil, err
	}
	client, err := h.NewClient()
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求异常: %v", err)
	}
	return resp, nil
}

func (h *Httpx) Decode(output any) (*http.Response, error) {
	resp, err := h.Request()
	if err != nil {
		return resp, err
	}
	return resp, Decode(resp, output)
}

func (h *Httpx) Body() ([]byte, *http.Response, error) {
	resp, err := h.Request()
	if err != nil {
		return nil, nil, err
	}
	body, err := GetBody(resp)
	if err != nil {
		return nil, nil, err
	}
	return body, resp, nil
}

func NewHttpTransport(proxyUrl string) (*http.Transport, error) {
	if proxyUrl == "" {
		return nil, nil
	}
	parsedUrl, err := url.Parse(proxyUrl)
	if err != nil {
		return nil, fmt.Errorf("proxy url parse error: %v", err)
	}
	switch parsedUrl.Scheme {
	case "http", "https":
		return &http.Transport{
			Proxy: http.ProxyURL(parsedUrl),
		}, nil
	case "socks5", "socks5h":
		dialer, err := proxy.SOCKS5("tcp", strings.Replace(proxyUrl, parsedUrl.Scheme+"://", "", 1), nil, proxy.Direct)
		if err != nil {
			return nil, fmt.Errorf("proxy %s error: %v", parsedUrl.Scheme, err)
		}
		dialContext := func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.Dial(network, addr)
		}
		return &http.Transport{
			DialContext: dialContext,
		}, nil
	default:
		return nil, fmt.Errorf("proxy url scheme error: %s", parsedUrl.Scheme)
	}
}

func IsValidHttpUrl(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

// JoinUrl 将基准 Url 和相对 Url 拼接在一起，返回完整的 Url
func JoinUrl(baseUrl, relativeUrl string) (string, error) {
	base, err := url.Parse(baseUrl)
	if err != nil {
		return "", fmt.Errorf("error parsing base Url: %v", err)
	}

	rel, err := url.Parse(relativeUrl)
	if err != nil {
		return "", fmt.Errorf("error parsing relative Url: %v", err)
	}

	// 组合 Url
	resolvedURL := base.ResolveReference(rel)
	return resolvedURL.String(), nil
}

func GetBody(resp *http.Response) ([]byte, error) {
	if resp == nil || resp.Body == nil {
		return nil, nil
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体异常: %v", err)
	}
	return body, nil
}

func Decode(resp *http.Response, output any) error {
	body, err := GetBody(resp)
	if err != nil {
		return err
	}
	err = mapstructurex.WeakDecode(body, &output)
	if err != nil {
		return fmt.Errorf("解析响应体异常: %v", err)
	}
	return nil
}

func Close(resp *http.Response) {
	if resp == nil || resp.Body == nil {
		return
	}
	resp.Body.Close()
}
