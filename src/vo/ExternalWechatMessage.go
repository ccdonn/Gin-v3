package vo

type ExternalWechatMessage struct {
	Suc bool `json:"suc"`
	// msg  string     `json:"msg"`
	Data []ExternalWechatVO `json:"data"`
}

type ExternalWechatVO struct {
	Channel string `json:"channel"`
	Value   string `json:"value"`
}
