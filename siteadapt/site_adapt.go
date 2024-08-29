package siteadapt

import (
	"encoding/json"
	"fmt"
	"github.com/heibizi/go-media-kit/core/utils/mapstructurex"
	"io"
)

func NewSiteAdaptor(sc Config) *SiteAdaptor {
	return &SiteAdaptor{
		sc: sc,
	}
}

type (
	SiteAdaptor struct {
		sc Config
	}
	// RequestSiteResult 站点请求结果
	result struct {
		List     []map[string]any // 列表数据
		Raw      []byte           // 原始数据
		Data     map[string]any   // 键值对数据
		NextPage string           // 下一页
		RequestInfo
	}
	RequestInfo struct {
		Domain     string // 域名
		RequestUrl string // 请求地址
		StatusCode int    // 状态码
	}
	Result struct {
		NextPage string
		Raw      []byte
		RequestInfo
	}
	ResultFunc func(result Result)
)

// List 获取列表数据
func (sa *SiteAdaptor) List(requestSiteParams RequestSiteParams, output any, fn ResultFunc) error {
	r, err := sa.requestSite(requestSiteParams)
	if err != nil {
		return fmt.Errorf("请求异常: %v", err)
	}
	err = mapstructurex.WeakDecode(r.List, output)
	if err != nil {
		return fmt.Errorf("解析异常: %v", err)
	}
	if fn != nil {
		fn(Result{
			NextPage:    r.NextPage,
			Raw:         r.Raw,
			RequestInfo: r.RequestInfo,
		})
	}
	return nil
}

// Data 获取对象数据
func (sa *SiteAdaptor) Data(requestSiteParams RequestSiteParams, output any, fn ResultFunc) error {
	r, err := sa.requestSite(requestSiteParams)
	if err != nil {
		return fmt.Errorf("请求异常: %v", err)
	}
	err = mapstructurex.WeakDecode(r.Data, output)
	if err != nil {
		return fmt.Errorf("解析异常: %v", err)
	}
	if fn != nil {
		fn(Result{
			Raw:         r.Raw,
			RequestInfo: r.RequestInfo,
		})
	}
	return nil
}

// Raw 获取原始数据
func (sa *SiteAdaptor) Raw(requestSiteParams RequestSiteParams, fn ResultFunc) error {
	r, err := sa.requestSite(requestSiteParams)
	if err != nil {
		return fmt.Errorf("请求异常: %v", err)
	}
	if fn != nil {
		fn(Result{
			Raw:         r.Raw,
			RequestInfo: r.RequestInfo,
		})
	}
	return nil
}

// Json 直接 JSON 转 Struct
func (sa *SiteAdaptor) Json(requestSiteParams RequestSiteParams, output interface{}) error {
	r, err := sa.requestSite(requestSiteParams)
	if err != nil {
		return fmt.Errorf("请求异常: %v", err)
	}
	var input map[string]interface{}
	err = json.Unmarshal(r.Raw, &input)
	if err != nil {
		return err
	}
	err = mapstructurex.WeakDecode(input, &output)
	if err != nil {
		return fmt.Errorf("解析异常: %v", err)
	}
	return nil
}

// 请求站点
func (sa *SiteAdaptor) requestSite(rsp RequestSiteParams) (result, error) {
	reqId := rsp.ReqId
	rd, err := sa.getRd(reqId, rsp)
	if err != nil {
		return result{}, err
	}
	sc := sa.sc
	if rsp.Domain != "" {
		sc.Domain = rsp.Domain
	}
	if rsp.Api != "" {
		sc.Api = rsp.Api
	}
	if rsp.Path != "" {
		rd.Path = rsp.Path
	}
	rsp.Rd = &rd
	sr := newSiteRequest(sc, rsp)
	requestInfo := RequestInfo{
		Domain: sc.Domain,
	}
	var data []byte
	if rsp.SiteData != nil {
		data = rsp.SiteData
	} else {
		resp, err := sr.request(rd)
		if err != nil {
			return result{}, err
		}
		requestInfo.StatusCode = resp.StatusCode
		requestInfo.RequestUrl = sr.requestUrl
		defer resp.Body.Close()
		success := false
		if rd.SuccessStatusCodes == nil {
			if resp.StatusCode == 200 {
				success = true
			}
		} else {
			for _, successStatusCode := range rd.SuccessStatusCodes {
				if resp.StatusCode == successStatusCode {
					success = true
					break
				}
			}
		}
		if !success {
			return result{}, sa.newError("%s 请求失败, 状态码为: %v, 异常: %v", sr.requestUrl, resp.StatusCode, err)
		}
		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return result{}, sa.newError("请求异常: %v", err)
		}
	}
	// 解析响应
	parsedData, err := newParserHelper(data, sa.sc, rd).parse()
	if err != nil {
		return result{}, sa.newError("解析响应异常: %v", err)
	}
	result := result{
		RequestInfo: requestInfo,
		Raw:         data,
	}
	nextPage := parsedData[string(fieldNameNextPage)]
	if nextPage != nil {
		result.NextPage = nextPage.(string)
	}
	if parsedData[string(fieldNameList)] != nil {
		result.List = parsedData[string(fieldNameList)].([]map[string]any)
	} else {
		result.Data = parsedData
	}
	return result, nil
}

func (sa *SiteAdaptor) newError(format string, v ...any) error {
	return fmt.Errorf("站点(%s) %s", sa.sc.Name, fmt.Sprintf(format, v...))
}

func (sa *SiteAdaptor) getRd(reqId string, params RequestSiteParams) (RequestDefinition, error) {
	if params.Rd != nil {
		return *params.Rd, nil
	}
	scRd, exists := sa.sc.RequestDefinitions[reqId]
	if exists {
		return scRd, nil
	}
	return RequestDefinition{}, sa.newError("未适配该请求: %s", reqId)
}
