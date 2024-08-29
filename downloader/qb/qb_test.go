package qb_test

import (
	"encoding/json"
	"github.com/heibizi/go-media-kit/downloader"
	"github.com/heibizi/go-media-kit/downloader/qb"
	"os"
	"testing"
)

var client downloader.Client

func TestMain(m *testing.M) {
	client = downloader.NewClient(downloader.QbClientType, qb.Config{
		Host:     os.Getenv("GO_DOWNLOADER_QB_HOST"),
		Username: os.Getenv("GO_DOWNLOADER_QB_USERNAME"),
		Password: os.Getenv("GO_DOWNLOADER_QB_PASSWORD"),
	})
	err := client.Login()
	if err != nil {
		panic(err)
	}
	m.Run()
}

func TestVersion(t *testing.T) {
	version, err := client.Version()
	log(version, err, t)
}

func TestAddTorrents(t *testing.T) {
	err := client.AddTorrents(qb.AddTorrentsInfo{
		Torrents: []string{os.Getenv("GO_DOWNLOADER_TORRENT1"), os.Getenv("GO_DOWNLOADER_TORRENT2")},
		SavePath: "/media3/downloads",
		Category: "电影",
		Tags:     "test",
		Paused:   true,
	})
	log(nil, err, t)
}

func log(v any, err error, t *testing.T) {
	if err != nil {
		t.Log(err)
		return
	}
	j, _ := json.MarshalIndent(v, "", "    ")
	t.Log(string(j))
}
