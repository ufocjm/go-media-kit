package iyuu

import (
	"github.com/heibizi/go-media-kit/message"
)

func (c *Client) SendText(req message.TextMessageReq) error {
	return c.send(messageReq{
		Title:   req.Title,
		Content: req.Content,
	})
}
