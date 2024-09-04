package serverchan

import (
	"fmt"
	"github.com/heibizi/go-media-kit/message"
	"strings"
)

func (c *Client) SendList(req message.ListMessageReq) error {
	if len(req.Items) == 0 {
		return fmt.Errorf("消息体不能为空")
	}
	title := req.Items[0].Title
	contents := []string{req.Items[0].Content}
	if len(req.Items) > 1 {
		for _, item := range req.Items[1:] {
			contents = append(contents, item.Title)
			contents = append(contents, item.Content)
		}
	}
	content := strings.Join(contents, "\n\n")
	return c.SendText(message.TextMessageReq{Title: title, Content: content})
}
