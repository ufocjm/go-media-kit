package message

type (
	Client interface {
		// SendText 发送文本消息
		SendText(req TextMessageReq) error
		// SendImage 发送图文消息
		SendImage(req ImageMessageReq) error
		// SendList 发送列表消息
		SendList(req ListMessageReq) error
	}
)

type ClientType int

const (
	WorkWechat ClientType = iota
	IYUU       ClientType = iota
)
