package moviesubject

import (
	"fmt"
	"github.com/heibizi/go-media-kit/douban"
	"strconv"
	"strings"
)

type doubanService struct {
	client *douban.ApiClient
}

func (s *doubanService) movieShowItems(pageNum int, pageSize int) (Result, error) {
	// 北京 热度
	r, err := s.client.MovieMovieShowing("108288", (pageNum-1)*pageSize, pageSize, "recommend")
	if err != nil {
		return Result{}, err
	}
	var list []Media
	for _, item := range r.Items {
		list = append(list, Media{
			Id:       item.Id,
			Title:    item.Title,
			Type:     item.Type,
			Year:     getDoubanYearFrom(item.CardSubtitle),
			Vote:     item.Rating.Value,
			Image:    item.Pic.Large,
			Overview: item.CardSubtitle,
		})
	}
	return Result{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    r.Total,
		List:     list,
	}, nil
}

func (s *doubanService) movieHotGaiaItems(pageNum int, pageSize int) (Result, error) {
	// 北京 热度
	r, err := s.client.MovieHotGaia("108288", (pageNum-1)*pageSize, pageSize, "recommend")
	if err != nil {
		return Result{}, err
	}
	var list []Media
	for _, item := range r.Items {
		list = append(list, Media{
			Id:       item.Id,
			Title:    item.Title,
			Type:     item.Type,
			Year:     getDoubanYearFrom(item.CardSubtitle),
			Vote:     item.Rating.Value,
			Image:    item.Pic.Large,
			Overview: item.CardSubtitle,
		})
	}
	return Result{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    r.Total,
		List:     list,
	}, nil
}

func (s *doubanService) comingSoonItems(code string, pageNum int, pageSize int) (Result, error) {
	mediaType := douban.Movie
	area := ""
	if code == DoubanRanksMovieComingSoonDomestic.Code {
		area = "domestic"
	} else if code == DoubanRanksMovieComingSoonInternational.Code {
		area = "international"
	} else if code == DoubanRanksTvComingSoon.Code {
		mediaType = douban.Tv
	}
	r, err := s.client.ComingSoon(mediaType, (pageNum-1)*pageSize, pageSize, "hot", area)
	if err != nil {
		return Result{}, err
	}
	var list []Media
	for _, item := range r.Subjects {
		list = append(list, Media{
			Id:       item.Id,
			Title:    item.Title,
			Type:     item.Type,
			Year:     getDoubanYearFrom(item.CardSubtitle),
			Vote:     item.Rating.Value,
			Image:    item.CoverURL,
			Overview: item.CardSubtitle,
		})
	}
	return Result{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    r.Total,
		List:     list,
	}, nil
}

func getDoubanYearFrom(title string) int {
	i, err := strconv.Atoi(strings.TrimSpace(strings.Split(title, "/")[0]))
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return i
}
