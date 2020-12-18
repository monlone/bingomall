package control

import (
	"encoding/json"
	"errors"
	"fmt"
	"bingomall/constant"
	helper "bingomall/helpers"
	"bingomall/helpers/convention"
	"bingomall/library/gopay"
	"bingomall/library/gopay/client"
	"bingomall/library/gopay/common"
	payconstant "bingomall/library/gopay/constant"
	"bingomall/library/wxpay"
	"bingomall/models"
	"bingomall/repositories"
	service "bingomall/services"
	"bingomall/system"
	"github.com/chanxuehong/rand"
	"github.com/chanxuehong/wechat/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"
)

// app 的 wechat 支付
// @Summary app 的 wechat 支付,传 product_id 按 product_id 查
// @Tags PayController
// @Consumes formData
// @Produce json
// @Param Gin-Access-Token header string true "令牌"
// @Param product_id query string true "产品记录id"
// @Router /api/app/wechat_pay [post]
func AppWeChatPay(context *gin.Context) {
	userId := helper.GetUserID(context)
	initWechatClient()
	pay := &model.Pay{}
	err := context.Bind(pay)
	if err != nil {
		context.JSON(http.StatusOK, helper.JsonObject{
			Code:    1120,
			Message: "pay绑定参数错误",
		})
		return
	}
	db := helper.GetDBByName(constant.DBMerchant)
	tx := db.Begin()
	orderService := service.OrderServiceInstance(repositories.OrderRepositoryInstance(tx))
	order := orderService.GetByOrderIdAndUserId(pay.OrderId, userId)
	if order == nil {
		context.JSON(http.StatusOK, helper.JsonObject{
			Code:    1106,
			Message: "userToken错误",
		})
		context.Abort()
		return
	}
	outTradeNo := helper.GenerateId32()

	//TODO 先这样子保留，为了后面参考，要改的。这样子只支付了最后一个商品的价格，以后要在create里去创建价格，把价格写到order表，
	//不然要从新计算，其实应该重新拉产品列表来计算，因为有的产品可能没有库存了，或者久不支付已经下架了，不然支付是有时效性的，比如12小时内必须支付

	moneyTotal := order.Pay // 计算商品总价,从order里面去取，这个要改
	walletService := service.WalletServiceInstance(repositories.WalletRepositoryInstance(tx))
	scoreUser := walletService.GetScoreByUserID(userId)

	if pay.Score > scoreUser {
		context.JSON(http.StatusOK,
			&helper.JsonObject{
				Code:    4204,
				Message: "积分不足",
			})
		return
	}

	if pay.Score > 0 {
		if pay.Score >= uint64(moneyTotal) { //前端有多少积分传了多少积分，所以要写个判断
			pay.Score = uint64(moneyTotal)
		}
		scoreUser = scoreUser - pay.Score
		moneyTotal = moneyTotal - pay.Score

		err = walletService.ReduceScoreByUserID(userId, pay.Score)
		if err != nil {
			tx.Rollback()
			context.JSON(http.StatusOK,
				&helper.JsonObject{
					Code:    4410,
					Message: err.Error(),
					Content: "扣除积分事务错误",
				})
			return
		}
	}

	moneyLog := service.MoneyLogServiceInstance(repositories.MoneyLogRepositoryInstance(tx))
	moneyLogData := &model.MoneyLog{}
	moneyLogData.TransactionID = outTradeNo
	moneyLogData.UserID = userId
	moneyLogData.Type = constant.MoneyPay
	moneyLogData.OrderId = pay.OrderId
	moneyLogData.Cost = pay.Score
	moneyLogData.Income = moneyTotal
	moneyLogData.RelationUserID = userId
	moneyLogData.Describe = "消费积分购买商品"

	helper.ServiceLogger.Println("消费积分购买商品，添加到moneyLog表，moneyLogData：", helper.Json(moneyLogData))
	err = moneyLog.Save(moneyLogData) // TODO 要是用户不支付，积分是要退回的，为了提高转化率，暂时没有做。
	if err != nil {
		tx.Rollback()
		context.JSON(http.StatusOK,
			&helper.JsonObject{
				Code:    4412,
				Message: err.Error(),
				Content: "moneyLog事务错误",
			})
		return
	}

	charge := new(common.Charge)
	charge.PayMethod = payconstant.WECHAT_APP
	charge.MoneyFee = moneyTotal
	charge.Describe = "商城产品" //TODO 这个要改，https://pay.weixin.qq.com/wiki/doc/api/app/app.php?chapter=9_1 的body字段
	charge.TradeNum = outTradeNo
	charge.CallbackURL = constant.WechatPayCallback
	charge.ProfitSharing = "Y"
	data := make(map[string]string)

	if moneyTotal > 0 {
		data, err = gopay.Pay(charge)
		if err != nil {
			context.JSON(http.StatusOK, helper.JsonObject{
				Code:    1106,
				Message: err.Error(),
			})
			context.Abort()
		}
	} else {
		data["prepayid"] = ""
		data["total_fee"] = "0"
	}

	if charge.MoneyFee == 0 { //用户的积分已经够支付了，不用微信支付了
		order.Status = constant.OrderPaySuccess
	}

	err = orderService.Update(order)

	orderProductService := service.OrderProductServiceInstance(repositories.OrderProductRepositoryInstance(tx))
	err = orderProductService.UpdateOrderProductByOrderId(order.ID)
	if err != nil {
		tx.Rollback()
		context.JSON(http.StatusOK,
			&helper.JsonObject{
				Code:    4411,
				Message: err.Error(),
				Content: "order事务错误",
			})
		return
	}
	tx.Commit()
	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: "ok",
			Content: data,
		})
}

// app 的 wechat 支付回调
// @Summary app 的 wechat 支付回调
// @Tags PayController
// @Accept json
// @Produce json
// @Success 200 {object} model.AppWechatPay
// @Router /api/app/wechat_pay_callback [get]
func WechatAppPayCallback(context *gin.Context) {
	result, err := gopay.WeChatAppCallback(context.Writer, context.Request)
	if err != nil {
		helper.ErrorLogger.Println("gopay.WeChatAppCallback err:", helper.Json(err))
		wechatPayCallbackError(context.Writer)
		return
	}
	helper.ServiceLogger.Println("WechatAppPayCallback:", helper.Json(result))

	if result != nil && result.ReturnCode != "SUCCESS" {
		wechatPayCallbackError(context.Writer)
	}

	db := helper.GetDBByName(constant.DBMerchant)
	tx := db.Begin()
	orderService := service.OrderServiceInstance(repositories.OrderRepositoryInstance(tx))
	order := orderService.GetByOutTradNo(result.OutTradeNO)
	helper.ServiceLogger.Println("WechatAppPayCallback, order:", helper.Json(order))
	if order == nil {
		tx.Rollback()
		helper.ErrorLogger.Println("gopay.WeChatAppCallback order err:", helper.Json(err))
		wechatPayCallbackError(context.Writer)
		return
	}

	if strconv.FormatUint(order.Pay, 10) != result.TotalFee {
		helper.ErrorLogger.Println("gopay.WeChatAppCallback TotalFee not equal, result.TotalFee:", result.TotalFee, ",order.Pay:", order.Pay)
		wechatPayCallbackError(context.Writer)
		return
	}

	if order.Status == constant.OrderWaitedForPay || order.Status == constant.OrderWaitedForCheck {
		order.Status = constant.OrderPaySuccess
		order.TransactionID = result.TransactionID
		//order.ShopDetail = nil //ShopDetail是指针的话，就可以用nil，不然就不行。原来定义为指针的
		//order.ProductDetail = nil
		err := orderService.SaveOrUpdate(order)
		if err != nil {
			tx.Rollback()
			helper.ErrorLogger.Println("gopay.WeChatAppCallback orderService.SaveOrUpdate err:", helper.Json(err))
			wechatPayCallbackError(context.Writer)
			return
		}

		money := float64(order.Pay)

		//TODO 以后再处理，这里只是记录个日志，以后加debug是否开启来处理
		//userService := service.UserServiceInstance(repositories.UserRepositoryInstance(helper.GetUserDB()))
		//user := userService.GetByUserID(order.UserID)
		//helper.ServiceLogger.Println("WechatAppPayCallback, user:", helper.Json(user), ",money:", money,
		//	",order:", helper.Json(order))

		err = profitSharingToSelf(tx, money, order) // 用户自己消费返现
		if err != nil {
			tx.Rollback()
			helper.ErrorLogger.Println("gopay.WeChatAppCallback profitSharingToSelf err:", helper.Json(err))
			wechatPayCallbackError(context.Writer)
			return
		}
	}

	tx.Commit()
	wechatPayCallbackSuccess(context.Writer)
}

func initWechatClient() {
	client.InitWxAppClient(&client.WechatAppClient{
		AppID: constant.WxAppId,
		MchID: constant.MchID,
		Key:   constant.WxApiKey,
	})
}

func initWechatMiniProgramClient() {
	client.InitWxMiniProgramClient(&client.WechatMiniProgramClient{
		AppID: constant.WxAppId,
		MchID: constant.MchID,
		Key:   constant.WxApiKey,
	})
}

func wechatPayCallbackSuccess(w http.ResponseWriter) {
	var returnCode = "SUCCESS"
	var returnMsg = "Ok"
	formatStr := `<xml><return_code><![CDATA[%s]]></return_code>
                  <return_msg>![CDATA[%s]]</return_msg></xml>`
	returnBody := fmt.Sprintf(formatStr, returnCode, returnMsg)
	code, _ := w.Write([]byte(returnBody))
	fmt.Println(code)
}

func wechatPayCallbackError(w http.ResponseWriter) {
	var returnCode = "FAIL"
	var returnMsg = ""
	formatStr := `<xml><return_code><![CDATA[%s]]></return_code>
                  <return_msg>![CDATA[%s]]</return_msg></xml>`
	returnBody := fmt.Sprintf(formatStr, returnCode, returnMsg)
	code, _ := w.Write([]byte(returnBody))
	fmt.Println(code)
}

func AddReceiver(context *gin.Context) {
	claims, ok := context.Get("claims")
	if !ok {
		context.JSON(http.StatusOK, helper.JsonObject{
			Code:    6001,
			Message: "token错误",
		})
		return
	}
	userInfo := claims.(*system.CustomClaims)

	userService := service.UserServiceInstance(repositories.UserRepositoryInstance(helper.GetUserDB()))
	user := userService.GetByUserID(userInfo.ID)
	if user == nil || user.UnionId == "" {
		context.JSON(http.StatusOK,
			&helper.JsonObject{
				Code:    4223,
				Message: "userId 非法或者unionId不存在",
			})
		return
	}

	httpClient := util.DefaultHttpClient
	receiverMap := make(map[string]string)
	receiverMap["type"] = "PERSONAL_OPENID"
	receiverMap["account"] = user.OpenId
	receiverMap["relation_type"] = "USER"
	receiver, _ := json.Marshal(receiverMap)
	nonceStr := string(rand.NewHex())

	var param = make(map[string]string)
	param["mch_id"] = constant.MchID
	param["appid"] = constant.WxAppId
	param["nonce_str"] = nonceStr
	param["sign_type"] = "HMAC-SHA256"
	param["receiver"] = string(receiver)

	// 创建支付账户
	account := wxpay.NewAccount(constant.WxAppId, constant.MchID, constant.WxApiKey, false)

	// 新建微信支付客户端
	wechatClient := wxpay.NewClient(account)

	// 设置http请求超时时间
	wechatClient.SetHttpConnectTimeoutMs(2000)

	// 设置http读取信息流超时时间
	wechatClient.SetHttpReadTimeoutMs(1000)

	// 更改签名类型
	wechatClient.SetSignType(wxpay.HMACSHA256)

	sign := wechatClient.Sign(param)
	param["sign"] = sign

	data := strings.NewReader(postBody(param))

	httpResp, err := httpClient.Post("https://api.mch.weixin.qq.com/pay/profitsharingaddreceiver", constant.ApplicationXML, data)
	if err != nil {
		return
	}

	if httpResp.StatusCode != http.StatusOK {
		err = fmt.Errorf("http.Status: %s", httpResp.Status)
		return
	}

	str, _ := ioutil.ReadAll(httpResp.Body)
	fmt.Println("str:::", string(str))
}

func addWechatReceiver(userId uint64) (err error) {
	userService := service.UserServiceInstance(repositories.UserRepositoryInstance(helper.GetUserDB()))
	user := userService.GetByUserID(userId)
	if user == nil || user.UnionId == "" {
		return errors.New("userId非法或者unionId不存在")
	}

	if user.OpenId == "" {
		return errors.New("user的OpenId不存在")
	}

	httpClient := util.DefaultHttpClient
	receiverMap := make(map[string]string)
	receiverMap["type"] = "PERSONAL_OPENID"
	receiverMap["account"] = user.OpenId
	receiverMap["relation_type"] = "USER"
	receiver, _ := json.Marshal(receiverMap)
	nonceStr := string(rand.NewHex())

	var param = make(map[string]string)
	param["mch_id"] = constant.MchID
	param["appid"] = constant.WxAppId
	param["nonce_str"] = nonceStr
	param["sign_type"] = "HMAC-SHA256"
	param["receiver"] = string(receiver)

	// 创建支付账户
	account := wxpay.NewAccount(constant.WxAppId, constant.MchID, constant.WxApiKey, false)

	// 新建微信支付客户端
	wechatClient := wxpay.NewClient(account)

	// 设置http请求超时时间
	wechatClient.SetHttpConnectTimeoutMs(2000)

	// 设置http读取信息流超时时间
	wechatClient.SetHttpReadTimeoutMs(1000)

	// 更改签名类型
	wechatClient.SetSignType(wxpay.HMACSHA256)

	sign := wechatClient.Sign(param)
	param["sign"] = sign

	data := strings.NewReader(postBody(param))

	httpResp, err := httpClient.Post("https://api.mch.weixin.qq.com/pay/profitsharingaddreceiver", constant.ApplicationXML, data)
	if err != nil {
		return
	}

	if httpResp.StatusCode != http.StatusOK {
		err = fmt.Errorf("http.Status: %s", httpResp.Status)
		return
	}

	str, _ := ioutil.ReadAll(httpResp.Body)
	fmt.Println("str:::", string(str))

	return nil
}

func profitSharing(dbT *gorm.DB, order *model.Order, orderProduct *model.OrderProduct, money float64, checkShopId uint64, ) (err error) {
	helper.ServiceLogger.Println("order in profitSharing:", helper.Json(order))

	userService := service.UserServiceInstance(repositories.UserRepositoryInstance(helper.GetUserDB()))
	user := userService.GetByUserID(order.UserID)
	if user == nil || user.UnionId == "" {
		helper.ServiceLogger.Println("用户不存在")
		return
	}

	var receiveList []interface{}

	//分账给一级上级
	err, receiverMap := profitSharingToPerson(user, dbT, money, order)
	if err != nil {
		return
	}

	if _, ok := receiverMap["type"]; ok {
		//存在
		receiveList = append(receiveList, receiverMap)
	}

	// 分账给二级上线
	err, receiverMapSecond := profitSharingToPersonSecond(user, dbT, money, order)
	if err != nil {
		return
	}

	if _, ok := receiverMapSecond["type"]; ok {
		receiveList = append(receiveList, receiverMapSecond)
	}

	//分账给平台代核销商户
	err, receiveMerchant := profitSharingToMerchant(dbT, money, order, checkShopId)

	if err != nil {
		return
	}

	if _, ok := receiveMerchant["type"]; ok {
		receiveList = append(receiveList, receiveMerchant)
	}

	if len(receiveList) == 0 {
		return
	}

	account := wxpay.NewAccount(constant.WxAppId, constant.MchID, constant.WxApiKey, false)
	// 新建微信支付客户端
	wechatClient := wxpay.NewClient(account)
	// 设置证书
	wechatConfig := system.GetWechatConfig()
	account.SetCertData(wechatConfig.PayApiClientCert)
	// 设置http请求超时时间
	wechatClient.SetHttpConnectTimeoutMs(2000)
	// 设置http读取信息流超时时间
	wechatClient.SetHttpReadTimeoutMs(1000)
	// 更改签名类型
	wechatClient.SetSignType(wxpay.HMACSHA256)

	nonceStr := string(rand.NewHex())
	var param = make(map[string]string)
	outTradeNo := helper.GenerateId32() //最长只能有32位，所以不用helper.GenerateId36()。sss
	param["mch_id"] = constant.MchID
	param["appid"] = constant.WxAppId
	param["sign_type"] = "HMAC-SHA256"
	param["transaction_id"] = order.TransactionID
	param["out_order_no"] = outTradeNo
	receivers := helper.Json(receiveList)
	helper.ServiceLogger.Println("receivers in profitSharing receivers：", receivers)
	param["receivers"] = receivers
	param["nonce_str"] = nonceStr
	data, err := wechatClient.PostWithCert(constant.MultiProfitSharing, param)
	if err != nil {
		helper.ErrorLogger.Errorln("分账错误：", err, ",param:", helper.Json(param), ",data:", data)
		return
	}
	helper.ServiceLogger.Println("in function wechatClient.PostWithCert data：", data)

	err = profitSharingFinish(order)

	return
	//微信支付零钱到到购买者，要公众号开能90天以后才能用
	//account = wxpay.NewAccount(constant.MPWechatAppId, constant.MchID, constant.MPWxApiKey, false)
	////wechatClient = wxpay.NewClient(account)
	////account.SetCertData(wechatConfig.PayApiClientCert)
	//wechatClient.SetSignType(wxpay.MD5)
	//var paramPerson = make(map[string]string)
	//paramPerson["mchid"] = constant.MchID
	//paramPerson["mch_appid"] = constant.WxAppId
	//paramPerson["nonce_str"] = nonceStr
	//paramPerson["sign_type"] = "MD5"
	//paramPerson["partner_trade_no"] = order.TransactionID
	//paramPerson["out_order_no"] = outTradeNo
	//paramPerson["check_name"] = "NO_CHECK"
	//paramPerson["openid"] = userLevel.OpenId
	//paramPerson["amount"] = strconv.FormatInt(int64(money*float64(product.DiscountUser)/100), 10)
	//paramPerson["spbill_create_ip"] = uip.LocalIP()
	//paramPerson["desc"] = "person"
	//httpRespPerson, err := wechatClient.EnterprisePostWithCert(wxpay.EnterprisePayTransfersUrl, paramPerson)
	//fmt.Println("PayToClient httpRespPerson:", httpRespPerson, ",err:", err)

	//initWechatClient()
	//charge := new(common.Charge)
	//charge.PayMethod = payconstant.WECHAT_APP
	//charge.MoneyFee = money * percent // 这个sdk以元为单位的，系统是以分为单位的，所以要转换。坑
	//charge.Describe = "个人返现"
	//charge.TradeNum = outTradeNo
	//charge.OpenID = firstLevelOpenId
	//charge.ProfitSharing = "Y"
	//data, err := gopay.PayToClient(charge)
	//fmt.Println("PayToClient data:", data, ",err:", err)
}

//分账给上线，一级分账
//现在是分账到微信的
func profitSharingToPerson(user *model.User, dbT *gorm.DB, money float64, order *model.Order) (err error, receiverMap map[string]interface{}) {
	userService := service.UserServiceInstance(repositories.UserRepositoryInstance(helper.GetUserDB()))
	receiverMap = make(map[string]interface{})

	if user == nil || user.MultiLevel == 0 {
		helper.ServiceLogger.Println("分账用户的上线不存在,消费用户user：", helper.Json(user))
		return
	}
	userLevel := userService.GetByUserID(user.MultiLevel)

	if userLevel == nil || userLevel.OpenId == "" {
		helper.ServiceLogger.Println("分账用户的上线不存在,userLevel：", helper.Json(userLevel))
		return
	}

	err = addWechatReceiver(userLevel.ID)
	if err != nil {
		helper.ServiceLogger.Println("添加分账用户失败：", err)
		return
	}

	productList := GetProductListByOrderId(dbT, order.ID)
	for _, p := range productList {
		amount := uint64(math.Floor(money*float64(p.DiscountLevel*constant.FirstLevelProfit)/100 + 0.5))

		receiverMap["type"] = "PERSONAL_OPENID"
		receiverMap["account"] = userLevel.OpenId
		receiverMap["amount"] = amount
		receiverMap["description"] = "第一次分到个人" //分账到上线

		moneyLog := service.MoneyLogServiceInstance(repositories.MoneyLogRepositoryInstance(dbT))
		moneyLogData := &model.MoneyLog{}
		moneyLogData.TransactionID = order.TransactionID
		moneyLogData.UserID = order.UserID
		moneyLogData.Type = constant.MoneyLogUser
		moneyLogData.OrderId = order.ID
		moneyLogData.Income = uint64(money)
		moneyLogData.RelationUserID = user.MultiLevel
		moneyLogData.Describe = "一级下线消费提成"

		helper.ServiceLogger.Println("profitSharingToPerson添加到moneyLog表，moneyLogData：", helper.Json(moneyLogData))
		err = moneyLog.Save(moneyLogData)
		if err != nil {
			helper.ServiceLogger.Println("一级下线消费提成分账失败：", err)
			return
		}
	}

	return
}

//分账给二级上线
func profitSharingToPersonSecond(user *model.User, dbT *gorm.DB, money float64, order *model.Order) (err error, receiverMap map[string]interface{}) {
	userService := service.UserServiceInstance(repositories.UserRepositoryInstance(helper.GetUserDB()))

	if user == nil || user.MultiLevel == 0 {
		helper.ServiceLogger.Println("分账用户的上线不存在,消费用户user：", helper.Json(user))
		return
	}
	userLevelFirst := userService.GetByUserID(user.MultiLevel)

	if userLevelFirst == nil {
		helper.ServiceLogger.Println("分账用户的一级上线不存在,userLevelFirst：", helper.Json(userLevelFirst))
		return
	}
	userLevel := userService.GetByUserID(userLevelFirst.MultiLevel)
	if userLevel == nil || userLevel.OpenId == "" {
		helper.ServiceLogger.Println("分账用户的二级上线不存在,userLevel：", helper.Json(userLevel))
		return
	}

	err = addWechatReceiver(userLevel.ID)
	if err != nil {
		helper.ServiceLogger.Println("添加分账用户失败：", err)
		return
	}

	productList := GetProductListByOrderId(dbT, order.ID)
	for _, p := range productList {
		amount := uint64(math.Floor(money*float64(p.DiscountLevel*constant.SecondLevelProfit)/100 + 0.5))

		receiverMap = make(map[string]interface{})
		receiverMap["type"] = "PERSONAL_OPENID"
		receiverMap["account"] = userLevel.OpenId
		receiverMap["amount"] = amount
		receiverMap["description"] = "第二次分到个人" //分账到上线

		moneyLog := service.MoneyLogServiceInstance(repositories.MoneyLogRepositoryInstance(dbT))
		moneyLogData := &model.MoneyLog{}
		moneyLogData.TransactionID = order.TransactionID
		moneyLogData.UserID = order.UserID
		moneyLogData.Type = constant.MoneyLogUser
		moneyLogData.OrderId = order.ID
		moneyLogData.Income = uint64(money)
		moneyLogData.RelationUserID = user.MultiLevel
		moneyLogData.Describe = "二级下线消费提成"

		helper.ServiceLogger.Println("profitSharingToPersonSecond添加到moneyLog表，moneyLogData：", helper.Json(moneyLogData))
		err = moneyLog.Save(moneyLogData)
		if err != nil {
			helper.ServiceLogger.Println("二级下线消费提成分账失败：", err)
			return
		}
	}

	return
}

// 分账给代核销商户
func profitSharingToMerchant(dbT *gorm.DB, money float64, order *model.Order, checkShopId uint64, ) (err error, receiverMap map[string]interface{}) {
	receiverMap = make(map[string]interface{})
	userShopService := service.UserShopServiceInstance(repositories.UserShopRepositoryInstance(merchantDB))
	shop := userShopService.GetByShopId(checkShopId)
	if len(shop) == 0 {
		helper.ServiceLogger.Println("商店不存在,", "checkShopId:", checkShopId)
		return
	}

	MerchantUserIDs := helper.ModelObjectToSlice(shop, "UserID")
	MerchantUserID := MerchantUserIDs[0] //现在的逻辑一个商户只能开一个商户代平台核销的店，这个店由平台来上架东西，商户来核销，从中拿分成

	err = addWechatReceiver(MerchantUserID)
	if err != nil {
		helper.ServiceLogger.Println("添加分账用户失败：", err)
		return
	}

	userService := service.UserServiceInstance(repositories.UserRepositoryInstance(helper.GetUserDB()))
	merchantUserInfo := userService.GetByUserID(MerchantUserID)

	productList := GetProductListByOrderId(dbT, order.ID)
	for _, p := range productList {
		amount := uint64(math.Floor(money*float64(p.DiscountMerchant)/100 + 0.5))

		baseService := service.BaseServiceInstance(repositories.BaseRepositoryInstance(dbT))
		shopService := service.ShopServiceInstance(repositories.ShopRepositoryInstance(dbT))
		orderShop := shopService.GetByShopId(p.ShopId)
		if baseService.IsPlatformOfficial(orderShop) == true {
			amount = uint64(math.Floor(float64(amount)*p.DiscountMerchant/100 + 0.5))
		}
		receiverMap["type"] = "PERSONAL_OPENID"
		receiverMap["account"] = merchantUserInfo.OpenId
		receiverMap["amount"] = amount
		receiverMap["description"] = "分账到平台代核销商户" //分账到平台代核销商户

		moneyLog := service.MoneyLogServiceInstance(repositories.MoneyLogRepositoryInstance(dbT))
		moneyLogData := &model.MoneyLog{}
		moneyLogData.TransactionID = order.TransactionID
		moneyLogData.UserID = MerchantUserID
		moneyLogData.Type = constant.MoneyLogUser
		moneyLogData.OrderId = order.ID
		moneyLogData.Income = uint64(money)
		moneyLogData.RelationUserID = 0 //不好记录，暂时为空
		moneyLogData.Describe = "代平台核销返现"

		helper.ServiceLogger.Println("profitSharingToMerchant添加到moneyLog表，moneyLogData：", helper.Json(moneyLogData))
		err = moneyLog.Save(moneyLogData)
		if err != nil {
			helper.ServiceLogger.Println("代平台核销返现失败：", err)
		}
	}

	return
}

//分账完结
func profitSharingFinish(order *model.Order) (err error) {
	var receiveList []interface{}
	account := wxpay.NewAccount(constant.WxAppId, constant.MchID, constant.WxApiKey, false)
	wechatClient := wxpay.NewClient(account)
	wechatConfig := system.GetWechatConfig()
	account.SetCertData(wechatConfig.PayApiClientCert)
	wechatClient.SetHttpConnectTimeoutMs(2000)
	wechatClient.SetHttpReadTimeoutMs(1000)
	wechatClient.SetSignType(wxpay.HMACSHA256)

	nonceStr := string(rand.NewHex())
	var param = make(map[string]string)
	outTradeNo := helper.GenerateId32()  //最长只能有32位，所以不用helper.GenerateId36()。sss
	param["appid"] = constant.WxAppId    //微信分配的公众账号ID
	param["mch_id"] = constant.MchID     //微信支付分配的商户号
	param["sub_mch_id"] = constant.MchID //TODO 微信支付分配的子商户号,不知道是哪个，好玉米没有
	param["sign_type"] = "HMAC-SHA256"
	param["transaction_id"] = order.TransactionID
	param["out_order_no"] = outTradeNo
	param["nonce_str"] = nonceStr
	param["description"] = "完成分账"

	helper.ErrorLogger.Errorln("in function profitSharingToPerson receiveList：", helper.Json(receiveList))
	data, err := wechatClient.PostWithCert(constant.ProfitSharingFinish, param)
	if err != nil {
		helper.ErrorLogger.Errorln("分账结束线错误：", err, ",param:", helper.Json(param), ",data:", helper.Json(data))
	}
	helper.ErrorLogger.Errorln("in function profitSharingFinish data：", helper.Json(data))

	return
}

// 用户消费返现，暂时是返积分
func profitSharingToSelf(dbT *gorm.DB, money float64, order *model.Order) (err error) {
	wallet := service.WalletServiceInstance(repositories.WalletRepositoryInstance(dbT))
	walletData := wallet.GetByUserID(order.UserID)

	productList := GetProductListByOrderId(dbT, order.ID)
	for _, p := range productList {
		amount := uint64(math.Floor(money*float64(p.DiscountUser)/100 + 0.5))
		helper.ServiceLogger.Println("profitSharingToSelf, amount:", amount)

		if walletData == nil {
			walletData = &model.Wallet{}
			walletData.Score = amount
		} else {
			walletData.Score = amount + walletData.Score
		}
		moneyLog := service.MoneyLogServiceInstance(repositories.MoneyLogRepositoryInstance(dbT))
		moneyLogData := &model.MoneyLog{}
		moneyLogData.TransactionID = order.TransactionID
		moneyLogData.UserID = order.UserID
		moneyLogData.Type = constant.MoneyLogSelf
		moneyLogData.OrderId = order.ID
		moneyLogData.Income = uint64(money)
		moneyLogData.RelationUserID = order.UserID
		moneyLogData.Describe = "自己消费返现"

		helper.ServiceLogger.Println("添加到moneyLog表，moneyLogData：", helper.Json(moneyLogData))
		err = moneyLog.Save(moneyLogData)
		if err != nil {
			return
		}
	}

	walletData.UserID = order.UserID
	helper.ServiceLogger.Println("profitSharingToSelf, walletData:", walletData)
	err = wallet.SaveOrUpdate(walletData)
	if err != nil {
		return
	}

	return
}

func postBody(m map[string]string) string {
	formatData := "<xml>" +
		"<mch_id>%s</mch_id>" +
		"<appid>%s</appid>" +
		"<nonce_str>%s</nonce_str>" +
		"<sign>%s</sign>" +
		"<sign_type>%s</sign_type>" +
		"<receiver>%s</receiver>" +
		"</xml>"
	data := fmt.Sprintf(formatData, m["mch_id"], m["appid"], m["nonce_str"], m["sign"], m["sign_type"], m["receiver"])

	return data
}

// 微信小程序的 wechat 支付
// @Summary app 的 wechat 支付,传 productId 按 productId 查
// @Tags PayController
// @Consumes formData
// @Produce json
// @Param Gin-Access-Token header string true "令牌"
// @Param productId query string true "产品记录id"
// @Router /api/app/micro_wechat_pay [post]
func MicroWeChatPay(context *gin.Context) {
	userId := helper.GetUserID(context)
	openId := helper.GetOpenId(context)
	initWechatMiniProgramClient()
	pay := &model.Pay{}
	err := context.Bind(pay)
	if err != nil {
		context.JSON(http.StatusOK, helper.JsonObject{
			Code:    1120,
			Message: "pay绑定参数错误",
		})
		return
	}
	outTradeNo := helper.GenerateId32() //最长只能有32位，所以不用helper.GenerateId36()。
	if len(pay.CouponList) > 0 {
		//TODO 要使用coupon
		return
	}

	db := helper.GetDBByName(constant.DBMerchant)
	orderService := service.OrderServiceInstance(repositories.OrderRepositoryInstance(db))
	var order *model.Order
	order = orderService.GetByOrderIdAndUserId(pay.OrderId, userId)
	if order == nil {
		helper.ErrorLogger.Errorln("order:", order, ",pay:", pay, "userId:", userId)
		helper.ErrorLogger.Errorln("orderId err!")
		context.JSON(http.StatusOK, helper.JsonObject{
			Code:    1106,
			Message: "orderId err!",
		})
		return
	}

	//TODO 先这样子保留，为了后面参考，要改的。这样子只支付了最后一个商品的价格，以后要在create里去创建价格，把价格写到order表，
	//不然要从新计算，其实应该重新拉产品列表来计算，因为有的产品可能没有库存了，或者久不支付已经下架了，不然支付是有时效性的，比如12小时内必须支付
	//productService := service.ProductServiceInstance(repositories.ProductRepositoryInstance(helper.GetDBByName(constant.DBMerchant)))
	//product := productService.GetByProductId(order.ProductDetail.ProductId)
	//for _, p := range productList {
	//	product = p
	//}
	//var describe string
	var moneyTotal uint64
	for _, op := range order.OrderProduct {
		//describe += order.ProductDetail.Title
		moneyTotal += op.Money
	}

	tx := db.Begin()
	walletService := service.WalletServiceInstance(repositories.WalletRepositoryInstance(tx))
	scoreUser := walletService.GetScoreByUserID(userId)

	if pay.Score > scoreUser {
		context.JSON(http.StatusOK,
			&helper.JsonObject{
				Code:    4204,
				Message: "积分不足",
			})
		return
	}

	if pay.Score > 0 {
		if pay.Score >= moneyTotal { //前端传了多少积分，所以要写个判断
			pay.Score = moneyTotal
		}
		scoreUser = scoreUser - pay.Score
		moneyTotal = moneyTotal - pay.Score

		err = walletService.ReduceScoreByUserID(userId, pay.Score)
		if err != nil {
			tx.Rollback()
			context.JSON(http.StatusOK,
				&helper.JsonObject{
					Code:    4410,
					Message: err.Error(),
					Content: "扣除积分事务错误",
				})
			return
		}
	}

	moneyLog := service.MoneyLogServiceInstance(repositories.MoneyLogRepositoryInstance(tx))
	moneyLogData := &model.MoneyLog{}
	moneyLogData.TransactionID = outTradeNo
	moneyLogData.UserID = userId
	moneyLogData.Type = constant.MoneyPay
	moneyLogData.OrderId = pay.OrderId
	moneyLogData.Cost = moneyTotal
	moneyLogData.CostScore = pay.Score
	moneyLogData.RelationUserID = userId
	moneyLogData.Describe = "购买商品"

	helper.ServiceLogger.Println("消费积分购买商品，添加到moneyLog表，moneyLogData：", helper.Json(moneyLogData))
	err = moneyLog.Save(moneyLogData) // TODO 要是用户不支付，积分是要退回的，但是为了提高支付率，暂时不退回，以后再做。
	if err != nil {
		tx.Rollback()
		context.JSON(http.StatusOK,
			&helper.JsonObject{
				Code:    4412,
				Message: err.Error(),
				Content: "moneyLog事务错误",
			})
		return
	}

	charge := new(common.Charge)
	charge.PayMethod = payconstant.WECHAT_MINI_PROGRAM
	charge.MoneyFee = moneyTotal
	charge.Describe = "支付定单：" + convention.Uint64ToString(order.ID)
	charge.TradeNum = outTradeNo
	charge.CallbackURL = constant.WechatPayCallback
	charge.ProfitSharing = "Y"
	charge.OpenID = openId
	data := make(map[string]string)

	if moneyTotal > 0 {
		data, err = gopay.Pay(charge)
		if err != nil {
			context.JSON(http.StatusOK, helper.JsonObject{
				Code:    1106,
				Message: err.Error(),
			})
			context.Abort()
		}
	} else {
		data["prepayid"] = ""
		data["total_fee"] = "0"
	}

	order.Status = constant.OrderWaitedForPay

	err = orderService.Update(order)

	if err != nil {
		tx.Rollback()
		context.JSON(http.StatusOK,
			&helper.JsonObject{
				Code:    4411,
				Message: err.Error(),
				Content: "order事务错误",
			})
	}

	tx.Commit()
	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: "ok",
			Content: data,
		})
}

//小程序客户端支付成功后，更新订单状态
func UpdateOrderAfterPay(context *gin.Context) {
	userId := helper.GetUserID(context)
	db := helper.GetDBByName(constant.DBMerchant)
	tx := db.Begin()
	orderService := service.OrderServiceInstance(repositories.OrderRepositoryInstance(tx))
	var order *model.Order
	pay := &model.Pay{}
	err := context.Bind(pay)
	if err != nil {
		context.JSON(http.StatusOK, helper.JsonObject{
			Code:    1120,
			Message: "pay绑定参数错误",
		})
		return
	}
	order = orderService.GetByOrderIdAndUserId(pay.OrderId, userId)
	if order == nil {
		context.JSON(http.StatusOK, helper.JsonObject{
			Code:    1121,
			Message: "支付订单错误",
		})
		return
	}

	order.Status = constant.OrderWaitedForCheck
	err = orderService.Update(order)

	if err != nil {
		tx.Rollback()
		context.JSON(http.StatusOK,
			&helper.JsonObject{
				Code:    4411,
				Message: err.Error(),
				Content: "order事务错误",
			})
	}

	tx.Commit()
	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: "ok",
			Content: "支付成功！",
		})
}
