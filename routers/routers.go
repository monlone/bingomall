package routers

import (
	"bingomall/controls"
	_ "bingomall/docs"
	"bingomall/system"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"
)

// api 路由注册要鉴权,C端用户
func RegisterApiRoutes(router *gin.Engine) {
	api := router.Group("api")
	// 鉴权
	api.Use(system.JWTAuth())
	//api.GET("get_user_page", control.GetUserPage)
	api.GET("user/info", control.UserInfo)
	api.POST("user/save", control.SaveUser)
	api.POST("user/logout", control.Logout)

	//api.POST("app/user/delete", control.DeleteUser)
	//api.GET("get_all_users", control.GetAllUsers)
	//api.POST("save_role", control.SaveRole)

	api.GET("user/score", control.Score)   //获取用户积分
	api.GET("user/wallet", control.Wallet) //获取用户的钱包

	api.POST("pay/wechat_pay", control.AppWeChatPay)         // APP中拉起 wechat 支付
	api.POST("pay/wechat_pay_success", control.AppWeChatPay) // APP支付成功后的请求的

	api.POST("user/order_create", control.CreateOrder) //微信小程序用户下单
	api.GET("user/order_list", control.OrderList)      //微信小程序获取用户订单列表
	api.GET("user/order_detail", control.OrderDetail)
	api.POST("user/order_verification", control.VerificationOrder) // 核销
	api.POST("pay/micro_wechat_pay", control.MicroWeChatPay)       //微信小程序支付
	api.POST("pay/update_order", control.UpdateOrderAfterPay)      //微信小程序支付成功后，客户端更新定单状态

	api.GET("user/order_statistics", control.Statistics)

	api.GET("user/order/merchant_list", control.MerchantOrderDetailList) // 商户已核销、待核销列表
	api.GET("user/order/liquidate_total", control.LiquidateTotal)        // 商户清账
	api.POST("user/order/liquidate", control.Liquidate)

	api.GET("order/merchant_employee_list", control.MerchantEmployeeOrderDetailList) // 店员查看
	api.GET("order/user_list", control.UserOrderDetailList)                          // 用户每月账单

	// 文件下载
	//router.GET("export_user_infos", control.ExportUserInfos)

	api.POST("app/receiver", control.AddReceiver)

	api.POST("app/feedback/add", control.SaveFeedback) //添加反馈

	api.GET("app/money_log/list", control.MoneyLogList) //获取商户结算列表或用户分账列表

	api.GET("auth/checkToken", control.CheckTokenWeChat) //checkToken

	api.GET("user/address_list", control.AddressList)       //用户收货地址列表
	api.GET("user/address_detail", control.AddressDetail)   //用户收货地址
	api.GET("user/address_default", control.AddressDefault) //用户默认收货地址

	api.POST("user/address_add", control.SaveAddress)               //添加用户收货地址
	api.POST("user/address_update", control.SaveAddress)            //更新用户收货地址
	api.POST("user/set_default_address", control.SetDefaultAddress) //设置用户默认收货地址
	api.POST("user/shoppingCart/add", control.AddToCart)            //添加到购物车
	api.GET("user/shoppingCart/info", control.ShoppingCartInfo)     //获取用户购物车信息

	api.GET("user/cashLog", control.MoneyLogList) //获取用户购物车信息
}

// app 路由注册
func RegisterAppRoutes(router *gin.Engine) {
	app := router.Group("app")
	// 鉴权
	app.Use(system.JWTAuth())
	app.GET("hello", func(context *gin.Context) {
		context.String(http.StatusOK, "Hello APP")
	})
}

// 注册其他需要鉴权的接口,管理员
func RegisterAuthRoutes(router *gin.Engine) {
	api := router.Group("admin")
	api.Use(system.JWTAuth())
	api.POST("api/merchant/save", control.SaveMerchant)
	api.POST("api/shop/save", control.SaveShop)
	api.GET("api/shop/list", control.ShopList)
	api.POST("api/product/save", control.SaveProduct)
	api.POST("api/banner/save", control.SaveBanner)
	api.PATCH("api/banner/save", control.SaveBanner)
	api.DELETE("api/banner/remove", control.RemoveBanner)
	api.GET("api/banner/list_all", control.BannerListAll)
	api.GET("api/product/list_all", control.ProductListAll)
	api.GET("api/merchant/money_log", control.MerchantMoneyLog) //获取全部商户结算列表
	api.POST("api/money_log/update", control.UpdateMoneyLog)    //更新商户结算状态
}

// 注册不需要鉴权的 接口
func RegisterOpenRoutes(router *gin.Engine) {
	//中间件，只有在中间件之后注册的路由才会走中间件
	router.POST("api/login", control.Login)
	router.POST("api/register", control.Register)
	router.POST("api/version/check", control.Check)

	// 使用gin-swagger 中间件
	router.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// wechat web授权
	//router.GET("api/web/wechat_auth", control.Auth)
	//router.GET("api/web/wechat_exchange_token", control.ExchangeToken)
	router.GET("api/web/wechat_mp_register", control.WechatMPRegister) //推荐下线，分销(微信公众号网页)注册
	router.Any("api/web/mpcallback", control.WechatMPCallBack)         //分销注册（微信公众号）回调,微信不支持下划线

	router.POST("api/web/merchant/wechat_mp_register", control.MerchantWechatMPRegister)           //销售员分享给商户注册,申请入驻
	router.POST("api/web/merchant/register", control.MerchantRegister)                             //商户自主注册
	router.GET("api/web/merchant/wechat_mp_register_token", control.MerchantWechatMPRegisterToken) //生成销售员分享给商户注册token
	//router.Any("api/web/mpmerchantcb", control.MerchantWechatMPCallBack)                           //分销注册（微信公众号）回调,微信不支持下划线

	// app wechat 获取token
	router.POST("api/app/wechat_exchange_token", control.AppExchangeToken) //app拉起微信授权登录,这是第一步，新用户会加入到表20200723

	//wechat mini program 注册
	router.POST("api/auth/registerByWeChat", control.RegisterByWeChat) //通过微信小程序注册
	router.POST("api/auth/loginByWeChat", control.LoginByWeChat)       //通过微信小程序注册并拉起微信授权

	router.GET("api/shop/category_list", control.CategoryList)

	router.GET("api/shop_list", control.ShopList) //首页
	router.GET("api/shop_list_nearby", control.ShopListNearby)
	router.GET("api/shop_detail", control.ShopDetail)
	router.GET("api/shop_detail_with_product", control.ShopDetailWithProduct)

	router.GET("api/app/merchant_list", control.MerchantList)
	router.GET("api/app/merchant_detail", control.MerchantDetail)
	router.GET("api/app/merchant_detail_with_shop", control.MerchantDetailWithShop)

	router.GET("api/shop/product_list", control.ProductList)
	router.GET("api/shop/product_detail", control.ProductDetail)
	router.GET("api/shop/bargain_set", control.BargainProduct)  //TODO 获取砍价详情，待实现
	router.POST("api/shop/product_price", control.ProductPrice) //通过sku获取商品价格

	router.GET("api/shop/banner_detail", control.BannerDetail)
	router.GET("api/shop/banner_list", control.BannerList)

	router.GET("api/common/province_list", control.ProvinceList)
	router.GET("api/common/city_list", control.CityList)
	router.GET("api/common/area_list", control.AreaList)

	router.GET("api/shop/product_dynamic", control.ProductList) //购买的人还买了,这个要改

	router.POST("api/app/wechat_pay/callback", control.WechatAppPayCallback) // 微信支付成功后腾讯的回调
}
