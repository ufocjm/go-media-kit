package moviesubject

import (
	"time"
)

type Media struct {
	Id           string  `json:"id"`
	Title        string  `json:"title"`
	Type         string  `json:"type"`
	Year         int     `json:"year"`
	Vote         float64 `json:"vote"`
	Image        string  `json:"image"`
	ImageDynamic string  `json:"imageDynamic"`
	Overview     string  `json:"overview"`
}

type Result struct {
	PageNum  int     `json:"pageNum"`
	PageSize int     `json:"pageSize"`
	Total    int64   `json:"total"`
	List     []Media `json:"list"`
}

type Subject struct {
	Code string
	Name string
}

type Category struct {
	Code   string
	Name   string
	Custom bool
}

type ItemsRequest struct {
	Category string
	Subject  string
}

type TmdbApiParams struct {
	CustomBaseURL   string        // 自定义 api 地址，默认为 api.themoviedb.org
	ImageBaseURL    string        // 图片地址，默认为 image.tmdb.org
	ImageWidth      string        // 图片宽度，默认为 w500
	ApiKey          string        // api key
	Language        string        // 语言，默认为 zh
	Region          string        // 地区，默认为 CN
	IncludeAdult    bool          // 是否包含成人内容，默认为 false
	Timeout         time.Duration // 超时时间
	MaxIdleConn     int           // 最大空闲连接数
	IdleConnTimeout time.Duration // 空闲连接超时时间
	Proxy           string        // 代理
}
