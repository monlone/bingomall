package model

type CallbackInfo struct {
	AccessToken  string `json:"access_token" weChat:"access_token"`
	RefreshToken string `json:"refresh_token" weChat:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in" weChat:"expires_in"`
	OpenId       string `json:"openid" weChat:"openid"`
	UnionId      string `json:"unionid" weChat:"unionid"`
	Nickname     string `json:"nickname" weChat:"nickname"`
	HeadImageURL string `json:"head_image_url,omitempty" weChat:"HeadImageURL"`
}

type WechatUserInfo struct {
	AccessToken  string `json:"access_token" weChat:"access_token"`
	OpenId       string `json:"openid" weChat:"openid"`
	HeadImageURL string `json:"head_image_url,omitempty" weChat:"HeadImageURL"`
	Nickname     string `json:"nickname" weChat:"nickname"`
	UserType     uint8  `json:"user_type"`
	UserID       string `json:"user_id"`
}

type AppWechatPay struct {
	Appid     string `json:"appid"`
	Noncestr  string `json:"noncestr"`
	Package   string `json:"package"`
	Partnerid uint64 `json:"partnerid"`
	PaySign   string `json:"paySign"`
	Prepayid  string `json:"prepayid"`
	Timestamp string `json:"timestamp"`
}
