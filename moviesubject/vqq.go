package moviesubject

import (
	"errors"
	"github.com/heibizi/go-media-kit/core/utils/netx"
	"github.com/heibizi/go-media-kit/core/utils/stringx"
	"github.com/tidwall/gjson"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type vqqService struct {
}

func (s *vqqService) Items(code string, pageNum int) (Result, error) {
	if service, ok := vqqServiceRegistry[code]; ok {
		return service(s, code, pageNum)
	}
	return Result{}, errors.New("invalid vqq code")
}

var vqqServiceRegistry = map[string]func(s *vqqService, code string, pageNum int) (Result, error){
	VqqMoviePopular.Code: func(s *vqqService, code string, pageNum int) (Result, error) {
		return s.getPageData(MediaTypeMovie, "100173", "75", pageNum)
	},
	VqqMovieNowPlaying.Code: func(s *vqqService, code string, pageNum int) (Result, error) {
		return s.getPageData(MediaTypeMovie, "100173", "83", pageNum)
	},
	VqqMovieTopRated.Code: func(s *vqqService, code string, pageNum int) (Result, error) {
		return s.getPageData(MediaTypeMovie, "100173", "81", pageNum)
	},
	VqqTvPopular.Code: func(s *vqqService, code string, pageNum int) (Result, error) {
		return s.getPageData(MediaTypeTv, "100113", "75", pageNum)
	},
	VqqTvNowPlaying.Code: func(s *vqqService, code string, pageNum int) (Result, error) {
		return s.getPageData(MediaTypeTv, "100113", "79", pageNum)
	},
	VqqTvTopRated.Code: func(s *vqqService, code string, pageNum int) (Result, error) {
		return s.getPageData(MediaTypeTv, "100113", "16", pageNum)
	},
}

func (s *vqqService) getPageData(mediaType string, channelId string, sort string, pageNum int) (Result, error) {
	urlStr := "https://pbaccess.video.qq.com/trpc.universal_backend_service.page_server_rpc.PageServer/GetPageData"
	p := url.Values{
		"video_appid":            {"1000005"},
		"vplatform":              {"2"},
		"vversion_name":          {"8.9.10"},
		"new_mark_label_enabled": {"1"},
	}
	data := map[string]any{
		"page_params": map[string]any{
			"channel_id":    channelId,
			"filter_params": "sort=" + sort + "&ifeature=-1&ipay=-1&iarea=-1&iyear=-1&recommend_4=-1",
			"page_type":     "channel_operation",
			"page_id":       "channel_list_second_page",
		},
		"page_context": map[string]any{
			"page_index": strconv.Itoa(pageNum),
		},
	}
	resp, err := netx.NewHttpx(netx.HttpRequestConfig{Url: urlStr, Params: p, Referer: "https://v.qq.com/"}).Post(nil, data)
	if err != nil {
		return Result{}, err
	}
	defer netx.Close(resp)
	body, err := netx.GetBody(resp)
	if err != nil {
		return Result{}, err
	}
	r := gjson.ParseBytes(body)
	d := r.Get("data.module_list_datas.1.module_datas.0")
	var list []Media
	for _, result := range d.Get("item_data_lists.item_datas").Array() {
		if result.Get("item_type").Int() != 2 {
			continue
		}
		itemParams := result.Get("item_params")
		vote := 0.0
		if mediaType == MediaTypeMovie {
			imgTag := itemParams.Get("uni_imgtag").String()
			imgTag = strings.Trim(imgTag, "\"")
			re := regexp.MustCompile(`\d+(\.\d+)?`)
			match := re.FindString(gjson.Parse(imgTag).Get("tag_4.text").String())
			vote = stringx.ParseFloat64(match)
		} else {
			re := regexp.MustCompile(`评分 (\d+(\.\d+)?)`)
			matches := re.FindStringSubmatch(itemParams.Get("third_title").String())
			if len(matches) > 1 {
				vote = stringx.ParseFloat64(matches[1])
			}
		}
		list = append(list, Media{
			Id:       itemParams.Get("cid").String(),
			Title:    itemParams.Get("title").String(),
			Type:     mediaType,
			Year:     getTmdbYearFrom(itemParams.Get("publish_date").String()),
			Vote:     vote,
			Image:    itemParams.Get("new_pic_vt").String(),
			Overview: itemParams.Get("second_title").String(),
		})
	}
	return Result{
		PageNum:  pageNum,
		PageSize: int(d.Get("module_params.page_size").Int()),
		Total:    d.Get("module_params.total_video").Int(),
		List:     list,
	}, nil
}
