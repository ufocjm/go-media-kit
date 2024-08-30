package btsite

import (
	"context"
	"fmt"
	"github.com/heibizi/go-media-kit/core/utils/netx"
	"github.com/heibizi/go-media-kit/siteadapt"
	"strconv"
	"strings"
	"time"
)

type (
	// mtClient 馒头客户端
	mtClient struct {
		ctx  context.Context
		site *Site
		*npClient
	}
	mtMyPeerStatus struct {
		Leecher int `mapstructure:"leecher"`
		Seeder  int `mapstructure:"seeder"`
	}
	mtMsgNotifyStatistic struct {
		Count  int `mapstructure:"count"`
		UnMake int `mapstructure:"un_make"`
	}
	mtProfile struct {
		CreatedDate      int64   `mapstructure:"created_date,omitempty"`
		LastModifiedDate int64   `mapstructure:"last_modified_date,omitempty"`
		Username         string  `mapstructure:"username,omitempty"`
		Uploaded         int64   `mapstructure:"uploaded,omitempty"`
		Downloaded       int64   `mapstructure:"downloaded,omitempty"`
		ShareRate        float64 `mapstructure:"share_rate,omitempty"`
		Bonus            float64 `mapstructure:"bonus,omitempty"`
		Role             string  `mapstructure:"role,omitempty"`
	}
	mtUserTorrent struct {
		Size int64 `mapstructure:"size"`
	}
	mtSysRole struct {
		Id      string `mapstructure:"id"`
		NameChs string `mapstructure:"name_chs"`
		NameEng string `mapstructure:"name_eng"`
		Image   string `mapstructure:"image"`
	}
)

const (
	requestIdMTMyPeerStatus       requestId = "my_peer_status"
	requestIdMTMsgNotifyStatistic requestId = "msg_notify_statistic"
	requestIdMTProfile            requestId = "profile"
	requestIdMTUserTorrentList    requestId = "user_torrent_list"
	requestIdMTSysRoleList        requestId = "sys_role_list"
	requestIdMTGenDLToken         requestId = "gen_dl_token"
)

func (c *mtClient) UserBasicInfo() (UserBasicInfo, error) {
	mp, err := c.memberProfile()
	if err != nil {
		return UserBasicInfo{}, err
	}
	return UserBasicInfo{
		IsLogin:    len(mp.Username) > 0,
		Id:         c.site.UserId,
		Name:       mp.Username,
		Ratio:      mp.ShareRate,
		Uploaded:   mp.Uploaded,
		Downloaded: mp.Downloaded,
		Bonus:      mp.Bonus,
	}, nil
}

func (c *mtClient) UserDetails() (UserDetails, error) {
	mp, err := c.memberProfile()
	if err != nil {
		return UserDetails{}, err
	}
	level := ""
	levelIcon := ""
	srl, err := c.sysRoleList()
	if err != nil {
		return UserDetails{}, err
	}
	if srl != nil {
		for _, role := range srl {
			if mp.Role == role.Id {
				level = role.NameChs + " " + role.NameEng
				levelIcon = role.Image
			}
		}
	}
	return UserDetails{
		Level:        level,
		LevelIcon:    levelIcon,
		JoinAt:       mp.CreatedDate,
		LastAccessed: mp.LastModifiedDate,
	}, nil
}

func (c *mtClient) Search(searchParams SearchParams) ([]SearchTorrent, error) {
	sh := SiteHelper
	sc, err := sh.GetConfigByCode(c.site.Code)
	if err != nil {
		return nil, newError(c.site, err, "未获取到站点配置")
	}
	if searchParams.PageSize == 0 {
		searchParams.PageSize = 100
	}
	body := map[string]any{
		"mode":    "normal",
		"visible": 1,
		// 馒头从 1 开始
		"pageNumber": searchParams.PageNum + 1,
		"pageSize":   searchParams.PageSize,
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
			var categories []string
			for _, cat := range sc.Categories.Movie {
				categories = append(categories, cat.Id)
			}
			body["categories"] = categories
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
		var labels []string
		for _, label := range torrent.Labels {
			for _, s := range strings.Split(label, "|") {
				labels = append(labels, s)
			}
		}
		torrent.PageUrl = pageUrl
		torrent.Labels = labels
		torrents = append(torrents, torrent)
	}
	return torrents, nil
}

func (c *mtClient) SeedingStatistics() (SeedingStatistics, error) {
	seeding := SeedingStatistics{}
	for pageNumber := 1; ; pageNumber++ {
		tl, err := c.userTorrentList(pageNumber)
		if err != nil {
			return seeding, err
		}
		if tl == nil {
			break
		}
		for _, ut := range tl {
			seeding.Count = seeding.Count + 1
			seeding.Size = seeding.Size + ut.Size
		}
		time.Sleep(500 * time.Millisecond)
	}
	return seeding, nil
}

func (c *mtClient) MyHr() ([]HrTorrent, error) {
	// 无 hr
	return nil, nil
}

func (c *mtClient) SignIn() (SignInResult, error) {
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
	return SignInResult{
		Code:    SignInCodeSuccess,
		Message: "模拟登录成功",
	}, nil
}

type mtMessage struct {
	Message
	UnRead bool `mapstructure:"unread"`
}
type mtNotice struct {
	Title   string `mapstructure:"title,omitempty"`   // 标题
	Date    string `mapstructure:"date,omitempty"`    // 日期
	Content string `mapstructure:"content,omitempty"` // 内容
}

func (c *mtClient) Notice() (string, error) {
	var ns []mtNotice
	err := list(requestSiteParams{
		ctx:      c.ctx,
		site:     c.site,
		reqId:    requestIdNotice,
		siteData: c.homeData,
	}, &ns, nil)
	if err != nil {
		return "", newError(c.site, err, "解析公告失败")
	}
	var notices []string
	for _, n := range ns {
		notices = append(notices, fmt.Sprintf("%s(%s)\n%s", n.Title, n.Date, n.Content))
	}
	if len(notices) > 0 {
		return strings.Join(notices, "\n\n"), nil
	}
	return "", nil
}

func (c *mtClient) Messages(page int) ([]Message, error) {
	var o []mtMessage
	fromDataEnv := map[string]string{"pageNumber": strconv.Itoa(page)}
	err := list(requestSiteParams{
		ctx:   c.ctx,
		site:  c.site,
		reqId: requestIdMessages,
		env:   fromDataEnv,
	}, &o, nil)
	if err != nil {
		return nil, newError(c.site, err, "获取消息异常")
	}
	var ids []string
	var list []Message
	for _, message := range o {
		if message.UnRead {
			ids = append(ids, message.Id)
		}
		list = append(list, message.Message)
	}
	if len(ids) > 0 {
		err = c.markAsRead(ids)
		if err != nil {
			return nil, newError(c.site, err, "标记为已读异常")
		}
	}
	return list, nil
}

// MessageDetail 消息详情
func (c *mtClient) MessageDetail(message Message) (string, error) {
	return message.Content, nil
}

func (c *mtClient) GetDownloadUrl(torrent SearchTorrent) (string, error) {
	formDataEnv := map[string]string{"id": torrent.Id}
	m := make(map[string]any)
	err := data(requestSiteParams{
		ctx:   c.ctx,
		site:  c.site,
		reqId: requestIdMTGenDLToken,
		env:   formDataEnv,
	}, &m, nil)
	if err != nil {
		return "", err
	}
	return m["url"].(string), nil
}

func (c *mtClient) GetSubtitleDownloadUrl(_ string) (string, error) {
	return "", nil
}

func (c *mtClient) myPeerStatus() (mtMyPeerStatus, error) {
	o := mtMyPeerStatus{}
	err := data(requestSiteParams{
		ctx:   c.ctx,
		site:  c.site,
		reqId: requestIdMTMyPeerStatus,
	}, &o, nil)
	if err != nil {
		return o, newError(c.site, err, "做种信息异常")
	}
	return o, nil
}

func (c *mtClient) msgNotifyStatistic() (mtMsgNotifyStatistic, error) {
	o := mtMsgNotifyStatistic{}
	err := data(requestSiteParams{
		ctx:   c.ctx,
		site:  c.site,
		reqId: requestIdMTMsgNotifyStatistic,
	}, &o, nil)
	if err != nil {
		return o, newError(c.site, err, "消息通知统计异常")
	}
	return o, nil
}

// memberProfile 获取个人资料
func (c *mtClient) memberProfile() (mtProfile, error) {
	o := mtProfile{}
	err := data(requestSiteParams{
		ctx:   c.ctx,
		site:  c.site,
		reqId: requestIdMTProfile,
	}, &o, nil)
	if err != nil {
		return o, newError(c.site, err, "个人资料异常")
	}
	return o, nil
}

func (c *mtClient) userTorrentList(pageNumber int) ([]mtUserTorrent, error) {
	var o []mtUserTorrent
	var body = make(map[string]any)
	body["userid"] = c.site.UserId
	body["type"] = "SEEDING"
	body["pageNumber"] = pageNumber
	body["pageSize"] = 100
	err := list(requestSiteParams{
		ctx:   c.ctx,
		site:  c.site,
		reqId: requestIdMTUserTorrentList,
		body:  body,
	}, &o, nil)
	if err != nil {
		return nil, newError(c.site, err, "用户做种列表异常")
	}
	return o, nil
}

func (c *mtClient) sysRoleList() ([]mtSysRole, error) {
	var o []mtSysRole
	err := list(requestSiteParams{
		ctx:   c.ctx,
		site:  c.site,
		reqId: requestIdMTSysRoleList,
	}, &o, nil)
	if err != nil {
		return nil, newError(c.site, err, "角色列表异常")
	}
	return o, nil
}

// markAsRead 未读消息设为已读
func (c *mtClient) markAsRead(ids []string) error {
	r := markAsReadResult{}
	err := data(requestSiteParams{
		ctx:   c.ctx,
		site:  c.site,
		reqId: requestIdMarkAsRead,
		env:   map[string]string{"ids": strings.Join(ids, ",")},
	}, &r, nil)
	if err != nil {
		return newError(c.site, err, "未读消息设为已读异常")
	}
	if r.Success {
		return nil
	}
	return newError(c.site, err, "未读消息设为已读失败: %s", r.Message)
}
