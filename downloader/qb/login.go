package qb

import (
	"fmt"
	"github.com/heibizi/go-media-kit/core/utils/netx"
	"net/url"
)

func (c *Client) Login() error {
	resp, err := netx.NewHttpx(netx.HttpRequestConfig{
		Ctx:     c.config.Ctx,
		Url:     c.config.Host + "/api/v2/auth/login",
		Referer: c.config.Host,
	}).Post(url.Values{"username": {c.config.Username}, "password": {c.config.Password}}, nil)
	if err != nil {
		return fmt.Errorf("登录失败: %v", err)
	}
	if resp.StatusCode == 200 {
		c.ck = resp.Header.Get("Set-Cookie")
		return nil
	}
	return fmt.Errorf("登录失败: %s", resp.Status)
}
