package qywx

import "github.com/heibizi/go-media-kit/message"

func (c *Client) SendImage(req message.ImageMessageReq) error {
	if req.Url == "" {
		return c.SendText(req.TextMessageReq)
	}
	return c.send(imageMessage{
		ToUser:  c.userId(),
		MsgType: "news",
		AgentId: c.AgentId,
		News: news{
			Articles: []article{
				{
					Title:       req.Title,
					Description: req.Content,
					PicUrl:      req.Url,
					Url:         req.Link,
				},
			},
		},
	}, true)
}
