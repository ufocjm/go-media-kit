package moviesubject_test

import (
	"encoding/json"
	"github.com/heibizi/go-media-kit/moviesubject"
	"os"
	"testing"
)

var ms = moviesubject.NewMovieSubject()

func TestMain(m *testing.M) {
	_ = ms.SetTmdbApiParams(moviesubject.TmdbApiParams{
		ApiKey: os.Getenv("GO_MOVIESUBJECT_TMDB_APIKEY"),
		Proxy:  os.Getenv("GO_MOVIESUBJECT_TMDB_PROXY"),
	})
	m.Run()
}

func TestSubjectInfos(t *testing.T) {
	subjects := ms.SubjectInfos()
	t.Log(subjects)
}

func TestSubjectCollection(t *testing.T) {
	subject := moviesubject.ItemsRequest{Category: moviesubject.CategoryDoubanYearRanks.Code, Subject: moviesubject.DoubanRanksMovieHotGaia.Code}
	items, err := ms.Items(subject, 2, 20)
	log(items, err, t)
}

func TestTmdb(t *testing.T) {
	subject := moviesubject.ItemsRequest{Category: moviesubject.CategoryTmdb.Code, Subject: moviesubject.TmdbMoviePopular.Code}
	items, err := ms.Items(subject, 1, 0)
	log(items, err, t)
}

func TestVqq(t *testing.T) {
	subject := moviesubject.ItemsRequest{Category: moviesubject.CategoryVqq.Code, Subject: moviesubject.VqqMoviePopular.Code}
	items, err := ms.Items(subject, 1, 0)
	log(items, err, t)
}

func TestIqy(t *testing.T) {
	subject := moviesubject.ItemsRequest{Category: moviesubject.CategoryIqy.Code, Subject: moviesubject.IqyMoviePopular.Code}
	items, err := ms.Items(subject, 1, 60)
	log(items, err, t)
}

func TestYk(t *testing.T) {
	subject := moviesubject.ItemsRequest{Category: moviesubject.CategoryYk.Code, Subject: moviesubject.YkMoviePopular.Code}
	items, err := ms.Items(subject, 1, 0)
	log(items, err, t)
}

func TestMg(t *testing.T) {
	subject := moviesubject.ItemsRequest{Category: moviesubject.CategoryMg.Code, Subject: moviesubject.MgTvVarietyPopular.Code}
	items, err := ms.Items(subject, 1, 20)
	log(items, err, t)
}

func log(v any, err error, t *testing.T) {
	if err != nil {
		t.Log(err)
		return
	}
	j, _ := json.MarshalIndent(v, "", "    ")
	t.Log(string(j))
}
