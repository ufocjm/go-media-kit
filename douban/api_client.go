package douban

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"github.com/heibizi/go-media-kit/core/utils/mapstructurex"
	"github.com/heibizi/go-media-kit/core/utils/netx"
	"github.com/tidwall/gjson"
	"net/url"
	"strings"
	"time"
)

const apiKey = "054022eaeae0b00e0fc068c0c0a2102a"
const apiUrl = "https://frodo.douban.com/api/v2"
const secretKey = "bf7dddc7c9cfe6f7"
const ua = "MicroMessenger/"
const referer = "https://servicewechat.com/wx2f9b06c1de1ccfca/91/page-frame.html"

type ApiClient struct {
}

func NewApiClient() *ApiClient {
	return &ApiClient{}
}

func (c *ApiClient) sign(urlStr string, ts string, method string) string {
	parsedURL, _ := url.Parse(urlStr)
	urlPath := parsedURL.Path
	rawSign := strings.Join([]string{strings.ToUpper(method), url.QueryEscape(urlPath), ts}, "&")
	h := hmac.New(sha1.New, []byte(secretKey))
	h.Write([]byte(rawSign))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return signature
}

func (c *ApiClient) request(urlStr string, method string, params map[string]string) ([]byte, error) {
	ts := time.Now().Format("20060102")
	p := url.Values{
		"apiKey": {apiKey},
		"os_rom": {"android"},
		"_sig":   {c.sign(urlStr, ts, method)},
		"_ts":    {ts},
	}
	if params != nil {
		for k, v := range params {
			p.Add(k, v)
		}
	}
	resp, err := netx.NewHttpx(netx.HttpRequestParams{
		Ctx:     nil,
		Url:     urlStr,
		Params:  p,
		Cookie:  "",
		Ua:      ua,
		Referer: referer,
	}).Request()
	if err != nil {
		return nil, err
	}
	defer netx.Close(resp)
	body, err := netx.GetBody(resp)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (c *ApiClient) result(data []byte, o interface{}) error {
	result := gjson.Parse(string(data))
	err := mapstructurex.WeakDecode(result.Value(), &o)
	if err != nil {
		return err
	}
	return nil
}

func (c *ApiClient) get(url string, params map[string]string, o interface{}) error {
	data, err := c.request(url, "GET", params)
	if err != nil {
		return err
	}
	return c.result(data, o)
}
