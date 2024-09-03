package qywx

import (
	"fmt"
	"github.com/heibizi/go-media-kit/message"
)

func (c *Client) SendList(req message.ListMessageReq) error {
	var articles []article
	for index, item := range req.Items {
		articles = append(articles, article{
			Title:       fmt.Sprintf("%d. %s", index+1, item.Title),
			Description: item.Content,
			PicUrl:      item.Url,
			Url:         item.Link,
		})
	}
	return c.send(imageMessage{
		ToUser:  c.userId(),
		MsgType: "news",
		AgentId: c.AgentId,
		News: news{
			Articles: articles,
		},
	}, true)
}
