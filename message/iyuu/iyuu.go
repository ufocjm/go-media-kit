package iyuu

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/heibizi/go-media-kit/core/utils/netx"
	"io"
	"net/http"
	"net/url"
)

type (
	Client struct {
		baseUrl string
		Config
	}

	Config struct {
		Ctx   context.Context
		Token string
	}

	resp struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}

	messageReq struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
)

func NewClient(config Config) *Client {
	client := &Client{
		Config: config,
	}
	client.baseUrl = "https://iyuu.cn/%s.send"
	return client
}

func (c *Client) send(req messageReq) error {
	u := fmt.Sprintf(c.baseUrl, c.Token)
	res, err := netx.NewHttpx(netx.HttpRequestParams{
		Ctx:    c.Ctx,
		Method: http.MethodGet,
		Url:    u,
		Params: url.Values{"text": {req.Title}, "desp": {req.Content}},
	}).Request()
	if err != nil {
		return fmt.Errorf("发送消息异常: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusOK {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("读取响应体异常: %v", err)
		}
		var sendResp resp
		err = json.Unmarshal(body, &sendResp)
		if err != nil {
			return fmt.Errorf("解析响应体异常: %v", err)
		}
		if sendResp.ErrCode == 0 {
			return nil
		} else {
			return fmt.Errorf("发送消息异常: %v", sendResp.ErrMsg)
		}
	} else {
		return fmt.Errorf("错误码：%d，错误原因：%s", res.StatusCode, res.Status)
	}
}
