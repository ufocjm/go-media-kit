package qb

import (
	"fmt"
	"github.com/heibizi/go-media-kit/core/utils/netx"
	"net/http"
	"net/url"
)

func (c *Client) Login() error {
	resp, err := netx.NewHttpx(netx.HttpRequestParams{
		Ctx:      c.config.Ctx,
		Method:   http.MethodPost,
		Url:      c.config.Host + "/api/v2/auth/login",
		Referer:  c.config.Host,
		FormData: url.Values{"username": {c.config.Username}, "password": {c.config.Password}},
	}).Request()
	if err != nil {
		return fmt.Errorf("登录失败: %v", err)
	}
	if resp.StatusCode == 200 {
		c.ck = resp.Header.Get("Set-Cookie")
		return nil
	}
	return fmt.Errorf("登录失败: %s", resp.Status)
}
