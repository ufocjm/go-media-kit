package qywx

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/heibizi/go-media-kit/core/utils/netx"
	"github.com/redis/go-redis/v9"
	"io"
	"net/http"
	"net/url"
	"time"
)

type (
	Client struct {
		baseUrl string
		Config
	}

	Config struct {
		Ctx         context.Context
		CorpId      string
		CorpSecret  string
		AgentId     string
		Proxy       string
		UserId      string
		RedisClient redis.UniversalClient
	}

	resp struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}

	accessTokenResp struct {
		ErrCode     int    `json:"errcode"`
		ErrMsg      string `json:"errmsg"`
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}

	messageReq struct {
		ToUser               string      `json:"touser"`
		MsgType              string      `json:"msgtype"`
		AgentId              string      `json:"agentid"`
		Text                 textMessage `json:"text"`
		Safe                 int         `json:"safe"`
		EnableIdTrans        int         `json:"enable_id_trans"`
		EnableDuplicateCheck int         `json:"enable_duplicate_check"`
	}

	textMessage struct {
		Content string `json:"content"`
	}

	imageMessage struct {
		ToUser  string `json:"touser"`
		MsgType string `json:"msgtype"`
		AgentId string `json:"agentid"`
		News    news   `json:"news"`
	}

	news struct {
		Articles []article `json:"articles"`
	}

	article struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		PicUrl      string `json:"picurl"`
		Url         string `json:"url"`
	}
)

func NewClient(config Config) *Client {
	client := &Client{
		Config: config,
	}
	if config.Proxy != "" {
		client.baseUrl = config.Proxy
	} else {
		client.baseUrl = "https://qyapi.weixin.qq.com"
	}
	return client
}

func (c *Client) send(data any, retry bool) error {
	u, err := url.JoinPath(c.baseUrl, "/cgi-bin/message/send")
	if err != nil {
		return fmt.Errorf("拼接 url 异常: %v", err)
	}
	token, err := c.getAccessToken(false)
	if err != nil {
		return fmt.Errorf("获取 accessToken 异常: %v", err)
	}
	res, err := netx.NewHttpx(netx.HttpRequestParams{
		Ctx:    c.Ctx,
		Method: http.MethodPost,
		Url:    u,
		Params: url.Values{"access_token": {token}},
		Body:   data,
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
			if sendResp.ErrCode == 42001 {
				// 处理 AccessToken 过期的情况
				_, err := c.getAccessToken(true)
				if err != nil {
					return err
				}
				// 重试一次
				if retry {
					return c.send(data, false)
				} else {
					return fmt.Errorf("accessToken 已过期: %v", sendResp.ErrMsg)
				}
			}
			return fmt.Errorf("发送消息异常: %v", sendResp.ErrMsg)
		}
	} else {
		return fmt.Errorf("错误码：%d，错误原因：%s", res.StatusCode, res.Status)
	}
}

func (c *Client) getAccessToken(force bool) (string, error) {
	k := fmt.Sprintf("qywx:accessToken:%s:%s", c.CorpId, c.AgentId)
	if force {
		c.RedisClient.Del(c.Ctx, k)
	} else {
		accessToken := c.RedisClient.Get(context.Background(), k).Val()
		if accessToken != "" {
			return accessToken, nil
		}
	}
	u, err := url.JoinPath(c.baseUrl, "/cgi-bin/gettoken")
	if err != nil {
		return "", fmt.Errorf("拼接 url 异常: %v", err)
	}
	resp, err := netx.NewHttpx(netx.HttpRequestParams{
		Ctx:    c.Ctx,
		Method: http.MethodGet,
		Url:    u,
		Params: url.Values{"corpid": {c.CorpId}, "corpsecret": {c.CorpSecret}},
	}).Request()
	if err != nil {
		return "", fmt.Errorf("获取 accessToken 异常: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("获取 accessToken 异常: %v", resp.Status)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取 accessToken 响应异常: %v", err)
	}
	var atResp accessTokenResp
	err = json.Unmarshal(data, &atResp)
	if err != nil {
		return "", fmt.Errorf("解析 accessToken 响应异常: %v", err)
	}
	if atResp.ErrCode != 0 {
		return "", fmt.Errorf("获取 accessToken 异常，错误码: %d，错误信息: %s", atResp.ErrCode, atResp.ErrMsg)
	}
	c.RedisClient.Set(context.Background(), k, atResp.AccessToken, time.Duration(atResp.ExpiresIn)*time.Second)
	return atResp.AccessToken, nil
}

func (c *Client) userId() string {
	if c.UserId != "" {
		return c.UserId
	}
	return "@all"
}
