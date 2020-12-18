package client

import (
	"errors"
	"fmt"
	"bingomall/library/gopay/common"
	"bingomall/library/gopay/util"
	"strings"
	"time"
)

var defaultWechatAppClient *WechatAppClient

func InitWxAppClient(c *WechatAppClient) {
	defaultWechatAppClient = c
}

// DefaultWechatAppClient 默认微信app客户端
func DefaultWechatAppClient() *WechatAppClient {
	return defaultWechatAppClient
}

// WechatAppClient 微信app支付
type WechatAppClient struct {
	AppID       string       // 公众账号ID
	MchID       string       // 商户号ID
	Key         string       // 密钥
	PrivateKey  []byte       // 私钥文件内容
	PublicKey   []byte       // 公钥文件内容
	httpsClient *HTTPSClient // 双向证书链接
}

// Pay 支付
func (wechatApp *WechatAppClient) Pay(charge *common.Charge) (map[string]string, error) {
	var m = make(map[string]string)
	m["appid"] = wechatApp.AppID
	m["mch_id"] = wechatApp.MchID
	m["nonce_str"] = util.RandomStr()
	m["body"] = TruncatedText(charge.Describe, 32)
	m["out_trade_no"] = charge.TradeNum
	m["total_fee"] = WechatMoneyFeeToString(charge.MoneyFee)
	m["spbill_create_ip"] = util.LocalIP()
	m["notify_url"] = charge.CallbackURL
	m["trade_type"] = "APP"
	m["sign_type"] = "MD5"
	if charge.ProfitSharing != "" {
		m["profit_sharing"] = charge.ProfitSharing
	}

	sign, err := WechatGenSign(wechatApp.Key, m)
	if err != nil {
		return map[string]string{}, errors.New("WechatApp.sign: " + err.Error())
	}

	m["sign"] = sign

	xmlRe, err := PostWechat("https://api.mch.weixin.qq.com/pay/unifiedorder", m, nil)
	if err != nil {
		return map[string]string{}, err
	}

	var c = make(map[string]string)
	c["appid"] = wechatApp.AppID
	c["partnerid"] = wechatApp.MchID
	c["prepayid"] = xmlRe.PrepayID
	c["package"] = "Sign=WXPay"
	c["noncestr"] = util.RandomStr()
	c["timestamp"] = fmt.Sprintf("%d", time.Now().Unix())

	sign2, err := WechatGenSign(wechatApp.Key, c)
	if err != nil {
		return map[string]string{}, errors.New("WechatApp.paySign: " + err.Error())
	}
	c["paySign"] = strings.ToUpper(sign2)

	return c, nil
}

// 支付到用户的微信账号
func (wechatApp *WechatAppClient) PayToClient(charge *common.Charge) (map[string]string, error) {
	return WechatCompanyChange(wechatApp.AppID, wechatApp.MchID, wechatApp.Key, wechatApp.httpsClient, charge)
}

// QueryOrder 查询订单
func (wechatApp *WechatAppClient) QueryOrder(tradeNum string) (common.WeChatQueryResult, error) {
	var m = make(map[string]string)
	m["appid"] = wechatApp.AppID
	m["mch_id"] = wechatApp.MchID
	m["out_trade_no"] = tradeNum
	m["nonce_str"] = util.RandomStr()

	sign, err := WechatGenSign(wechatApp.Key, m)
	if err != nil {
		return common.WeChatQueryResult{}, err
	}

	m["sign"] = sign

	return PostWechat("https://api.mch.weixin.qq.com/pay/orderquery", m, nil)
}
