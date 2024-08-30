package moviesubject

import (
	"errors"
	"github.com/heibizi/go-media-kit/core/utils/netx"
	"github.com/tidwall/gjson"
	"net/url"
	"strconv"
)

type mgService struct{}

type mgListParams struct {
	channelId  string
	sort       string
	chargeInfo string
	fitAge     string
	pageNum    int
	pageSize   int
}

var mgServiceRegistry = map[string]func(s *mgService, code string, pageNum int, pageSize int) (Result, error){
	MgTvVarietyPopular.Code: func(s *mgService, code string, pageNum int, pageSize int) (Result, error) {
		return s.list(MediaTypeTv, mgListParams{channelId: "1", sort: "c2", chargeInfo: "a1", pageNum: pageNum, pageSize: pageSize})
	},
	MgTvVarietyNewly.Code: func(s *mgService, code string, pageNum int, pageSize int) (Result, error) {
		return s.list(MediaTypeTv, mgListParams{channelId: "1", sort: "c1", chargeInfo: "a1", pageNum: pageNum, pageSize: pageSize})
	},
	MgTvPopular.Code: func(s *mgService, code string, pageNum int, pageSize int) (Result, error) {
		return s.list(MediaTypeTv, mgListParams{channelId: "2", sort: "c2", chargeInfo: "a1", pageNum: pageNum, pageSize: pageSize})
	},
	MgTvNewly.Code: func(s *mgService, code string, pageNum int, pageSize int) (Result, error) {
		return s.list(MediaTypeTv, mgListParams{channelId: "2", sort: "c1", chargeInfo: "a1", pageNum: pageNum, pageSize: pageSize})
	},
	MgMoviePopular.Code: func(s *mgService, code string, pageNum int, pageSize int) (Result, error) {
		return s.list(MediaTypeMovie, mgListParams{channelId: "3", sort: "c2", chargeInfo: "a1", pageNum: pageNum, pageSize: pageSize})
	},
	MgMovieNewly.Code: func(s *mgService, code string, pageNum int, pageSize int) (Result, error) {
		return s.list(MediaTypeMovie, mgListParams{channelId: "3", sort: "c1", chargeInfo: "a1", pageNum: pageNum, pageSize: pageSize})
	},
	MgVipVariety.Code: func(s *mgService, code string, pageNum int, pageSize int) (Result, error) {
		return s.list(MediaTypeTv, mgListParams{channelId: "1", sort: "c2", chargeInfo: "b2", pageNum: pageNum, pageSize: pageSize})
	},
	MgVipTv.Code: func(s *mgService, code string, pageNum int, pageSize int) (Result, error) {
		return s.list(MediaTypeTv, mgListParams{channelId: "2", sort: "c2", chargeInfo: "b2", pageNum: pageNum, pageSize: pageSize})
	},
	MgVipMovie.Code: func(s *mgService, code string, pageNum int, pageSize int) (Result, error) {
		return s.list(MediaTypeMovie, mgListParams{channelId: "3", sort: "c2", chargeInfo: "b2", pageNum: pageNum, pageSize: pageSize})
	},
	MgTvChild.Code: func(s *mgService, code string, pageNum int, pageSize int) (Result, error) {
		return s.list(MediaTypeTv, mgListParams{channelId: "10", sort: "c2", fitAge: "197", pageNum: pageNum, pageSize: pageSize})
	},
	MgTvAnimation.Code: func(s *mgService, code string, pageNum int, pageSize int) (Result, error) {
		return s.list(MediaTypeTv, mgListParams{channelId: "50", sort: "c2", pageNum: pageNum, pageSize: pageSize})
	},
	MgTvDocument.Code: func(s *mgService, code string, pageNum int, pageSize int) (Result, error) {
		return s.list(MediaTypeTv, mgListParams{channelId: "51", sort: "c2", pageNum: pageNum, pageSize: pageSize})
	},
	MgTvEducation.Code: func(s *mgService, code string, pageNum int, pageSize int) (Result, error) {
		return s.list(MediaTypeTv, mgListParams{channelId: "115", sort: "c2", pageNum: pageNum, pageSize: pageSize})
	},
}

func (s *mgService) Items(code string, pageNum int, pageSize int) (Result, error) {
	if service, ok := mgServiceRegistry[code]; ok {
		return service(s, code, pageNum, pageSize)
	}
	return Result{}, errors.New("invalid mg code")
}

func (s *mgService) list(mediaType string, params mgListParams) (Result, error) {
	p := url.Values{
		"allowedRC": {"1"},
		"platform":  {"pcweb"},
		"channelId": {params.channelId},
		"pn":        {strconv.Itoa(params.pageNum)},
		"pc":        {strconv.Itoa(params.pageSize)},
		"hudong":    {"1"},
		"_support":  {"10000000"},
		"kind":      {"a1"},
		"edition":   {"a1"},
		"year":      {"all"},
		"sort":      {params.sort},
		"area":      {"a1"},
	}
	if params.chargeInfo != "" {
		p.Set("chargeInfo", params.chargeInfo)
	}
	if params.fitAge != "" {
		p.Set("fitAge", params.fitAge)
	}
	resp, err := netx.NewHttpx(netx.HttpRequestParams{
		Url:     "https://pianku.api.mgtv.com/rider/list/pcweb/v3",
		Params:  p,
		Referer: "https://www.mgtv.com/",
	}).Request()
	if err != nil {
		return Result{}, err
	}
	defer netx.Close(resp)
	body, err := netx.GetBody(resp)
	if err != nil {
		return Result{}, err
	}
	r := gjson.ParseBytes(body)
	var list []Media
	for _, result := range r.Get("data.hitDocs").Array() {
		list = append(list, Media{
			Id:           result.Get("clipId").String(),
			Title:        result.Get("title").String(),
			Type:         mediaType,
			Year:         int(result.Get("year").Int()),
			Vote:         result.Get("zhihuScore").Float(),
			Image:        result.Get("img").String(),
			ImageDynamic: "",
			Overview:     result.Get("story").String(),
		})
	}
	return Result{
		PageNum:  params.pageNum,
		PageSize: params.pageSize,
		Total:    r.Get("data.totalHits").Int(),
		List:     list,
	}, nil
}
