package btsite

type requestId string

const (
	requestIdUserBasicInfo     requestId = "user_basic_info"
	requestIdUserDetails       requestId = "user_details"
	requestIdSearch            requestId = "search"
	requestIdSeedingStatistics requestId = "seeding_statistics"
	requestIdFavicon           requestId = "favicon"
	requestIdMyHr              requestId = "my_hr"
	requestIdMessages          requestId = "messages"
	requestIdMessageDetail     requestId = "message_detail"
	requestIdMarkAsRead        requestId = "mark_as_read"
	requestIdNotice            requestId = "notice"
	requestIdSignIn            requestId = "sign_in"
	requestIdDetails           requestId = "details"
	requestIdGetSubtitleUrl    requestId = "get_subtitle_url"
	requestIdDownloadSubtitle  requestId = "download_subtitle"
	requestIdDownloadTorrent   requestId = "download_torrent"
)

type siteSchema string

const (
	siteSchemaNexusPHP siteSchema = "NexusPHP"
	siteSchemaMTorrent siteSchema = "mTorrent"
	siteSchemaZhuQue   siteSchema = "zhuQue"
)

type SignInCode int

const (
	SignInCodeSuccess   SignInCode = 0 // 签到成功
	SignInCodeSigned    SignInCode = 1 // 已经签到过
	SignInCodeFailure   SignInCode = 2 // 签到失败
	SignInCodeNeedLogin SignInCode = 3 // 未登录
)

var Movie = MediaType{
	Code: "movie",
	Name: "电影",
}

var Tv = MediaType{
	Code: "tv",
	Name: "电视剧",
}

var Anime = MediaType{
	Code: "anime",
	Name: "动漫",
}

var MediaTypes = []MediaType{
	Movie,
	Tv,
	Anime,
}
