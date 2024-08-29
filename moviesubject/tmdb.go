package moviesubject

import (
	"errors"
	"fmt"
	tmdb "github.com/cyruzin/golang-tmdb"
	"github.com/heibizi/go-media-kit/core/utils/netx"
	"net/http"
	"strconv"
	"time"
)

type tmdbService struct {
	client *tmdb.Client
	params TmdbApiParams
}

func (s *tmdbService) Items(code string, pageNum int) (Result, error) {
	if service, ok := tmdbServiceRegistry[code]; ok {
		return service(s, code, pageNum)
	}
	return Result{}, errors.New("invalid tmdb code")
}

var tmdbServiceRegistry = map[string]func(tmdbService *tmdbService, code string, pageNum int) (Result, error){
	TmdbMovieTrendingWeek.Code: func(tmdbService *tmdbService, code string, pageNum int) (Result, error) {
		return tmdbService.trendingItems(MediaTypeMovie, "week", pageNum)
	},
	TmdbMovieTrendingDay.Code: func(tmdbService *tmdbService, code string, pageNum int) (Result, error) {
		return tmdbService.trendingItems(MediaTypeMovie, "day", pageNum)
	},
	TmdbTvTrendingWeek.Code: func(tmdbService *tmdbService, code string, pageNum int) (Result, error) {
		return tmdbService.trendingItems(MediaTypeTv, "week", pageNum)
	},
	TmdbTvTrendingDay.Code: func(tmdbService *tmdbService, code string, pageNum int) (Result, error) {
		return tmdbService.trendingItems(MediaTypeTv, "day", pageNum)
	},
	TmdbMoviePopular.Code: func(tmdbService *tmdbService, code string, pageNum int) (Result, error) {
		return tmdbService.moviePopulars(pageNum)
	},
	TmdbMovieNowPlaying.Code: func(tmdbService *tmdbService, code string, pageNum int) (Result, error) {
		return tmdbService.movieNowPlayings(pageNum)
	},
	TmdbMovieTopRated.Code: func(tmdbService *tmdbService, code string, pageNum int) (Result, error) {
		return tmdbService.movieTopRated(pageNum)
	},
	TmdbMovieUpcoming.Code: func(tmdbService *tmdbService, code string, pageNum int) (Result, error) {
		return tmdbService.movieUpcoming(pageNum)
	},
	TmdbTvPopular.Code: func(tmdbService *tmdbService, code string, pageNum int) (Result, error) {
		return tmdbService.tvPopulars(pageNum)
	},
	TmdbTvTopRated.Code: func(tmdbService *tmdbService, code string, pageNum int) (Result, error) {
		return tmdbService.tvTopRated(pageNum)
	},
	TmdbTvAiringToday.Code: func(tmdbService *tmdbService, code string, pageNum int) (Result, error) {
		return tmdbService.tvAiringToday(pageNum)
	},
	TmdbTvOnTheAir.Code: func(tmdbService *tmdbService, code string, pageNum int) (Result, error) {
		return tmdbService.tvOnTheAir(pageNum)
	},
}

func newTmdbClient(params TmdbApiParams) (*tmdb.Client, error) {
	// apikey 为空才会报错，直接忽略
	tmdbClient, _ := tmdb.Init(params.ApiKey)
	tmdbClient.SetClientAutoRetry()
	if params.CustomBaseURL != "" {
		tmdbClient.SetCustomBaseURL(params.CustomBaseURL)
	}
	if params.Timeout == 0 {
		params.Timeout = time.Second * 60
	}
	if params.MaxIdleConn == 0 {
		params.MaxIdleConn = 10
	}
	if params.IdleConnTimeout == 0 {
		params.IdleConnTimeout = time.Second * 60
	}
	tp, err := netx.NewHttpTransport(params.Proxy)
	if err != nil {
		return nil, err
	}
	if tp == nil {
		tp = &http.Transport{}
	}
	tp.MaxIdleConns = params.MaxIdleConn
	tp.IdleConnTimeout = params.IdleConnTimeout
	customClient := http.Client{
		Timeout:   params.Timeout,
		Transport: tp,
	}
	tmdbClient.SetClientConfig(customClient)
	return tmdbClient, nil
}

func newTmdbOptions(params TmdbApiParams) map[string]string {
	options := make(map[string]string)
	if params.Language != "" {
		options["language"] = params.Language
	} else {
		options["language"] = "zh"
	}
	if params.Region != "" {
		options["region"] = params.Region
	} else {
		options["region"] = "CN"
	}
	options["include_adult"] = strconv.FormatBool(params.IncludeAdult)
	return options
}

func getTmdbYearFrom(date string) int {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		fmt.Println(err)
		return 2006
	}
	return t.Year()
}

func (s *tmdbService) getTmdbImage(url string) string {
	if url == "" {
		return ""
	}
	imageBaseUrl := "https://image.tmdb.org"
	if s.params.ImageBaseURL != "" {
		imageBaseUrl = s.params.ImageBaseURL
	}
	return fmt.Sprintf("%s/t/p/w500%s", imageBaseUrl, url)
}
