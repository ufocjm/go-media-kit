package moviesubject

import (
	"errors"
	"github.com/heibizi/go-media-kit/douban"
)

type MovieSubject struct {
	doubanService *doubanService
	tmdbService   *tmdbService
	vqqService    *vqqService
	iqyService    *iqyService
	ykService     *ykService
	mgService     *mgService
}

func NewMovieSubject() *MovieSubject {
	ms := MovieSubject{}
	ms.doubanService = &doubanService{douban.NewApiClient()}
	ms.vqqService = &vqqService{}
	ms.iqyService = &iqyService{}
	ms.ykService = &ykService{}
	ms.mgService = &mgService{}
	return &ms
}

func (ms *MovieSubject) SetTmdbApiParams(params TmdbApiParams) error {
	if params.ApiKey != "" {
		client, err := newTmdbClient(params)
		if err != nil {
			return err
		}
		ms.tmdbService = &tmdbService{
			client: client,
			params: params,
		}
	} else {
		ms.tmdbService = nil
	}
	return nil
}

func (ms *MovieSubject) SubjectInfos() []SubjectInfo {
	return subjectInfos
}

func (ms *MovieSubject) Categories() []Category {
	var categories []Category
	for _, subjectInfo := range subjectInfos {
		categories = append(categories, subjectInfo.Category)
	}
	return categories
}

func (ms *MovieSubject) Subjects(category string) []Subject {
	for _, subjectInfo := range subjectInfos {
		if subjectInfo.Category.Code == category {
			return subjectInfo.Subjects
		}
	}
	return nil
}

var serviceRegistry = map[string]func(ms *MovieSubject, code string, pageNum int, pageSize int) (Result, error){
	CategoryDoubanRanks.Code: func(ms *MovieSubject, code string, pageNum int, pageSize int) (Result, error) {
		return ms.doubanService.subjectCollectionItems(code, pageNum, pageSize)
	},
	CategoryDoubanYearRanks.Code: func(ms *MovieSubject, code string, pageNum int, pageSize int) (Result, error) {
		return ms.doubanService.subjectCollectionItems(code, pageNum, pageSize)
	},
	CategoryDoubanDouList.Code: func(ms *MovieSubject, code string, pageNum int, pageSize int) (Result, error) {
		return ms.doubanService.douListItems(code, pageNum, pageSize)
	},
	CategoryDoubanSubjectCollection.Code: func(ms *MovieSubject, code string, pageNum int, pageSize int) (Result, error) {
		return ms.doubanService.subjectCollectionItems(code, pageNum, pageSize)
	},
	CategoryTmdb.Code: func(ms *MovieSubject, code string, pageNum int, pageSize int) (Result, error) {
		return ms.tmdbService.Items(code, pageNum)
	},
	CategoryVqq.Code: func(ms *MovieSubject, code string, pageNum int, pageSize int) (Result, error) {
		return ms.vqqService.Items(code, pageNum)
	},
	CategoryIqy.Code: func(ms *MovieSubject, code string, pageNum int, pageSize int) (Result, error) {
		return ms.iqyService.Items(code, pageNum, pageSize)
	},
	CategoryYk.Code: func(ms *MovieSubject, code string, pageNum int, pageSize int) (Result, error) {
		return ms.ykService.Items(code, pageNum)
	},
	CategoryMg.Code: func(ms *MovieSubject, code string, pageNum int, pageSize int) (Result, error) {
		return ms.mgService.Items(code, pageNum, pageSize)
	},
}

func (ms *MovieSubject) Items(req ItemsRequest, pageNum int, pageSize int) (Result, error) {
	if service, ok := serviceRegistry[req.Category]; ok {
		return service(ms, req.Subject, pageNum, pageSize)
	}
	return Result{}, errors.New("category not found")
}
