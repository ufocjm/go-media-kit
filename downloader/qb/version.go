package qb

import (
	"fmt"
	"github.com/heibizi/go-media-kit/core/utils/netx"
)

func (c *Client) Version() (version string, err error) {
	resp, err := netx.NewHttpx(netx.HttpRequestConfig{
		Ctx:     c.config.Ctx,
		Url:     c.config.Host + "/api/v2/app/version",
		Referer: c.config.Host,
		Cookie:  c.ck,
	}).Get()
	if err != nil || resp.StatusCode != 200 {
		return "", fmt.Errorf("获取版本失败: %v", err)
	}
	body, err := netx.GetBody(resp)
	if err != nil {
		return "", fmt.Errorf("获取版本失败: %v", err)
	}
	return string(body), nil
}

func NewClient(config Config) *Client {
	return &Client{config: config}
}
