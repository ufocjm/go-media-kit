package qywx

import (
	"fmt"
	"github.com/heibizi/go-media-kit/message"
	"strings"
)

func (c *Client) SendText(req message.TextMessageReq) error {
	var content string
	if req.Content != "" {
		content = fmt.Sprintf("%s\n%s", req.Title, strings.ReplaceAll(req.Content, "\n\n", "\n"))
	} else {
		content = req.Title
	}
	if req.Link != "" {
		content = fmt.Sprintf("%s\n\n<a href='%s'>查看详情</a>", content, req.Link)
	}
	return c.send(messageReq{
		ToUser:               c.userId(),
		MsgType:              "text",
		AgentID:              c.AgentId,
		Text:                 textMessage{content},
		Safe:                 0,
		EnableIDTrans:        0,
		EnableDuplicateCheck: 0,
	}, true)
}
