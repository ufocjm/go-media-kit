package moviesubject

const (
	MediaTypeMovie = "movie"
	MediaTypeTv    = "tv"
)

var (
	CategoryDoubanRanks             = Category{Code: "douban_ranks", Name: "豆瓣榜单"}
	CategoryDoubanYearRanks         = Category{Code: "douban_year_ranks", Name: "豆瓣2023年度榜单"}
	CategoryDoubanDouList           = Category{Code: "douban_dou_list", Name: "豆瓣豆列", Custom: true}
	CategoryDoubanSubjectCollection = Category{Code: "douban_subject_collection", Name: "豆瓣片单", Custom: true}
	CategoryTmdb                    = Category{Code: "tmdb", Name: "TMDB"}
	CategoryMg                      = Category{Code: "mg", Name: "芒果"}
	CategoryIqy                     = Category{Code: "iqy", Name: "爱奇艺"}
	CategoryYk                      = Category{Code: "yk", Name: "优酷"}
	CategoryVqq                     = Category{Code: "vqq", Name: "腾讯视频"}
	CategoryDisney                  = Category{Code: "disney", Name: "迪士尼"}
	CategoryNetflix                 = Category{Code: "netflix", Name: "网飞"}
)

var (
	DoubanRanksMovieShowing                 = Subject{"movie_showing", "影院热映"}
	DoubanRanksMovieHotGaia                 = Subject{"movie_hot_gaia", "豆瓣热门"}
	DoubanRanksMovieComingSoonDomestic      = Subject{"movie_coming_soon_domestic", "国内即将上映"}
	DoubanRanksMovieComingSoonInternational = Subject{"movie_coming_soon_international", "全球值得期待"}
	DoubanRanksMovieRealTimeHot             = Subject{"movie_real_time_hotest", "实时热门电影"}
	DoubanRanksMovieWeeklyBest              = Subject{"movie_weekly_best", "一周口碑电影榜"}
	DoubanRanksMovieTop250                  = Subject{"movie_top250", "豆瓣电影 Top250"}
	DoubanRanksMovieHotTop20                = Subject{"ECPE465QY", "热门电影Top20"}
	DoubanRanksMovieTopRatedTop20           = Subject{"EC7Q5H2QI", "高分电影Top20"}
	DoubanRanksMovieDarkHorseTop20          = Subject{"ECSU5CIVQ", "冷门佳作Top20"}
	DoubanRanksTvComingSoon                 = Subject{"tv_coming_soon", "即将播出"}
	DoubanRanksTvHot                        = Subject{"tv_hot", "近期热门剧集"}
	DoubanRanksTvDomestic                   = Subject{"tv_domestic", "近期热门国产剧"}
	DoubanRanksTvAmerican                   = Subject{"tv_american", "近期热门美剧"}
	DoubanRanksTvJapanese                   = Subject{"tv_japanese", "近期热门日剧"}
	DoubanRanksTvKorean                     = Subject{"tv_korean", "近期热门韩剧"}
	DoubanRanksTvAnimation                  = Subject{"tv_animation", "近期热门动画"}
	DoubanRanksTvShowHot                    = Subject{"show_hot", "近期热门综艺节目"}
	DoubanRanksTvShowDomestic               = Subject{"show_domestic", "近期热门国内综艺"}
	DoubanRanksTvShowForeign                = Subject{"show_foreign", "近期热门国外综艺"}
	DoubanRanksTvRealTimeHot                = Subject{"tv_real_time_hotest", "实时热门电视"}
	DoubanRanksTvChineseBestWeekly          = Subject{"tv_chinese_best_weekly", "华语口碑剧集榜"}
	DoubanRanksTvGlobalBestWeekly           = Subject{"tv_global_best_weekly", "全球口碑剧集榜"}
	DoubanRanksTvShowChineseBestWeekly      = Subject{"show_chinese_best_weekly", "国内口碑综艺榜"}
	DoubanRanksTvShowGlobalBestWeekly       = Subject{"show_global_best_weekly", "国外口碑综艺榜"}
)

var (
	DoubanYearRanksMovieDomestic    = Subject{"ECQ46F7XI", "华语电影"}
	DoubanYearRanksMovieForeign     = Subject{"ECFA6FLWQ", "外语电影"}
	DoubanYearRanksMovieDarkHorse   = Subject{"ECMY6GCCA", "冷门佳作"}
	DoubanYearRanksMovieJapanese    = Subject{"ECCU6MRTY", "日本电影"}
	DoubanYearRanksMovieBz          = Subject{"EC4Y6ALRA", "韩国电影"}
	DoubanYearRanksMovieComedy      = Subject{"ECCI6H3TA", "喜剧片"}
	DoubanYearRanksMovieLove        = Subject{"EC3A56FJA", "爱情片"}
	DoubanYearRanksMovieAnimation   = Subject{"ECHU6BXBI", "动画片"}
	DoubanYearRanksMovieDocumentary = Subject{"ECRM6A2JA", "纪录片"}
	DoubanYearRanksMovieHorrible    = Subject{"ECYI6DWVQ", "恐怖片"}
	DoubanYearRanksTvDomestic       = Subject{"ECTE6EOZA", "华语剧集"}
	DoubanYearRanksTvForeign        = Subject{"ECUI6CVAI", "英美新剧"}
	DoubanYearRanksTvDarkHorse      = Subject{"ECM46I42A", "英美续订剧"}
	DoubanYearRanksTvBz             = Subject{"EC246FT6Y", "韩国剧集"}
	DoubanYearRanksTvJapanese       = Subject{"ECPE6B6NI", "日本剧集"}
	DoubanYearRanksTvShow           = Subject{"EC7I6GR6A", "综艺"}
	DoubanYearRanksTvAnimation      = Subject{"EC3Q6JTOQ", "动画剧集"}
	DoubanYearRanksTvDocumentary    = Subject{"ECCU6NNCI", "纪录剧集"}
)

var (
	TmdbMovieTrendingWeek = Subject{"movie_trending_week", "电影-本周趋势"}
	TmdbMovieTrendingDay  = Subject{"movie_trending_day", "电影-今日趋势"}
	TmdbMoviePopular      = Subject{"movie_popular", "电影-热门"}
	TmdbMovieNowPlaying   = Subject{"movie_now_playing", "电影-正在上映"}
	TmdbMovieUpcoming     = Subject{"movie_upcoming", "电影-即将上映"}
	TmdbMovieTopRated     = Subject{"movie_top_rated", "电影-高分"}
	TmdbTvTrendingWeek    = Subject{"tv_trending_week", "剧集-本周趋势"}
	TmdbTvTrendingDay     = Subject{"tv_trending_day", "剧集-今日趋势"}
	TmdbTvPopular         = Subject{"tv_popular", "剧集-热门"}
	TmdbTvAiringToday     = Subject{"tv_airing_today", "剧集-今日播出"}
	TmdbTvOnTheAir        = Subject{"tv_on_the_air", "剧集-电视播出中"}
	TmdbTvTopRated        = Subject{"tv_top_rated", "剧集-高分"}
)

var (
	VqqMoviePopular    = Subject{"movie_popular", "电影-最热"}
	VqqMovieNowPlaying = Subject{"movie_now_playing", "电影-最新"}
	VqqMovieTopRated   = Subject{"movie_top_rated", "电影-高分好评"}
	VqqTvPopular       = Subject{"tv_popular", "电视剧-最热"}
	VqqTvNowPlaying    = Subject{"tv_now_playing", "电视剧-最新上架"}
	VqqTvTopRated      = Subject{"tv_top_rated", "电视剧-好评"}
)

var (
	IqyMovieComprehensive = Subject{"movie_comprehensive", "电影-综合"}
	IqyMoviePopular       = Subject{"movie_popular", "电影-最热"}
	IqyMovieNowPlaying    = Subject{"movie_now_playing", "电影-最新"}
	IqyMovieTopRated      = Subject{"movie_top_rated", "电影-高分"}
	IqyTvComprehensive    = Subject{"tv_comprehensive", "电视剧-综合"}
	IqyTvPopular          = Subject{"tv_popular", "电视剧-最热"}
	IqyTvNowPlaying       = Subject{"tv_now_playing", "电视剧-最新"}
	IqyTvTopRated         = Subject{"tv_top_rated", "电视剧-高分"}
	IqyVipMovie           = Subject{"vip_movie", "VIP-电影"}
	IqyVipTv              = Subject{"vip_tv", "VIP-电视剧"}
	IqyVipVarietyShow     = Subject{"vip_variety_show", "VIP-综艺"}
	IqyVipAnimation       = Subject{"vip_animation", "VIP-动漫"}
	IqyVipDocumentary     = Subject{"vip_documentary", "VIP-纪录片"}
	IqyVipChild           = Subject{"vip_child", "VIP-儿童"}
)

var (
	YkMovieComprehensive = Subject{"movie_comprehensive", "电影-综合"}
	YkMoviePopular       = Subject{"movie_popular", "电影-热度最高"}
	YkMovieNewly         = Subject{"movie_newly", "电影-最新上线"}
	YkMovieTopRated      = Subject{"movie_top_rated", "电影-最好评"}
	YkMovieMostPlayed    = Subject{"movie_most_played", "电影-最多播放"}
	YkTvComprehensive    = Subject{"tv_comprehensive", "电视剧-综合"}
	YkTvPopular          = Subject{"tv_popular", "电视剧-热度最高"}
	YkTvNewly            = Subject{"tv_newly", "电视剧-最新上线"}
	YkTvTopRated         = Subject{"tv_top_rated", "电视剧-最好评"}
	YkTvMostPlayed       = Subject{"tv_most_played", "电视剧-最多播放"}
	YkVipMovie           = Subject{"vip_movie", "VIP-电影"}
	YkVipTv              = Subject{"vip_tv", "VIP-电视剧"}
	YkVipVarietyShow     = Subject{"vip_variety_show", "VIP-综艺"}
	YkVipAnimation       = Subject{"vip_animation", "VIP-动漫"}
	YkVipDocumentary     = Subject{"vip_documentary", "VIP-纪录片"}
	YkVipChild           = Subject{"vip_child", "VIP-儿童"}
)

var (
	MgTvVarietyPopular = Subject{"tv_variety_popular", "综艺-最热"}
	MgTvVarietyNewly   = Subject{"tv_variety_newly", "综艺-最新"}
	MgTvPopular        = Subject{"tv_popular", "电视剧-最热"}
	MgTvNewly          = Subject{"tv_newly", "电视剧-最新"}
	MgMoviePopular     = Subject{"movie_popular", "电影-最热"}
	MgMovieNewly       = Subject{"movie_newly", "电影-最新"}
	MgTvChild          = Subject{"tv_child", "少儿"}
	MgTvAnimation      = Subject{"tv_animation", "动漫"}
	MgTvDocument       = Subject{"tv_documentary", "纪录片"}
	MgTvEducation      = Subject{"tv_education", "教育"}
	MgVipVariety       = Subject{"vip_variety", "VIP-综艺"}
	MgVipTv            = Subject{"vip_tv", "VIP-电视剧"}
	MgVipMovie         = Subject{"vip_movie", "VIP-电影"}
)

type SubjectInfo struct {
	Category Category
	Subjects []Subject
}

var subjectInfos = []SubjectInfo{
	{
		Category: CategoryDoubanRanks,
		Subjects: []Subject{
			DoubanRanksMovieShowing,
			DoubanRanksMovieHotGaia,
			DoubanRanksMovieComingSoonDomestic,
			DoubanRanksMovieComingSoonInternational,
			DoubanRanksMovieRealTimeHot,
			DoubanRanksMovieWeeklyBest,
			DoubanRanksMovieTop250,
			DoubanRanksMovieHotTop20,
			DoubanRanksMovieTopRatedTop20,
			DoubanRanksMovieDarkHorseTop20,
			DoubanRanksTvComingSoon,
			DoubanRanksTvHot,
			DoubanRanksTvDomestic,
			DoubanRanksTvAmerican,
			DoubanRanksTvJapanese,
			DoubanRanksTvKorean,
			DoubanRanksTvAnimation,
			DoubanRanksTvShowHot,
			DoubanRanksTvShowDomestic,
			DoubanRanksTvShowForeign,
			DoubanRanksTvRealTimeHot,
			DoubanRanksTvChineseBestWeekly,
			DoubanRanksTvGlobalBestWeekly,
			DoubanRanksTvShowChineseBestWeekly,
			DoubanRanksTvShowGlobalBestWeekly,
		},
	},
	{
		Category: CategoryDoubanYearRanks,
		Subjects: []Subject{
			DoubanYearRanksMovieDomestic,
			DoubanYearRanksMovieForeign,
			DoubanYearRanksMovieDarkHorse,
			DoubanYearRanksMovieJapanese,
			DoubanYearRanksMovieBz,
			DoubanYearRanksMovieComedy,
			DoubanYearRanksMovieLove,
			DoubanYearRanksMovieAnimation,
			DoubanYearRanksMovieDocumentary,
			DoubanYearRanksMovieHorrible,
			DoubanYearRanksTvDomestic,
			DoubanYearRanksTvForeign,
			DoubanYearRanksTvDarkHorse,
			DoubanYearRanksTvBz,
			DoubanYearRanksTvJapanese,
			DoubanYearRanksTvShow,
			DoubanYearRanksTvAnimation,
			DoubanYearRanksTvDocumentary,
		},
	},
	{
		Category: CategoryDoubanDouList,
		Subjects: []Subject{},
	},
	{
		Category: CategoryDoubanSubjectCollection,
		Subjects: []Subject{},
	},
	{
		Category: CategoryTmdb,
		Subjects: []Subject{
			TmdbMovieTrendingWeek,
			TmdbMovieTrendingDay,
			TmdbTvTrendingWeek,
			TmdbTvTrendingDay,
			TmdbMoviePopular,
			TmdbMovieNowPlaying,
			TmdbMovieUpcoming,
			TmdbMovieTopRated,
			TmdbTvPopular,
			TmdbTvAiringToday,
			TmdbTvOnTheAir,
			TmdbTvTopRated,
		},
	},
	{
		Category: CategoryMg,
		Subjects: []Subject{
			MgTvVarietyPopular,
			MgTvVarietyNewly,
			MgTvPopular,
			MgTvNewly,
			MgMoviePopular,
			MgMovieNewly,
			MgTvChild,
			MgTvAnimation,
			MgTvDocument,
			MgTvEducation,
			MgVipVariety,
			MgVipTv,
			MgVipMovie,
		},
	},
	{
		Category: CategoryIqy,
		Subjects: []Subject{
			IqyMovieComprehensive,
			IqyMoviePopular,
			IqyMovieNowPlaying,
			IqyMovieTopRated,
			IqyTvComprehensive,
			IqyTvPopular,
			IqyTvNowPlaying,
			IqyTvTopRated,
			IqyVipMovie,
			IqyVipTv,
			IqyVipVarietyShow,
			IqyVipAnimation,
			IqyVipDocumentary,
			IqyVipChild,
		},
	},
	{
		Category: CategoryYk,
		Subjects: []Subject{
			YkMovieComprehensive,
			YkMoviePopular,
			YkMovieNewly,
			YkMovieTopRated,
			YkMovieMostPlayed,
			YkTvComprehensive,
			YkTvPopular,
			YkTvNewly,
			YkTvTopRated,
			YkTvMostPlayed,
			YkVipMovie,
			YkVipTv,
			YkVipVarietyShow,
			YkVipAnimation,
			YkVipDocumentary,
			YkVipChild,
		},
	},
	{
		Category: CategoryVqq,
		Subjects: []Subject{
			VqqMoviePopular,
			VqqMovieNowPlaying,
			VqqMovieTopRated,
			VqqTvPopular,
			VqqTvNowPlaying,
			VqqTvTopRated,
		},
	},
}
