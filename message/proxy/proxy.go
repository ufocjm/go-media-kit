package proxy

import (
	"errors"
	"github.com/heibizi/go-media-kit/message"
	"github.com/heibizi/go-media-kit/message/qywx"
)

type ClientProxy struct {
	client message.Client
}

var clientFactoryRegistry = map[message.ClientType]func(config any) message.Client{
	message.WorkWechat: func(config any) message.Client {
		return qywx.NewClient(config.(qywx.Config))
	},
}

func NewClientProxy(clientType message.ClientType, config any) *ClientProxy {
	return &ClientProxy{
		client: clientFactoryRegistry[clientType](config),
	}
}

func (p *ClientProxy) Send(req any) error {
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
