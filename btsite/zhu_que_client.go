package btsite

import (
	"context"
	"fmt"
	"github.com/heibizi/go-media-kit/core/utils/netx"
	"github.com/heibizi/go-media-kit/core/utils/stringx"
	"github.com/heibizi/go-media-kit/siteadapt"
	"strings"
	"time"
)

type (
	zhuQueClient struct {
		ctx  context.Context
		site *Site
		*npClient
	}
)

func (c *zhuQueClient) Search(searchParams SearchParams) ([]SearchTorrent, error) {
	sh := SiteHelper
	sc, err := sh.GetConfigByCode(c.site.Code)
	if err != nil {
		return nil, newError(c.site, err, "未获取到站点配置")
	}
	if searchParams.PageSize == 0 {
		searchParams.PageSize = 20
	}
	body := map[string]any{
		"connStatus":     0,
		"downloadStatus": 0,
		"more":           false,
		"page":           searchParams.PageNum,
		"size":           searchParams.PageSize,
		"type":           "title",
	}
	if searchParams.MediaType != nil {
		mediaType := *searchParams.MediaType
		var cats []AdaptMediaCat
		if mediaType == Movie {
			cats = sc.Categories.Movie
		} else if mediaType == Tv {
			cats = sc.Categories.TV
		}
		if len(cats) > 0 {
			var categories []int
			for _, cat := range sc.Categories.Movie {
				categories = append(categories, stringx.ParseInt(cat.Id))
			}
			body["category"] = categories
		}
	}
	if len(searchParams.Keyword) > 0 {
		body["keyword"] = searchParams.Keyword
	}
	site := c.site
	var searchTorrents []SearchTorrent
	domain := ""
	err = list(requestSiteParams{
		ctx:   c.ctx,
		site:  site,
		reqId: requestIdSearch,
		body:  body,
	}, &searchTorrents, func(result siteadapt.Result) {
		domain = result.Domain
	})
	if err != nil {
		return nil, newError(site, err, "搜索异常")
	}
	var torrents []SearchTorrent
	for _, torrent := range searchTorrents {
		pageUrl, err := netx.JoinUrl(domain, torrent.PageUrl)
		if err != nil {
			return nil, err
		}
		torrent.PageUrl = pageUrl
		torrents = append(torrents, torrent)
	}
	return torrents, nil
}

type zhuQueNotice struct {
	Title      string `mapstructure:"title"`
	Content    string `mapstructure:"content"`
	CreateTime int64  `mapstructure:"create_time"`
}

func (c *zhuQueClient) Notice() (string, error) {
	var ns []zhuQueNotice
	err := list(requestSiteParams{
		ctx:   c.ctx,
		site:  c.site,
		reqId: requestIdNotice,
	}, &ns, nil)
	if err != nil {
		return "", newError(c.site, err, "解析公告失败")
	}
	if len(ns) == 0 {
		return "", nil
	}
	var l []string
	for _, n := range ns {
		date := time.Unix(n.CreateTime, 0).Format("2006-01-02 15:04:05")
		l = append(l, fmt.Sprintf("%s - %s\n%s", date, n.Title, n.Content))
	}
	return strings.Join(l, "\n"), nil
}

func (c *zhuQueClient) GetSubtitleDownloadUrl(_ string) (string, error) {
	return "", nil
}
