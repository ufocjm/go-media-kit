package douban

import (
	"fmt"
	"net/url"
	"strconv"
)

func (c *ApiClient) SubjectCollectionItems(id string, start int, count int) (SubjectCollectionItemsResult, error) {
	u, _ := url.JoinPath(apiUrl, "/subject_collection/", id, "/items")
	var o SubjectCollectionItemsResult
	params := map[string]string{"start": strconv.Itoa(start), "count": strconv.Itoa(count)}
	err := c.get(u, params, &o)
	if err != nil {
		return o, fmt.Errorf("获取 subject collection items 失败: %v", err)
	}
	return o, nil
}

type SubjectCollectionItemsResult struct {
	ApiResult              `mapstructure:",squash"`
	SubjectCollectionItems []SubjectCollectionItem `json:"subjectCollectionItems" mapstructure:"subject_collection_items"`
	SubjectCollection      SubjectCollectionResult `json:"subjectCollection" mapstructure:"subject_collection"`
}

type HonorInfo struct {
	Kind  string `json:"kind" mapstructure:"kind"`
	URI   string `json:"uri" mapstructure:"uri"`
	Rank  int    `json:"rank" mapstructure:"rank"`
	Title string `json:"title" mapstructure:"title"`
}

type SubjectCollectionItem struct {
	Comment           string      `json:"comment" mapstructure:"comment"`
	Rating            Rating      `json:"rating" mapstructure:"rating"`
	ControversyReason string      `json:"controversyReason" mapstructure:"controversy_reason"`
	Pic               Pic         `json:"pic" mapstructure:"pic"`
	Rank              int         `json:"rank" mapstructure:"rank"`
	URI               string      `json:"uri" mapstructure:"uri"`
	IsShow            bool        `json:"isShow" mapstructure:"is_show"`
	VendorIcons       []string    `json:"vendorIcons" mapstructure:"vendor_icons"`
	CardSubtitle      string      `json:"cardSubtitle" mapstructure:"card_subtitle"`
	Id                string      `json:"id" mapstructure:"id"`
	Title             string      `json:"title" mapstructure:"title"`
	HasLineWatch      bool        `json:"hasLineWatch" mapstructure:"has_linewatch"`
	IsReleased        bool        `json:"isReleased" mapstructure:"is_released"`
	Interest          interface{} `json:"interest" mapstructure:"interest"`
	ColorScheme       ColorScheme `json:"colorScheme" mapstructure:"color_scheme"`
	Type              string      `json:"type" mapstructure:"type"`
	Description       string      `json:"description" mapstructure:"description"`
	//Tags              map[string]string `json:"tags" mapstructure:"tags"`
	CoverURL         string      `json:"coverUrl" mapstructure:"cover_url"`
	Photos           []string    `json:"photos" mapstructure:"photos"`
	Actions          []string    `json:"actions" mapstructure:"actions"`
	SharingURL       string      `json:"sharingUrl" mapstructure:"sharing_url"`
	URL              string      `json:"url" mapstructure:"url"`
	HonorInfos       []HonorInfo `json:"honorInfos" mapstructure:"honor_infos"`
	GoodRatingStats  int         `json:"goodRatingStats" mapstructure:"good_rating_stats"`
	Subtype          string      `json:"subtype" mapstructure:"subtype"`
	NullRatingReason string      `json:"nullRatingReason" mapstructure:"null_rating_reason"`
}
