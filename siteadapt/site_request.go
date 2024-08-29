package siteadapt

import (
	"encoding/json"
	"fmt"
	"github.com/heibizi/go-media-kit/core/utils/netx"
	"net/http"
	"net/url"
	"strings"
)

type (
	siteRequest struct {
		sc  Config
		rsp RequestSiteParams

		requestUrl string
	}
)

func newSiteRequest(sc Config, rsp RequestSiteParams) *siteRequest {
	return &siteRequest{
		sc:  sc,
		rsp: rsp,
	}
}

func (cr *siteRequest) request(rd RequestDefinition) (*http.Response, error) {
	// 占位符替换
	rd.Params = cr.replaceUrlValues(rd.Params)
	for k, vs := range cr.rsp.Params {
		for _, v := range vs {
			rd.Params.Add(k, v)
		}
	}
	rd.FormData = cr.replaceUrlValues(rd.FormData)
	if cr.rsp.FormData != nil {
		for k, vs := range cr.rsp.FormData {
			for _, v := range vs {
				rd.FormData.Add(k, v)
			}
		}
	}
	requestUrl, err := cr.newRequestUrl(rd)
	if err != nil {
		return nil, cr.newError("创建请求地址失败", err)
	}
	requestUrl = cr.replacePlaceholders(requestUrl)
	header := http.Header{}
	// header 的优先级：自定义请求头 > 站点配置请求头 > 默认请求头
	// 默认请求头
	if cr.rsp.UA != "" {
		header.Set("User-Agent", cr.rsp.UA)
	}
	if len(cr.rsp.Cookie) > 0 {
		header.Set("Cookie", cr.rsp.Cookie)
	}
	// 站点配置请求头
	if rd.Headers != nil {
		for k, v := range rd.Headers {
			v = cr.replacePlaceholders(v)
			header.Set(k, v)
		}
	}
	// 自定义请求头
	for k, v := range cr.rsp.Headers {
		header.Set(k, v)
	}
	// 校验必填请求头
	if len(rd.RequiredHeaders) > 0 {
		for _, requiredHeader := range rd.RequiredHeaders {
			if len(header.Get(requiredHeader)) == 0 {
				return nil, cr.newError("请求头 %s 未设置", requiredHeader)
			}
		}
	}
	var body any
	if rd.Body != nil {
		bodyJsonStr, err := json.Marshal(rd.Body)
		if err != nil {
			return nil, cr.newError("序列化请求体异常: %v", err)
		}
		body = cr.replacePlaceholders(string(bodyJsonStr))
	} else {
		body = cr.rsp.Body
	}
	resp, err := netx.NewHttpx(netx.HttpRequestConfig{
		Ctx:    cr.rsp.Ctx,
		Url:    requestUrl,
		Params: rd.Params,
		Header: header,
		Proxy:  cr.rsp.Proxy,
	}).Request(rd.Method, rd.FormData, body)
	if err != nil {
		return nil, cr.newError("请求站点 %s 异常: %v", requestUrl, err)
	}
	cr.requestUrl = requestUrl
	return resp, nil
}

func (cr *siteRequest) newError(format string, v ...any) error {
	return fmt.Errorf("站点(%s)%s", cr.sc.Name, fmt.Sprintf(format, v...))
}

func (cr *siteRequest) replacePlaceholders(str string) string {
	for placeholder, replacement := range cr.rsp.Env {
		str = strings.ReplaceAll(str, "{"+placeholder+"}", replacement)
	}
	return str
}
func (cr *siteRequest) replaceUrlValues(values url.Values) url.Values {
	newValues := url.Values{}
	for key, vals := range values {
		for _, val := range vals {
			newValues.Add(key, cr.replacePlaceholders(val))
		}
	}
	return newValues
}

func (cr *siteRequest) newRequestUrl(rd RequestDefinition) (string, error) {
	if netx.IsValidHttpUrl(rd.Path) {
		return rd.Path, nil
	} else {
		baseUrl := ""
		if rd.UseApi {
			baseUrl = cr.sc.Api
		} else {
			baseUrl = cr.sc.Domain
		}
		return netx.JoinURL(baseUrl, rd.Path)
	}
}
