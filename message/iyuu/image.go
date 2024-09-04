package iyuu

import "github.com/heibizi/go-media-kit/message"

func (c *Client) SendImage(req message.ImageMessageReq) error {
	return c.SendText(message.TextMessageReq{Title: req.Title, Content: req.Content})
}
