package btsite_test

import (
	"context"
	"encoding/json"
	"github.com/heibizi/go-media-kit/btsite"
	"os"
	"testing"
)

var client btsite.Client

func TestMain(m *testing.M) {
	btsite.InitConfig(os.Getenv("GO_BTSITE_CONFIGS_PATH"))
	client, _ = btsite.NewClient(context.Background(), &btsite.Site{
		Code:      os.Getenv("GO_BTSITE_CODE"),
		Name:      os.Getenv("GO_BTSITE_NAME"),
		UserAgent: os.Getenv("GO_BTSITE_UA"),
		Cookie:    os.Getenv("GO_BTSITE_COOKIE"),
		RssUrl:    os.Getenv("GO_BTSITE_RSS_URL"),
	})
	m.Run()
}

func TestUserBasicInfo(t *testing.T) {
	info, err := client.UserBasicInfo()
	log(info, err, t)
	notice, err := client.Notice()
	log(notice, err, t)
}

func TestUserDetails(t *testing.T) {
	details, err := client.UserDetails()
	log(details, err, t)
}

func TestSeedingStatistics(t *testing.T) {
	statistics, err := client.SeedingStatistics()
	log(statistics, err, t)
}

func TestFavicon(t *testing.T) {
	favicon, err := client.Favicon()
	log(favicon, err, t)
}

func TestMyHr(t *testing.T) {
	hr, err := client.MyHr()
	log(hr, err, t)
}

func TestMessage(t *testing.T) {
	messages, err := client.Messages(1)
	for _, message := range messages {
		if message.Content == "" && message.Link != "" {
			content, err := client.MessageDetail(message)
			t.Log(content, err)
		}
	}
	log(messages, err, t)
}

func TestNotice(t *testing.T) {
	notice, err := client.Notice()
	log(notice, err, t)
}

func TestSignIn(t *testing.T) {
	r, err := client.SignIn()
	log(r, err, t)
}

func TestDetails(t *testing.T) {
	details, err := client.Details(os.Getenv("GO_BTSITE_TORRENT_ID"))
	log(details, err, t)
}

func TestSearch(t *testing.T) {
	torrents, err := client.Search(btsite.SearchParams{
		Keyword:   "流浪地球",
		MediaType: nil,
		PageNum:   0,
	})
	log(torrents, err, t)
	if len(torrents) > 0 {
		url, err := client.GetDownloadUrl(torrents[0])
		if err != nil {
			t.Log(err)
		} else {
			log(url, err, t)
			data, err := client.DownloadFile(url)
			log(data, err, t)
		}
	}
}

func TestRss(t *testing.T) {
	rss, err := client.Rss()
	log(rss, err, t)
}

func TestDownloadSubtitle(t *testing.T) {
	subtitle, err := client.GetSubtitleDownloadUrl(os.Getenv("GO_BTSITE_TORRENT_ID"))
	log(subtitle, err, t)
}

func log(v any, err error, t *testing.T) {
	if err != nil {
		t.Log(err)
		return
	}
	j, _ := json.MarshalIndent(v, "", "    ")
	t.Log(string(j))
}
