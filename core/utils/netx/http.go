package netx

import (
	"context"
	"fmt"
	"github.com/heibizi/go-media-kit/core/utils/mapstructurex"
	"golang.org/x/net/proxy"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
)

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

// JoinURL 将基准 Url 和相对 Url 拼接在一起，返回完整的 Url
func JoinURL(baseURL, relativeURL string) (string, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("error parsing base Url: %v", err)
	}

	rel, err := url.Parse(relativeURL)
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
