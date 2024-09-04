package proxy

import (
	"errors"
	"github.com/heibizi/go-media-kit/message"
	"github.com/heibizi/go-media-kit/message/iyuu"
	"github.com/heibizi/go-media-kit/message/qywx"
)

type ClientProxy struct {
	client message.Client
}

func NewClientProxy(config any) *ClientProxy {
	var client message.Client
	switch t := config.(type) {
	case qywx.Config:
		client = qywx.NewClient(t)
	case iyuu.Config:
		client = iyuu.NewClient(t)
	}

	return &ClientProxy{
		client: client,
	}
}

func (p *ClientProxy) Send(req any) error {
	if p.client == nil {
		return errors.New("invalid config type")
	}
	switch m := req.(type) {
	case message.TextMessageReq:
		return p.client.SendText(m)
	case message.ImageMessageReq:
		return p.client.SendImage(m)
	case message.ListMessageReq:
		return p.client.SendList(m)
	default:
		return errors.New("unsupported message type")
	}
}
