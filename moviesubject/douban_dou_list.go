package moviesubject

func (s *doubanService) douListItems(code string, pageNum int, pageSize int) (Result, error) {
	r, err := s.client.DouListItems(code, (pageNum-1)*pageSize, pageSize)
	if err != nil {
		return Result{}, err
	}
	var list []Media
	for _, item := range r.Items {
		list = append(list, Media{
			Id:       item.Id,
			Title:    item.Title,
			Type:     item.Type,
			Year:     getDoubanYearFrom(item.Subtitle),
			Vote:     item.Rating.Value,
			Image:    item.CoverURL,
			Overview: item.Subtitle,
		})
	}
	return Result{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    r.Total,
		List:     list,
	}, nil
}
