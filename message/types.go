package message

type (
	TextMessageReq struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Link    string `json:"link"`
	}

	ImageMessageReq struct {
		TextMessageReq
		Url string `json:"url"`
	}

	ListMessageReq struct {
		Items []ImageMessageReq `json:"items"`
	}
)
