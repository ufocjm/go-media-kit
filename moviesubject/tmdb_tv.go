package moviesubject

import (
	"fmt"
	"github.com/heibizi/go-media-kit/core/utils/mathx"
	"strconv"
)

func (s *tmdbService) tvPopulars(pageNum int) (Result, error) {
	options := newTmdbOptions(s.params)
	options["page"] = strconv.Itoa(pageNum)
	t, err := s.client.GetTVPopular(options)
	if err != nil {
		return Result{}, fmt.Errorf("get tmdb tv popular err: %v", err)
	}
	var list []Media
	for _, item := range t.Results {
		list = append(list, Media{
			Id:       strconv.FormatInt(item.ID, 10),
			Title:    item.Name,
			Type:     MediaTypeTv,
			Year:     getTmdbYearFrom(item.FirstAirDate),
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

func (s *tmdbService) tvTopRated(pageNum int) (Result, error) {
	options := newTmdbOptions(s.params)
	options["page"] = strconv.Itoa(pageNum)
	t, err := s.client.GetTVTopRated(options)
	if err != nil {
		return Result{}, fmt.Errorf("get tmdb tv top rated err: %v", err)
	}
	var list []Media
	for _, item := range t.Results {
		list = append(list, Media{
			Id:       strconv.FormatInt(item.ID, 10),
			Title:    item.Name,
			Type:     MediaTypeTv,
			Year:     getTmdbYearFrom(item.FirstAirDate),
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

func (s *tmdbService) tvOnTheAir(pageNum int) (Result, error) {
	options := newTmdbOptions(s.params)
	options["page"] = strconv.Itoa(pageNum)
	t, err := s.client.GetTVOnTheAir(options)
	if err != nil {
		return Result{}, fmt.Errorf("get tmdb tv on the air err: %v", err)
	}
	var list []Media
	for _, item := range t.Results {
		list = append(list, Media{
			Id:       strconv.FormatInt(item.ID, 10),
			Title:    item.Name,
			Type:     MediaTypeTv,
			Year:     getTmdbYearFrom(item.FirstAirDate),
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

func (s *tmdbService) tvAiringToday(pageNum int) (Result, error) {
	options := newTmdbOptions(s.params)
	options["page"] = strconv.Itoa(pageNum)
	t, err := s.client.GetTVAiringToday(options)
	if err != nil {
		return Result{}, fmt.Errorf("get tmdb tv on the air err: %v", err)
	}
	var list []Media
	for _, item := range t.Results {
		list = append(list, Media{
			Id:       strconv.FormatInt(item.ID, 10),
			Title:    item.Name,
			Type:     MediaTypeTv,
			Year:     getTmdbYearFrom(item.FirstAirDate),
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
