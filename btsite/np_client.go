package btsite

import (
	"context"
	"encoding/xml"
	"fmt"
	"github.com/heibizi/go-media-kit/core/utils/netx"
	"github.com/heibizi/go-media-kit/core/utils/stringx"
	"github.com/heibizi/go-media-kit/siteadapt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type (
	// npClient NexusPHP 客户端
	npClient struct {
		ctx      context.Context
		site     *Site
		homeData []byte
	}
)

func (c *npClient) Favicon() ([]byte, error) {
	var data []byte
	err := raw(requestSiteParams{
		ctx:   c.ctx,
		site:  c.site,
		reqId: requestIdFavicon,
	}, func(result siteadapt.Result) {
		data = result.Raw
	})
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *npClient) UserBasicInfo() (UserBasicInfo, error) {
	var ud UserBasicInfo
	err := data(requestSiteParams{
		ctx:   c.ctx,
		site:  c.site,
		reqId: requestIdUserBasicInfo,
	}, &ud, func(result siteadapt.Result) {
		c.homeData = result.Raw
	})
	if err != nil {
		return ud, newError(c.site, err, "解析基础信息失败")
	}
	if ud.Ratio == 0 && ud.Downloaded > 0 {
		ratio := float64(10 ^ 3)
		ud.Ratio = float64(int(float64(ud.Uploaded)/float64(ud.Downloaded)*ratio)) / ratio
	}
	if ud.Bonus == 0 && (ud.Gold > 0 || ud.Silver > 0 || ud.Copper > 0) {
		totalCopper := ud.Gold*100*100 + ud.Silver*100 + ud.Copper
		ud.Bonus = totalCopper
	}
	return ud, nil
}

func (c *npClient) UserDetails() (UserDetails, error) {
	var ud UserDetails
	err := data(requestSiteParams{
		ctx:   c.ctx,
		site:  c.site,
		reqId: requestIdUserDetails,
	}, &ud, nil)
	if err != nil {
		return ud, newError(c.site, err, "解析用户详情信息异常")
	}
	return ud, nil
}

func (c *npClient) Search(searchParams SearchParams) ([]SearchTorrent, error) {
	sh := SiteHelper
	site := c.site
	sc, err := sh.GetConfigByCode(site.Code)
	if err != nil {
		return nil, newError(site, err, "未获取到站点配置")
	}
	var params url.Values
	var env = map[string]string{
		"keyword": searchParams.Keyword,
	}
	if len(searchParams.Keyword) > 0 {
		params = url.Values{
			"search_mode": {"0"},
			"notnewword":  {"1"},
		}
		if searchParams.PageNum > 0 {
			params.Set("page", strconv.Itoa(searchParams.PageNum-1))
		}
		if searchParams.MediaType != nil {
			category := &sc.Categories
			var cats []AdaptMediaCat
			mediaType := *searchParams.MediaType
			if mediaType == Movie {
				cats = category.Movie
			} else if mediaType == Tv {
				cats = category.TV
			}
			for _, cat := range cats {
				field := category.Field
				if len(field) > 0 {
					value := params.Get(field)
					params.Add(field, value+category.Delimiter+cat.Id)
				} else {
					params.Add(cat.Id, "1")
				}
			}
		}
	} else {
		params = url.Values{
			"page": {fmt.Sprintf("%d", searchParams.PageNum)},
		}
	}
	var torrents []SearchTorrent
	requestUrl := ""
	err = list(requestSiteParams{
		ctx:    c.ctx,
		site:   site,
		reqId:  requestIdSearch,
		params: params,
		env:    env,
	}, &torrents, func(result siteadapt.Result) {
		requestUrl = result.RequestUrl
	})
	if err != nil {
		return nil, newError(site, err, "解析种子列表失败")
	}
	var searchTorrents []SearchTorrent
	for _, torrent := range torrents {
		var pageUrl string
		if netx.IsValidHttpUrl(torrent.PageUrl) {
			pageUrl = torrent.PageUrl
		} else {
			detailLink, err := netx.JoinURL(requestUrl, torrent.PageUrl)
			if err != nil {
				return nil, newError(site, err, "搜索种子拼接 details 错误")
			}
			pageUrl = detailLink
		}
		if torrent.Enclosure != "" {
			if !netx.IsValidHttpUrl(torrent.Enclosure) && !strings.HasPrefix(torrent.Enclosure, "magnet") {
				enclosureLink, err := netx.JoinURL(requestUrl, torrent.Enclosure)
				if err != nil {
					return nil, newError(site, err, "搜索种子解析 download 错误")
				}
				torrent.Enclosure = enclosureLink
			}
		}
		var labels []string
		for _, label := range torrent.Labels {
			for _, s := range strings.Split(label, "|") {
				labels = append(labels, s)
			}
		}
		torrent.Labels = labels
		torrent.PageUrl = pageUrl
		searchTorrents = append(searchTorrents, torrent)
	}
	return searchTorrents, nil
}

func (c *npClient) SeedingStatistics() (SeedingStatistics, error) {
	sc, err := SiteHelper.GetConfigByCode(c.site.Code)
	if err != nil {
		return SeedingStatistics{}, err
	}
	rd, exists := sc.RequestDefinitions[string(requestIdSeedingStatistics)]
	// 如果定义了请求，且不走分页
	if exists && rd.List == nil {
		var ss SeedingStatistics
		err := data(requestSiteParams{
			ctx:   c.ctx,
			site:  c.site,
			reqId: requestIdSeedingStatistics,
		}, &ss, nil)
		if err != nil {
			return ss, newError(c.site, err, "做种信息失败")
		}
		return ss, nil
	} else {
		// 分页统计做种信息
		var seedingList []seeding
		page := 1
		currentPageSeedingList, nextPage, err := c.CurrentPageSeeding("", page)
		if err != nil {
			return SeedingStatistics{}, err
		}
		seedingList = append(seedingList, currentPageSeedingList...)
		for len(nextPage) > 0 {
			page++
			currentPageSeedingList, nextPageTmp, err := c.CurrentPageSeeding(nextPage, page)
			if err != nil {
				return SeedingStatistics{}, err
			}
			seedingList = append(seedingList, currentPageSeedingList...)
			nextPage = nextPageTmp
			time.Sleep(500 * time.Millisecond)
		}
		var size int64 = 0
		for _, seeding := range seedingList {
			size = size + seeding.Size
		}
		return SeedingStatistics{
			Count: len(seedingList),
			Size:  size,
		}, nil
	}
}

// CurrentPageSeeding 当前页做种信息以及下一页链接地址
func (c *npClient) CurrentPageSeeding(url string, page int) ([]seeding, string, error) {
	var seedingList []seeding
	nextPage := ""
	err := list(requestSiteParams{
		ctx:   c.ctx,
		site:  c.site,
		reqId: requestIdSeedingStatistics,
		path:  url,
		env:   map[string]string{"page": strconv.Itoa(page)},
	}, &seedingList, func(result siteadapt.Result) {
		nextPage = result.NextPage
	})
	if err != nil {
		return nil, "", newError(c.site, err, "解析做种信息列表失败")
	}
	return seedingList, nextPage, nil
}

func (c *npClient) MyHr() ([]HrTorrent, error) {
	sc, err := SiteHelper.GetConfigByCode(c.site.Code)
	if err != nil {
		return nil, err
	}
	if !sc.Price.HasHR {
		return nil, nil
	}
	var hrList []HrTorrent
	err = list(requestSiteParams{
		ctx:   c.ctx,
		site:  c.site,
		reqId: requestIdMyHr,
	}, &hrList, nil)
	if err != nil {
		return nil, newError(c.site, err, "HR列表失败")
	}
	return hrList, nil
}

func (c *npClient) Messages(page int) ([]Message, error) {
	var o []Message
	requestUrl := ""
	env := map[string]string{"page": strconv.Itoa(page - 1)}
	err := list(requestSiteParams{
		ctx:   c.ctx,
		site:  c.site,
		reqId: requestIdMessages,
		env:   env,
	}, &o, func(result siteadapt.Result) {
		requestUrl = result.RequestUrl
	})
	if err != nil {
		return nil, newError(c.site, err, "获取消息列表异常")
	}
	for i := range o {
		message := &o[i]
		detailUrl, err := netx.JoinURL(requestUrl, message.Link)
		if err != nil {
			return nil, err
		}
		message.Link = detailUrl
	}
	return o, nil
}

func (c *npClient) Notice() (string, error) {
	var n notice
	err := data(requestSiteParams{
		ctx:      c.ctx,
		site:     c.site,
		reqId:    requestIdNotice,
		siteData: c.homeData,
	}, &n, nil)
	if err != nil {
		return "", newError(c.site, err, "解析公告失败")
	}
	return n.Content, nil
}

func (c *npClient) Rss() ([]RssTorrent, error) {
	rd := siteadapt.RequestDefinition{
		Parser: "None",
		Method: http.MethodGet,
		Path:   c.site.RssUrl,
	}
	var data []byte
	err := raw(requestSiteParams{
		ctx:  c.ctx,
		site: c.site,
		rd:   &rd,
	}, func(result siteadapt.Result) {
		data = result.Raw
	})
	if err != nil {
		return nil, newError(c.site, err, "获取 RSS 数据异常")
	}
	var rss rssResult
	err = xml.Unmarshal(data, &rss)
	if err != nil {
		return nil, newError(c.site, err, "解析 RSS xml 数据异常")
	}
	var torrents []RssTorrent
	for _, item := range rss.Items {
		if len(item.Title) == 0 {
			continue
		}
		// todo 月月标题特殊处理
		//if siteDomain != "" {
		// Placeholder for special title processing
		//}
		link := item.Link
		enclosure := item.Enclosure.URL
		if len(enclosure) == 0 && len(link) == 0 {
			continue
		}
		if len(enclosure) == 0 && len(link) > 0 {
			enclosure = link
			link = ""
		}
		torrents = append(torrents, RssTorrent{
			Id:          item.Guid,
			Title:       item.Title,
			Enclosure:   enclosure,
			Size:        stringx.ParseInt64(item.Enclosure.Length),
			Description: item.Description,
			Link:        item.Link,
			PubDate:     stringx.TimeStamp(item.PubDate),
		})
	}
	return torrents, nil
}

// MessageDetail 消息详情
func (c *npClient) MessageDetail(message Message) (string, error) {
	var md Message
	err := data(requestSiteParams{
		ctx:   c.ctx,
		site:  c.site,
		reqId: requestIdMessageDetail,
		path:  message.Link,
	}, &md, nil)
	if err != nil {
		return "", newError(c.site, err, "用户未读消息详情失败")
	}
	return md.Content, nil
}

func (c *npClient) SignIn() (SignInResult, error) {
	// 尝试获取用户基础信息，既可以判断是否需要已登录也可以用于模拟登录
	ubi, err := c.UserBasicInfo()
	if err != nil {
		return SignInResult{}, err
	}
	if !ubi.IsLogin {
		return SignInResult{
			Code:    SignInCodeNeedLogin,
			Message: "未登录",
		}, nil
	}
	if ubi.SignedIn {
		return SignInResult{
			Code:    SignInCodeSigned,
			Message: "今日已签到",
		}, nil
	}
	sc, err := SiteHelper.GetConfigByCode(c.site.Code)
	if err != nil {
		return SignInResult{}, err
	}
	// 无需签到
	if !sc.Required.SignIn {
		return SignInResult{
			Code:    SignInCodeSuccess,
			Message: "模拟登录成功",
		}, nil
	}
	// 签到
	r := signInResult{}
	statusCode := 0
	err = data(requestSiteParams{
		ctx:   c.ctx,
		site:  c.site,
		reqId: requestIdSignIn,
	}, &r, func(result siteadapt.Result) {
		statusCode = result.StatusCode
	})
	if err != nil {
		return SignInResult{}, newError(c.site, err, "签到异常")
	}
	// 签到成功
	if r.SignedIn {
		return SignInResult{
			Code:    SignInCodeSuccess,
			Message: "签到成功",
		}, nil
	}
	// 签到失败
	if statusCode == 200 {
		return SignInResult{
			Code:    SignInCodeFailure,
			Message: fmt.Sprintf("签到失败，请检查该站点是否已适配"),
		}, nil
	}
	return SignInResult{
		Code:    SignInCodeFailure,
		Message: fmt.Sprintf("签到失败，状态码：%d", statusCode),
	}, nil
}

func (c *npClient) GetDownloadUrl(torrent SearchTorrent) (string, error) {
	return torrent.Enclosure, nil
}

func (c *npClient) Details(id string) (TorrentDetail, error) {
	env := map[string]string{"id": id}
	torrentDetail := TorrentDetail{}
	err := data(requestSiteParams{
		ctx:   c.ctx,
		site:  c.site,
		reqId: requestIdDetails,
		env:   env,
	}, &torrentDetail, nil)
	if err != nil {
		return torrentDetail, newError(c.site, err, "获取详情异常")
	}
	return torrentDetail, nil
}

func (c *npClient) GetSubtitleDownloadUrl(id string) (string, error) {
	env := map[string]string{"id": id}
	subtitleUrlResult := downloadSubtitleResult{}
	err := data(requestSiteParams{
		ctx:   c.ctx,
		site:  c.site,
		reqId: requestIdGetSubtitleUrl,
		env:   env,
	}, &subtitleUrlResult, nil)
	if err != nil {
		return "", newError(c.site, err, "获取字幕下载链接异常")
	}
	return subtitleUrlResult.Url, nil
}

func (c *npClient) DownloadFile(url string) ([]byte, error) {
	if url == "" {
		return nil, nil
	}
	rd := siteadapt.RequestDefinition{
		Parser: "None",
		Method: http.MethodGet,
		Path:   url,
	}
	var data []byte
	err := raw(requestSiteParams{
		ctx:  c.ctx,
		site: c.site,
		rd:   &rd,
	}, func(result siteadapt.Result) {
		data = result.Raw
	})
	if err != nil {
		return nil, newError(c.site, err, "下载文件异常")
	}
	return data, nil
}
