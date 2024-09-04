package message_test

import (
	"context"
	"github.com/heibizi/go-media-kit/message"
	"github.com/heibizi/go-media-kit/message/proxy"
	"github.com/heibizi/go-media-kit/message/qywx"
	"github.com/redis/go-redis/v9"
	"os"
	"testing"
)

var clientProxy *proxy.ClientProxy

func TestMain(m *testing.M) {
	clientProxy = proxy.NewClientProxy(qywx.Config{
		Ctx:         context.Background(),
		CorpId:      os.Getenv("GO_QYWX_CORP_ID"),
		CorpSecret:  os.Getenv("GO_QYWX_CORP_SECRET"),
		AgentId:     os.Getenv("GO_QYWX_AGENT_ID"),
		Proxy:       os.Getenv("GO_QYWX_PROXY"),
		UserId:      "",
		RedisClient: redis.NewClient(&redis.Options{Addr: os.Getenv("GO_QYWX_REDIS_ADDR")}),
	})
	m.Run()
}

func TestText(t *testing.T) {
	err := clientProxy.Send(message.TextMessageReq{
		Title:   "这是一条测试消息",
		Content: "test",
		Link:    "https://www.baidu.com",
	})
	if err != nil {
		t.Error(err)
	}
}

func TestImage(t *testing.T) {
	err := clientProxy.Send(message.ImageMessageReq{
		TextMessageReq: message.TextMessageReq{
			Title:   "这是一条测试消息",
			Content: "test",
			Link:    "https://www.baidu.com",
		},
		Url: "https://wwcdn.weixin.qq.com/node/wework/images/Pic_right@2x.7a03a9d992.png",
	})
	if err != nil {
		t.Error(err)
	}
}

func TestList(t *testing.T) {
	err := clientProxy.Send(message.ListMessageReq{Items: []message.ImageMessageReq{
		{TextMessageReq: message.TextMessageReq{
			Title:   "这是一条测试消息",
			Content: "test",
			Link:    "https://www.baidu.com",
		},
			Url: "https://wwcdn.weixin.qq.com/node/wework/images/Pic_right@2x.7a03a9d992.png"},
		{TextMessageReq: message.TextMessageReq{
			Title:   "这是一条测试消息",
			Content: "test",
			Link:    "https://www.baidu.com",
		},
			Url: "https://wwcdn.weixin.qq.com/node/wework/images/Pic_right@2x.7a03a9d992.png"},
	}})
	if err != nil {
		t.Error(err)
	}
}
