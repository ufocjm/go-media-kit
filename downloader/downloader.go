package downloader

import (
	"github.com/heibizi/go-media-kit/downloader/qb"
)

type Client interface {
	Login() error
	Version() (string, error)
	AddTorrents(req any) error
}

type ClientType string

const (
	QbClientType ClientType = "qb"
)

var clientFactoryRegistry = map[ClientType]func(config any) Client{
	QbClientType: func(config any) Client {
		return qb.NewClient(config.(qb.Config))
	},
}

func NewClient(clientType ClientType, config any) Client {
	return clientFactoryRegistry[clientType](config)
}
