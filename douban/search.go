package douban

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// Search 搜索 同时支持 movie tv 搜索，不需要指定 mediaType
func (c *ApiClient) Search(q string, start int, count int) (SearchResult, error) {
	requestUrl, _ := url.JoinPath(apiUrl, "/search/movie")
	var o SearchResult
	params := map[string]string{"q": q, "start": strconv.Itoa(start), "count": strconv.Itoa(count)}
	err := c.get(requestUrl, params, &o)
	if err != nil {
		return o, fmt.Errorf("search 异常")
	}
	items := make([]SearchItem, 0)
	for _, item := range o.Items {
		pic, err := resizePic(item.Target.CoverURL, 381)
		if err != nil {
			return o, err
		}
		item.Target.CoverURL = pic
		items = append(items, item)
	}
	o.Items = items
	return o, nil
}

type SearchResult struct {
	ApiResult `mapstructure:",squash"`
	Banned    string       `json:"banned" mapstructure:"banned"`
	Items     []SearchItem `json:"items" mapstructure:"items"`
}

func resizePic(url string, minHeight int) (string, error) {
	re := regexp.MustCompile(`/h/(\d+)`)
	matches := re.FindStringSubmatch(url)
	if len(matches) == 0 {
		return url, fmt.Errorf("no height parameter found in Url")
	}
	height, err := strconv.Atoi(matches[1])
	if err != nil {
		return url, fmt.Errorf("invalid height value: %v", err)
	}
	if height < minHeight {
		url = strings.Replace(url, fmt.Sprintf("/h/%d", height), fmt.Sprintf("/h/%d", minHeight), 1)
	}
	return url, nil
}
