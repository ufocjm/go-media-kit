package iyuu

import (
	"fmt"
	"github.com/heibizi/go-media-kit/message"
)

func (c *Client) SendList(req message.ListMessageReq) error {
	return fmt.Errorf("IYUU暂不支持列表消息发送")
}
