package moviesubject

import (
	"encoding/json"
	"errors"
	"github.com/heibizi/go-media-kit/core/utils/netx"
	"github.com/tidwall/gjson"
	"net/url"
	"regexp"
	"strconv"
)

type ykService struct {
}
type ykDataParams struct {
	pageNum int
	params  map[string]string
}

func (s *ykService) Items(code string, pageNum int) (Result, error) {
	if service, ok := ykServiceRegistry[code]; ok {
		return service(s, code, pageNum)
	}
	return Result{}, errors.New("invalid iqy code")
}

var ykServiceRegistry = map[string]func(s *ykService, code string, pageNum int) (Result, error){
	YkMovieComprehensive.Code: func(s *ykService, code string, pageNum int) (Result, error) {
		return s.data(MediaTypeMovie, ykDataParams{pageNum: pageNum, params: map[string]string{"type": "电影"}})
	},
	YkMoviePopular.Code: func(s *ykService, code string, pageNum int) (Result, error) {
		return s.data(MediaTypeMovie, ykDataParams{pageNum: pageNum, params: map[string]string{"type": "电影", "sort": "7"}})
	},
	YkMovieNewly.Code: func(s *ykService, code string, pageNum int) (Result, error) {
		return s.data(MediaTypeMovie, ykDataParams{pageNum: pageNum, params: map[string]string{"type": "电影", "sort": "1"}})
	},
	YkMovieTopRated.Code: func(s *ykService, code string, pageNum int) (Result, error) {
		return s.data(MediaTypeMovie, ykDataParams{pageNum: pageNum, params: map[string]string{"type": "电影", "sort": "3"}})
	},
	YkMovieMostPlayed.Code: func(s *ykService, code string, pageNum int) (Result, error) {
		return s.data(MediaTypeMovie, ykDataParams{pageNum: pageNum, params: map[string]string{"type": "电影", "sort": "2"}})
	},
	YkTvComprehensive.Code: func(s *ykService, code string, pageNum int) (Result, error) {
		return s.data(MediaTypeTv, ykDataParams{pageNum: pageNum, params: map[string]string{"type": "电视剧"}})
	},
	YkTvPopular.Code: func(s *ykService, code string, pageNum int) (Result, error) {
		return s.data(MediaTypeTv, ykDataParams{pageNum: pageNum, params: map[string]string{"type": "电视剧", "sort": "7"}})
	},
	YkTvNewly.Code: func(s *ykService, code string, pageNum int) (Result, error) {
		return s.data(MediaTypeTv, ykDataParams{pageNum: pageNum, params: map[string]string{"type": "电视剧", "sort": "1"}})
	},
	YkTvTopRated.Code: func(s *ykService, code string, pageNum int) (Result, error) {
		return s.data(MediaTypeTv, ykDataParams{pageNum: pageNum, params: map[string]string{"type": "电视剧", "sort": "3"}})
	},
	YkTvMostPlayed.Code: func(s *ykService, code string, pageNum int) (Result, error) {
		return s.data(MediaTypeTv, ykDataParams{pageNum: pageNum, params: map[string]string{"type": "电视剧", "sort": "2"}})
	},
	YkVipMovie.Code: func(s *ykService, code string, pageNum int) (Result, error) {
		return s.data(MediaTypeMovie, ykDataParams{pageNum: pageNum,
			params: map[string]string{"type": "电影", "payType": "2"}})
	},
	YkVipTv.Code: func(s *ykService, code string, pageNum int) (Result, error) {
		return s.data(MediaTypeTv, ykDataParams{pageNum: pageNum,
			params: map[string]string{"type": "电视剧", "payType": "2"}})
	},
	YkVipVarietyShow.Code: func(s *ykService, code string, pageNum int) (Result, error) {
		return s.data(MediaTypeTv, ykDataParams{pageNum: pageNum,
			params: map[string]string{"type": "综艺", "payType": "2"}})
	},
	YkVipAnimation.Code: func(s *ykService, code string, pageNum int) (Result, error) {
		return s.data(MediaTypeTv, ykDataParams{pageNum: pageNum,
			params: map[string]string{"type": "动漫", "payType": "2"}})
	},
	YkVipDocumentary.Code: func(s *ykService, code string, pageNum int) (Result, error) {
		return s.data(MediaTypeTv, ykDataParams{pageNum: pageNum,
			params: map[string]string{"type": "纪录片", "payType": "2"}})
	},
	YkVipChild.Code: func(s *ykService, code string, pageNum int) (Result, error) {
		return s.data(MediaTypeTv, ykDataParams{pageNum: pageNum,
			params: map[string]string{"type": "少儿", "payType": "2"}})
	},
}

func (s *ykService) data(mediaType string, params ykDataParams) (Result, error) {
	filter, _ := json.Marshal(params.params)
	p := url.Values{
		"params":        {string(filter)},
		"optionRefresh": {"1"},
		"pageNo":        {strconv.Itoa(params.pageNum)},
	}
	resp, err := netx.NewHttpx(netx.HttpRequestParams{Url: "https://youku.com/category/data", Params: p, Referer: "https://youku.com/"}).Request()
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
	for _, result := range r.Get("data.filterData.listData").Array() {
		videoLink := result.Get("videoLink").String()
		id := ""
		if videoLink != "" {
			re := regexp.MustCompile(`id_([a-zA-Z0-9]+)\.html`)
			matches := re.FindStringSubmatch(videoLink)
			if len(matches) >= 2 {
				id = matches[1]
			}
		}
		list = append(list, Media{
			Id:           id,
			Title:        result.Get("title").String(),
			Type:         mediaType,
			Year:         0,
			Vote:         result.Get("summary").Float(),
			Image:        result.Get("img").String(),
			ImageDynamic: "",
			Overview:     result.Get("subTitle").String(),
		})
	}
	// 优酷暂时不支持分页
	return Result{
		PageNum:  params.pageNum,
		PageSize: len(list),
		Total:    int64(len(list)),
		List:     list,
	}, nil
}
