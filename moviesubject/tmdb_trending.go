package moviesubject

import (
	"fmt"
	"github.com/heibizi/go-media-kit/core/utils/mathx"
	"strconv"
)

func (s *tmdbService) trendingItems(mediaType string, timeWindow string, pageNum int) (Result, error) {
	options := newTmdbOptions(s.params)
	options["page"] = strconv.Itoa(pageNum)
	t, err := s.client.GetTrending(mediaType, timeWindow, options)
	if err != nil {
		return Result{}, fmt.Errorf("获取 tmdb 趋势异常: %v", err)
	}
	var list []Media
	for _, item := range t.Results {
		list = append(list, Media{
			Id:       strconv.FormatInt(item.ID, 10),
			Title:    item.Title,
			Type:     mediaType,
			Year:     getTmdbYearFrom(item.ReleaseDate),
			Vote:     mathx.Round(float64(item.VoteAverage), 2),
			Image:    s.getTmdbImage(item.PosterPath),
			Overview: item.Overview,
		})
	}
	return Result{
		PageNum:  pageNum,
		PageSize: len(list),
		Total:    t.TotalResults,
		List:     list,
	}, nil
}
