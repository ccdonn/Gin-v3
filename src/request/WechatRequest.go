package request

// WechatRequest :
type WechatRequest struct {
	Wechat   string `form:"wc"`
	BrandKey string `form:"brandKey"`
}
