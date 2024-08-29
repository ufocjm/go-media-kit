package moviesubject

import (
	"fmt"
	"github.com/heibizi/go-media-kit/core/utils/mathx"
	"strconv"
)

func (s *tmdbService) moviePopulars(pageNum int) (Result, error) {
	options := newTmdbOptions(s.params)
	options["page"] = strconv.Itoa(pageNum)
	t, err := s.client.GetMoviePopular(options)
	if err != nil {
		return Result{}, fmt.Errorf("get tmdb movie popular err: %v", err)
	}
	var list []Media
	for _, item := range t.Results {
		list = append(list, Media{
			Id:       strconv.FormatInt(item.ID, 10),
			Title:    item.Title,
			Type:     MediaTypeMovie,
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

func (s *tmdbService) movieNowPlayings(pageNum int) (Result, error) {
	options := newTmdbOptions(s.params)
	options["page"] = strconv.Itoa(pageNum)
	t, err := s.client.GetMovieNowPlaying(options)
	if err != nil {
		return Result{}, fmt.Errorf("get tmdb movie now playing err: %v", err)
	}
	var list []Media
	for _, item := range t.Results {
		list = append(list, Media{
			Id:       strconv.FormatInt(item.ID, 10),
			Title:    item.Title,
			Type:     MediaTypeMovie,
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

func (s *tmdbService) movieTopRated(pageNum int) (Result, error) {
	options := newTmdbOptions(s.params)
	options["page"] = strconv.Itoa(pageNum)
	t, err := s.client.GetMovieTopRated(options)
	if err != nil {
		return Result{}, fmt.Errorf("get tmdb movie top rated err: %v", err)
	}
	var list []Media
	for _, item := range t.Results {
		list = append(list, Media{
			Id:       strconv.FormatInt(item.ID, 10),
			Title:    item.Title,
			Type:     MediaTypeMovie,
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

func (s *tmdbService) movieUpcoming(pageNum int) (Result, error) {
	options := newTmdbOptions(s.params)
	options["page"] = strconv.Itoa(pageNum)
	t, err := s.client.GetMovieUpcoming(options)
	if err != nil {
		return Result{}, fmt.Errorf("get tmdb movie upcoming err: %v", err)
	}
	var list []Media
	for _, item := range t.Results {
		list = append(list, Media{
			Id:       strconv.FormatInt(item.ID, 10),
			Title:    item.Title,
			Type:     MediaTypeMovie,
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
