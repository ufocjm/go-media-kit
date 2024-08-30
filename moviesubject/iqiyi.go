package moviesubject

import (
	"encoding/json"
	"errors"
	"github.com/heibizi/go-media-kit/core/utils/netx"
	"github.com/tidwall/gjson"
	"net/url"
	"strconv"
	"strings"
)

type iqyService struct {
}
type iqyDataParams struct {
	channelId string
	mode      string
	pageNum   int
	pageSize  int
	filter    map[string]string
}

func (s *iqyService) Items(code string, pageNum int, pageSize int) (Result, error) {
	if service, ok := iqyServiceRegistry[code]; ok {
		return service(s, code, pageNum, pageSize)
	}
	return Result{}, errors.New("invalid iqy code")
}

var iqyServiceRegistry = map[string]func(s *iqyService, code string, pageNum int, pageSize int) (Result, error){
	IqyMovieComprehensive.Code: func(s *iqyService, code string, pageNum int, pageSize int) (Result, error) {
		return s.data(MediaTypeMovie, iqyDataParams{channelId: "1", pageNum: pageNum, pageSize: pageSize, filter: map[string]string{"mode": "24"}})
	},
	IqyMoviePopular.Code: func(s *iqyService, code string, pageNum int, pageSize int) (Result, error) {
		return s.data(MediaTypeMovie, iqyDataParams{channelId: "1", pageNum: pageNum, pageSize: pageSize, filter: map[string]string{"mode": "11"}})
	},
	IqyMovieNowPlaying.Code: func(s *iqyService, code string, pageNum int, pageSize int) (Result, error) {
		return s.data(MediaTypeMovie, iqyDataParams{channelId: "1", pageNum: pageNum, pageSize: pageSize, filter: map[string]string{"mode": "4"}})
	},
	IqyMovieTopRated.Code: func(s *iqyService, code string, pageNum int, pageSize int) (Result, error) {
		return s.data(MediaTypeMovie, iqyDataParams{channelId: "1", pageNum: pageNum, pageSize: pageSize, filter: map[string]string{"mode": "8"}})
	},
	IqyTvComprehensive.Code: func(s *iqyService, code string, pageNum int, pageSize int) (Result, error) {
		return s.data(MediaTypeTv, iqyDataParams{channelId: "2", pageNum: pageNum, pageSize: pageSize, filter: map[string]string{"mode": "24"}})
	},
	IqyTvPopular.Code: func(s *iqyService, code string, pageNum int, pageSize int) (Result, error) {
		return s.data(MediaTypeTv, iqyDataParams{channelId: "2", pageNum: pageNum, pageSize: pageSize, filter: map[string]string{"mode": "11"}})
	},
	IqyTvNowPlaying.Code: func(s *iqyService, code string, pageNum int, pageSize int) (Result, error) {
		return s.data(MediaTypeTv, iqyDataParams{channelId: "2", pageNum: pageNum, pageSize: pageSize, filter: map[string]string{"mode": "4"}})
	},
	IqyTvTopRated.Code: func(s *iqyService, code string, pageNum int, pageSize int) (Result, error) {
		return s.data(MediaTypeTv, iqyDataParams{channelId: "2", pageNum: pageNum, pageSize: pageSize, filter: map[string]string{"mode": "8"}})
	},
	IqyVipMovie.Code: func(s *iqyService, code string, pageNum int, pageSize int) (Result, error) {
		return s.data(MediaTypeMovie, iqyDataParams{channelId: "1", pageNum: pageNum, pageSize: pageSize,
			filter: map[string]string{"mode": "24", "charge_control_paymark": "1_1_1", "is_purchase": "1"}})
	},
	IqyVipTv.Code: func(s *iqyService, code string, pageNum int, pageSize int) (Result, error) {
		return s.data(MediaTypeTv, iqyDataParams{channelId: "2", pageNum: pageNum, pageSize: pageSize,
			filter: map[string]string{"mode": "24", "charge_control_paymark": "1_1_1", "is_purchase": "1"}})
	},
	IqyVipVarietyShow.Code: func(s *iqyService, code string, pageNum int, pageSize int) (Result, error) {
		return s.data(MediaTypeTv, iqyDataParams{channelId: "6", pageNum: pageNum, pageSize: pageSize,
			filter: map[string]string{"mode": "24", "charge_control_paymark": "1_1_1", "is_purchase": "1"}})
	},
	IqyVipAnimation.Code: func(s *iqyService, code string, pageNum int, pageSize int) (Result, error) {
		return s.data(MediaTypeTv, iqyDataParams{channelId: "4", pageNum: pageNum, pageSize: pageSize,
			filter: map[string]string{"mode": "24", "charge_control_paymark": "1_1_1", "is_purchase": "1"}})
	},
	IqyVipDocumentary.Code: func(s *iqyService, code string, pageNum int, pageSize int) (Result, error) {
		return s.data(MediaTypeTv, iqyDataParams{channelId: "3", pageNum: pageNum, pageSize: pageSize,
			filter: map[string]string{"mode": "24", "charge_control_paymark": "1_1_1", "is_purchase": "1"}})
	},
	IqyVipChild.Code: func(s *iqyService, code string, pageNum int, pageSize int) (Result, error) {
		return s.data(MediaTypeTv, iqyDataParams{channelId: "15", pageNum: pageNum, pageSize: pageSize,
			filter: map[string]string{"mode": "24", "charge_control_paymark": "1_1_1", "is_purchase": "1"}})
	},
}

func (s *iqyService) data(mediaType string, params iqyDataParams) (Result, error) {
	filter, _ := json.Marshal(params.filter)
	p := url.Values{
		"ret_num":    {strconv.Itoa(params.pageSize)},
		"channel_id": {params.channelId},
		"page_id":    {strconv.Itoa(params.pageNum)},
		"vip":        {"0"},
		"filter":     {strings.ReplaceAll(string(filter), `"`, `\"`)},
	}
	resp, err := netx.NewHttpx(netx.HttpRequestParams{
		Url:     "https://mesh.if.iqiyi.com/portal/lw/videolib/data",
		Params:  p,
		Referer: "https://www.iqiyi.com/",
	}).Request()
	if err != nil {
		return Result{}, err
	}
	body, err := netx.GetBody(resp)
	defer netx.Close(resp)
	if err != nil {
		return Result{}, err
	}
	r := gjson.ParseBytes(body)
	var list []Media
	for _, result := range r.Get("data").Array() {
		list = append(list, Media{
			Id:           result.Get("tv_id").String(),
			Title:        result.Get("title").String(),
			Type:         mediaType,
			Year:         int(result.Get("date.year").Int()),
			Vote:         result.Get("sns_score").Float(),
			Image:        result.Get("image_url_normal").String(),
			ImageDynamic: result.Get("image_url_dynamic_normal").String(),
			Overview:     result.Get("description").String(),
		})
	}
	return Result{
		PageNum:  params.pageNum,
		PageSize: params.pageSize,
		Total:    r.Get("extension.max_result_num").Int(),
		List:     list,
	}, nil
}
