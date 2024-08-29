package moviesubject

var doubanSubjectCollectionRegistry = map[string]func(s *doubanService, code string, pageNum int, pageSize int) (Result, error){
	DoubanRanksTvComingSoon.Code: func(s *doubanService, code string, pageNum int, pageSize int) (Result, error) {
		return s.comingSoonItems(code, pageNum, pageSize)
	},
	DoubanRanksMovieComingSoonDomestic.Code: func(s *doubanService, code string, pageNum int, pageSize int) (Result, error) {
		return s.comingSoonItems(code, pageNum, pageSize)
	},
	DoubanRanksMovieComingSoonInternational.Code: func(s *doubanService, code string, pageNum int, pageSize int) (Result, error) {
		return s.comingSoonItems(code, pageNum, pageSize)
	},
	DoubanRanksMovieShowing.Code: func(s *doubanService, code string, pageNum int, pageSize int) (Result, error) {
		return s.movieShowItems(pageNum, pageSize)
	},
	DoubanRanksMovieHotGaia.Code: func(s *doubanService, code string, pageNum int, pageSize int) (Result, error) {
		return s.movieHotGaiaItems(pageNum, pageSize)
	},
}

func (s *doubanService) subjectCollectionItems(code string, pageNum int, pageSize int) (Result, error) {
	if f, ok := doubanSubjectCollectionRegistry[code]; ok {
		return f(s, code, pageNum, pageSize)
	}
	r, err := s.client.SubjectCollectionItems(code, (pageNum-1)*pageSize, pageSize)
	if err != nil {
		return Result{}, err
	}
	var list []Media
	for _, item := range r.SubjectCollectionItems {
		image := item.CoverURL
		if item.CoverURL == "" {
			image = item.Pic.Large
		}
		list = append(list, Media{
			Id:       item.Id,
			Title:    item.Title,
			Type:     item.Type,
			Year:     getDoubanYearFrom(item.CardSubtitle),
			Vote:     item.Rating.Value,
			Image:    image,
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
