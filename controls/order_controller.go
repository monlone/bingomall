package control

import (
	"fmt"
	"bingomall/constant"
	helper "bingomall/helpers"
	"bingomall/helpers/convention"
	model "bingomall/models"
	"bingomall/repositories"
	service "bingomall/services"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"strconv"
)

type createOrderPost struct {
	PayProductList []*model.ProductOrder `json:"payProductList" form:"payProductList"`
	Remark         string                `json:"remark"`
	DeliveryType   int8                  `json:"deliveryType"`
	Calculate      bool                  `json:"calculate"`
	BuyChannel     string                `json:"buyChannel"`
	Prepare        bool                  `json:"prepare"`
}

// 获取 order 列表
// @Summary 获取 order 列表
// @Tags OrderController
// @Accept json
// @Produce json
// @Success 200 {object} model.ListsResponse
// @Router /api/app/order/list [get]
func OrderList(context *gin.Context) {
	userId := helper.GetUserID(context)
	orderService := service.OrderServiceInstance(repositories.OrderRepositoryInstance(helper.GetDBByName(constant.DBMerchant)))

	orderListObject := model.OrderListObject{}
	_ = context.Bind(&orderListObject)
	pageObject := &model.PageObject{}
	_ = context.ShouldBindQuery(pageObject)
	if pageObject.Page == 0 {
		pageObject.Page = 1
	}
	if pageObject.PageSize == 0 {
		pageObject.PageSize = 20
	}

	order := &model.Order{}
	order.UserID = userId
	order.ID = orderListObject.OrderId
	order.Status = pageObject.Status
	orderList := orderService.GetPageByMonth(pageObject.Page, pageObject.PageSize, order, pageObject.Month)

	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: "ok",
			Content: orderList,
		})
}

func OrderListAll(context *gin.Context) {
	var order *model.Order
	orderService := service.OrderServiceInstance(repositories.OrderRepositoryInstance(helper.GetDBByName(constant.DBMerchant)))
	pageStr := context.DefaultQuery("page", "0")
	page, _ := strconv.Atoi(pageStr)
	pageSizeStr := context.DefaultQuery("page_size", constant.PageSize)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	orderList := orderService.GetPage(page, pageSize, order)

	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: "ok",
			Content: orderList,
		})
}

// 获取 order 详情
// @Summary 获取 order 详情
// @Tags OrderController
// @Accept json
// @Produce json
// @Success 200 {object} model.Order
// @Router /api/app/order_detail [get]
func OrderDetail(context *gin.Context) {
	orderService := service.OrderServiceInstance(repositories.OrderRepositoryInstance(helper.GetDBByName(constant.DBMerchant)))
	orderId := context.Query("order_id")
	userId := helper.GetUserID(context)
	order := orderService.GetByOrderIdAndUserId(convention.StringToUint64(orderId), userId)

	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: "ok",
			Content: order,
		})
}

func Statistics(context *gin.Context) {
	orderService := service.OrderServiceInstance(repositories.OrderRepositoryInstance(helper.GetDBByName(constant.DBMerchant)))
	userId := helper.GetUserID(context)
	orders := orderService.Statistics(userId)
	//0:待支付，1：已支付，待回调确认，2：支付成功，已回调确认，3:线下商品商户已经核销，4：已经发货，5：已收货，6：纠纷中，7:已退款，8：已完成
	data := make(map[string]int, 0)
	if len(orders) > 0 {
		for _, v := range orders {
			//待支付
			if v.Status == constant.OrderWaitedForPay {
				data["noPay"] = v.Total
			}
			//待发货
			if v.Status == constant.OrderPaySuccess {
				data["noShipped"] = v.Total
			}
			//已发货，待收货
			if v.Status == constant.OrderShipped {
				data["noConfirm"] = v.Total
			}
			//已收货，待评价
			if v.Status == constant.OrderReceived {
				data["noReview"] = v.Total
			}
		}
	}
	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: "ok",
			Content: data,
		})
}

func LevelCost(userId uint64, productId uint64, cost *model.Cost) {
	userService := service.UserServiceInstance(repositories.UserRepositoryInstance(helper.GetUserDB()))
	user := userService.GetByUserID(userId)
	if user == nil || user.MultiLevel == 0 {
		return
	}
	db := helper.GetDBByName(constant.DBMerchant)
	tx := db.Begin()

	orderService := service.OrderServiceInstance(repositories.OrderRepositoryInstance(tx))
	productService := service.ProductServiceInstance(repositories.ProductRepositoryInstance(tx))
	//levelUserIDs := strings.Split(user.MultiLevel, ",")
	//for _, levelUserID := range levelUserIDs {
	//	order := orderService.GetLevelOrder(convention.StringToUint64(levelUserID))
	//	product := productService.GetByProductId(productId)
	//	if order == nil {
	//		money := float64(product.Price) * float64(product.DiscountLevel)
	//		order = &model.Order{}
	//		order.Pay = uint64(math.Floor(money + 0.5))
	//		order.UserID = userId
	//		err := orderService.Save(order)
	//		if err != nil {
	//			tx.Rollback()
	//			fmt.Println("err update in LevelCost function", helper.Json(levelUserIDs))
	//			break
	//		}
	//	} else {
	//		err := orderService.Update(order)
	//		if err != nil {
	//			tx.Rollback()
	//			fmt.Println("err update in LevelCost function", helper.Json(levelUserIDs))
	//			break
	//		}
	//	}
	//}

	order := orderService.GetLevelOrder(user.MultiLevel)
	product := productService.GetByProductId(productId)
	if order == nil {
		money := float64(product.Price) * float64(product.DiscountLevel)
		order = &model.Order{}
		order.Pay = uint64(math.Floor(money + 0.5))
		order.UserID = userId
		err := orderService.Save(order)
		if err != nil {
			tx.Rollback()
		}
	} else {
		err := orderService.Update(order)
		if err != nil {
			tx.Rollback()
		}
	}
	tx.Commit()
}

func addUserOrderDetail(cost *model.Cost, levelUserID uint64) error {
	db := helper.GetDBByName(constant.DBMerchant)
	tx := db.Begin()
	orderService := service.OrderServiceInstance(repositories.OrderRepositoryInstance(tx))
	order := orderService.GetByOrderId(cost.OrderId)
	shopService := service.ShopServiceInstance(repositories.ShopRepositoryInstance(tx))
	shop := shopService.GetByShopId(cost.ShopId)
	userOrderOrderDetail := &model.UserOrderDetail{}
	userOrderDetailService := service.UserOrderDetailServiceInstance(repositories.UserOrderDetailRepositoryInstance(tx))

	userOrderOrderDetail.Number = cost.Number
	userOrderOrderDetail.Money = cost.Money * cost.Number

	userOrderOrderDetail.MerchantId = shop.MerchantId
	userOrderOrderDetail.OrderId = order.ID
	userOrderOrderDetail.UserID = order.UserID
	userOrderOrderDetail.GoodsType = order.Type
	userOrderOrderDetail.ShopId = cost.ShopId
	userOrderOrderDetail.Type = constant.UserCost
	userOrderOrderDetail.CheckUserID = levelUserID
	err := userOrderDetailService.Save(userOrderOrderDetail)

	return err
}

// 编辑 order
// @Summary  order 核销
// @Tags OrderController
// @Accept json
// @Param check_user_id query string true "店员的user_id"
// @Param money query int64 true "当前核销的金额，为代金券时必传"
// @Param number query int true "当前核销的数量，为实物时必传"
// @Param shop_id query string true "支付的店铺id"
// @Produce json
// @Success 200 {object} model.Order
// @Router /api/order/verification [post]
func VerificationOrder(context *gin.Context) {
	//order表是用户的购买记录，同时也有商户的相关信息存在里面。现在钱的明细和order的明细是分开的。
	//user_order_detail记录了用户的order明细，每一条核销记录都有
	//merchant_order_detail记录了商户的order明细，每一条核销录都有
	//money_log记录了用户、商户与钱相关的明细
	cost := &model.Cost{}
	_ = context.Bind(cost)
	var err error
	if err = context.Bind(cost); err != nil {
		context.JSON(http.StatusUnprocessableEntity, helper.JsonObject{
			Code:    4234,
			Message: helper.StatusText(helper.BindModelErr),
			Content: err,
		})
		return
	}
	userId := helper.GetUserID(context)
	db := helper.GetDBByName(constant.DBMerchant)
	tx := db.Begin()
	orderProductService := service.OrderProductServiceInstance(repositories.OrderProductRepositoryInstance(tx))
	orderService := service.OrderServiceInstance(repositories.OrderRepositoryInstance(tx))
	order := orderService.GetByOrderIdAndUserId(cost.OrderId, userId)
	if order == nil {
		tx.Rollback()
		context.JSON(http.StatusUnprocessableEntity, helper.JsonObject{
			Code:    4234,
			Message: helper.StatusText(helper.BindModelErr),
			Content: err,
		})
		return
	}

	var orderProduct *model.OrderProduct
	orderProduct = orderProductService.GetByOrderProductId(cost.OrderId)

	if orderProduct == nil {
		tx.Rollback()
		context.JSON(http.StatusOK,
			&helper.JsonObject{
				Code:    4204,
				Message: "orderProduct 不存在",
			})
		return
	}
	if orderProduct.Status != constant.OrderPaySuccess {
		tx.Rollback()
		context.JSON(http.StatusOK,
			&helper.JsonObject{
				Code:    4205,
				Message: "order状态不对",
			})
		return
	}

	if orderProduct.Type == constant.OrderRealGoods {
		if cost.Number <= 0 || cost.Number > orderProduct.Number {
			tx.Rollback()
			helper.ErrorLogger.Errorln("Number 非法, cost.Number:", cost.Number, ",orderProduct.Number:", orderProduct.Number)
			context.JSON(http.StatusOK,
				&helper.JsonObject{
					Code:    4224,
					Message: "Number 非法",
				})
			return
		}
	}

	checkUserID := context.PostForm("checkUserID")
	userService := service.UserServiceInstance(repositories.UserRepositoryInstance(helper.GetUserDB()))
	checkUser := userService.GetByUserID(convention.StringToUint64(checkUserID))
	if checkUser == nil {
		helper.ErrorLogger.Errorln("核销店员非法, checkUserID:", checkUserID)
		tx.Rollback()
		context.JSON(http.StatusOK,
			&helper.JsonObject{
				Code:    4226,
				Message: "核销店员非法",
			})
		return
	}
	shopService := service.ShopServiceInstance(repositories.ShopRepositoryInstance(helper.GetDBByName(constant.DBMerchant)))
	shop := shopService.GetByShopId(cost.ShopId)
	if shop == nil {
		tx.Rollback()
		context.JSON(http.StatusOK,
			&helper.JsonObject{
				Code:    4227,
				Message: "核销店铺非法",
			})
		return
	}

	productService := service.ProductServiceInstance(repositories.ProductRepositoryInstance(helper.GetDBByName(constant.DBMerchant)))
	product := productService.GetByProductId(orderProduct.ID)

	// 添加商户记录
	merchantOrderDetail := &model.MerchantOrderDetail{}
	merchantOrderDetailService := service.MerchantOrderDetailServiceInstance(repositories.MerchantOrderDetailRepositoryInstance(tx))

	percent := float64(product.Discount / 100)
	describe := "商家核销自己产品"
	status := constant.MerchantOrderWaitedForLiquidate
	orderShop := shopService.GetByShopId(orderProduct.ShopId)
	baseService := service.BaseServiceInstance(repositories.BaseRepositoryInstance(tx))
	if baseService.IsPlatformOfficial(orderShop) == true {
		describe = "商家代核销平台产品"
		status = constant.MerchantOrderLiquidated
		percent = float64(product.DiscountMerchant / 100)
	}

	totalMoney := float64(0)
	if orderProduct.Type == constant.OrderRealGoods { // 实物商品有个数,扣除相应的个数
		totalMoney = float64(orderProduct.Price * cost.Number)
		merchantOrderDetail.Number = cost.Number
		merchantOrderDetail.Money = uint64(math.Floor(totalMoney*percent + 0.5)) //因为按百分比分钱就有小数了
	} else if orderProduct.Type == constant.OrderVirtualGoods { //虚拟商品没有个数,扣除相应消费的钱
		totalMoney = float64(cost.Money)
		merchantOrderDetail.Money = uint64(math.Floor(totalMoney*percent) + 0.5)
	}

	//添加核销商户的流水
	shopCheck := shopService.GetByShopId(cost.ShopId)
	merchantOrderDetail.Type = constant.MerchantCheckUserPay
	merchantOrderDetail.MerchantId = shopCheck.MerchantId
	merchantOrderDetail.ProductId = orderProduct.ID
	merchantOrderDetail.ShopId = cost.ShopId
	merchantOrderDetail.UserID = userId
	merchantOrderDetail.OrderId = orderProduct.OrderId
	merchantOrderDetail.GoodsType = orderProduct.Type
	merchantOrderDetail.CheckUserID = convention.StringToUint64(checkUserID)
	merchantOrderDetail.Status = uint8(status)
	merchantOrderDetail.Describe = describe
	err = merchantOrderDetailService.Save(merchantOrderDetail)
	if err != nil {
		helper.ErrorLogger.Errorln("4228, DB添加商户记录事务错误, err:", err.Error())
		tx.Rollback()
		context.JSON(http.StatusOK,
			&helper.JsonObject{
				Code:    4228,
				Message: "DB添加商户记录事务错误",
			})
		return
	}

	// 添加平台或者商家流水记录
	// 当shop是平台自己的店或者商户代销售的店时，要添加平台自己的流水,现在平台会卖很多产品，然后让不同的商户去帮忙平台核销，其实就是和菜鸟驿站一样。
	if baseService.IsPlatformOfficial(shop) == true {
		merchantOrderDetailPlatform := &model.MerchantOrderDetail{}
		percentTemp := float64(1 - product.DiscountMerchant/100 - product.DiscountLevel/100 - product.DiscountUser/100)
		merchantOrderDetailPlatform.Money = uint64(totalMoney - totalMoney*percentTemp + 0.5)
		merchantOrderDetailPlatform.Type = constant.MerchantCheckUserPay
		merchantOrderDetailPlatform.MerchantId = shop.MerchantId
		merchantOrderDetailPlatform.ProductId = orderProduct.ID
		merchantOrderDetailPlatform.ShopId = orderProduct.ShopId
		merchantOrderDetailPlatform.UserID = userId
		merchantOrderDetailPlatform.OrderId = orderProduct.OrderId
		merchantOrderDetailPlatform.GoodsType = orderProduct.Type
		merchantOrderDetailPlatform.CheckUserID = convention.StringToUint64(checkUserID)
		merchantOrderDetailPlatform.Status = constant.MerchantOrderWaitedForLiquidate
		merchantOrderDetailPlatform.Describe = "商家代核销平台商品，平台收益"
		merchantOrderDetailPlatform.Number = cost.Number
		err = merchantOrderDetailService.Save(merchantOrderDetailPlatform)
		if err != nil {
			helper.ErrorLogger.Errorln("4228, DB添加商户记录事务错误, err:", err)
			tx.Rollback()
			context.JSON(http.StatusOK,
				&helper.JsonObject{
					Code:    4229,
					Message: "DB添加商户记录事务错误02",
				})
			return
		}
	} else {
		//添加商户记录
		merchantOrderDetailPlatform := &model.MerchantOrderDetail{}
		merchantOrderDetailPlatform.Number = cost.Number
		percentTemp := float64(product.DiscountPlatform / 100)
		merchantOrderDetailPlatform.Money = uint64(totalMoney*percentTemp + 0.5)
		merchantOrderDetailPlatform.Type = constant.MerchantCheckUserPay
		merchantOrderDetailPlatform.MerchantId = shop.MerchantId
		merchantOrderDetailPlatform.ProductId = orderProduct.ID
		merchantOrderDetailPlatform.ShopId = orderProduct.ShopId
		merchantOrderDetailPlatform.UserID = userId
		merchantOrderDetailPlatform.OrderId = orderProduct.OrderId
		merchantOrderDetailPlatform.GoodsType = orderProduct.Type
		merchantOrderDetailPlatform.CheckUserID = convention.StringToUint64(checkUserID)
		merchantOrderDetailPlatform.Status = constant.MerchantOrderWaitedForLiquidate
		merchantOrderDetailPlatform.Describe = "商家代核销自己商品，平台收益"
		err = merchantOrderDetailService.Save(merchantOrderDetailPlatform)
		if err != nil {
			helper.ErrorLogger.Errorln("4228, DB添加商户记录事务错误, err:", err)
			tx.Rollback()
			context.JSON(http.StatusOK,
				&helper.JsonObject{
					Code:    4230,
					Message: "DB添加商户记录事务错误02",
				})
			return
		}
	}

	// 	添加用户记录
	userOrderOrderDetail := &model.UserOrderDetail{}
	userOrderDetailService := service.UserOrderDetailServiceInstance(repositories.UserOrderDetailRepositoryInstance(tx))

	if orderProduct.Type == constant.OrderRealGoods {
		userOrderOrderDetail.Number = cost.Number
		userOrderOrderDetail.Money = orderProduct.Price * uint64(cost.Number)
	} else if orderProduct.Type == constant.OrderVirtualGoods {
		userOrderOrderDetail.Money = cost.Money
	}

	userOrderOrderDetail.MerchantId = shop.MerchantId
	userOrderOrderDetail.ProductId = orderProduct.ID
	userOrderOrderDetail.OrderId = orderProduct.OrderId
	userOrderOrderDetail.UserID = userId
	userOrderOrderDetail.GoodsType = orderProduct.Type
	userOrderOrderDetail.ShopId = orderProduct.ShopId
	userOrderOrderDetail.Type = constant.UserCost
	userOrderOrderDetail.CheckUserID = convention.StringToUint64(checkUserID)
	err = userOrderDetailService.Save(userOrderOrderDetail)
	if err != nil {
		tx.Rollback()
		helper.ErrorLogger.Errorln("4229, DB添加用户记录事务错误, err:", err.Error())
		context.JSON(http.StatusOK,
			&helper.JsonObject{
				Code:    4231,
				Message: "DB添加用户记录事务错误",
			})
		return
	}

	if orderProduct.Type == constant.OrderRealGoods {
		orderProduct.Number = orderProduct.Number - cost.Number
		if orderProduct.Number == 0 {
			orderProduct.Status = constant.OrderVerified
		}
	} else if orderProduct.Type == constant.OrderVirtualGoods {
		orderProduct.Money = orderProduct.Money - cost.Money
		if orderProduct.Money == 0 {
			orderProduct.Status = constant.OrderVerified
		}
	}

	err = orderProductService.Update(orderProduct)
	if err != nil {
		tx.Rollback()
		helper.ErrorLogger.Errorln("4231, DB更新order记录事务错误, err:", err.Error())
		context.JSON(http.StatusOK,
			&helper.JsonObject{
				Code:    4232,
				Message: "DB更新order记录事务错误",
			})
		return
	}

	//TODO 以后改到队列不能开协程来处理，开了协程更新不了数据库，会有这个bug,这个要时间修复。
	//go profitSharing(order, float64(cost.Money*uint64(cost.Number)))
	helper.ServiceLogger.Println("in VerificationOrder func, order:", helper.Json(order),
		",totalMoney:", totalMoney, ",cost.shopId:", cost.ShopId)
	err = profitSharing(tx, order, orderProduct, float64(totalMoney), cost.ShopId)

	if err != nil {
		helper.ErrorLogger.Errorln("4230, DB更新order事务错误, err:", err.Error())
		tx.Rollback()
		context.JSON(http.StatusOK,
			&helper.JsonObject{
				Code:    4233,
				Message: "DB更新order事务错误",
			})
		return
	} else {
		tx.Commit()
	}

	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: helper.StatusText(helper.SaveStatusOK),
		})
}

func CreateOrderPrepare(context *gin.Context) {
	userId := helper.GetUserID(context)
	db := helper.GetDBByName(constant.DBMerchant)
	order := &model.Order{}
	orderReturn := &model.OrderReturn{}
	postData := &createOrderPost{}
	err := context.Bind(postData)
	if err != nil {
		context.JSON(http.StatusOK, helper.JsonObject{
			Code:    1120,
			Message: "绑定参数错误",
			Content: err,
		})
		return
	}
	var deliveryPrice uint64
	deliveryPrice = 1
	shopIds := make([]uint64, 1)
	var productIds []uint64

	for _, v := range postData.PayProductList {
		productIds = append(productIds, v.ProductId)
	}

	var productList []*model.Product
	productService := service.ProductServiceInstance(repositories.ProductRepositoryInstance(db))
	productList = productService.GetListByProductIds(productIds)

	for _, p := range productList {
		for _, v := range postData.PayProductList {
			if v.ProductId == p.ID {
				score := uint64(1) //TODO 这个要改的
				order.TotalAmount += p.Price * v.Number
				order.Pay = order.TotalAmount - score
				order.LogisticsAmount += deliveryPrice * v.Number
				shopIds = append(shopIds, p.ShopId)
			}
		}
	}

	order.OutTradeNo = helper.GenerateId32()
	order.Type = 0 //TODO 订单类型要改，先占个位，还没想好要怎么定义
	order.UserID = userId
	order.PrepayID = ""
	order.Platform = "wechatMini"
	order.Score = 0 //支付的积分

	orderReturn.Order = *order
	orderReturn.LogisticsAmount = 100
	orderReturn.IsNeedLogistics = true
	orderReturn.CouponList = nil

	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: helper.StatusText(helper.SaveStatusOK),
			Content: orderReturn,
		})
}

//创建定单
func CreateOrder(context *gin.Context) {
	userId := helper.GetUserID(context)
	db := helper.GetDBByName(constant.DBMerchant).Begin()
	order := &model.Order{}
	orderReturn := &model.OrderReturn{}
	postData := &createOrderPost{}
	err := context.Bind(postData)
	if err != nil {
		context.JSON(http.StatusOK, helper.JsonObject{
			Code:    1120,
			Message: "绑定参数错误",
			Content: err,
		})
		return
	}
	var deliveryPrice uint64
	deliveryPrice = 1
	shopIds := make([]uint64, 1)
	var productIds []uint64
	var shoppingCartIds []uint64

	for _, v := range postData.PayProductList {
		productIds = append(productIds, v.ProductId)
		shoppingCartIds = append(shoppingCartIds, v.ShoppingCartId)
	}

	shoppingCartService := service.ShoppingCartServiceInstance(repositories.ShoppingCartRepositoryInstance(db))
	shoppingCartList := shoppingCartService.GetShoppingCartByUserIDCartIds(userId, shoppingCartIds)

	for _, sc := range shoppingCartList {
		for _, v := range postData.PayProductList {
			if v.ShoppingCartId == sc.ID {
				score := uint64(0) //TODO 这个要改的
				order.TotalAmount += sc.Sku.Price * v.Number
				order.Pay = order.TotalAmount - score
				order.LogisticsAmount += deliveryPrice * v.Number
				shopIds = append(shopIds, sc.Product.ShopId)
			}
		}
	}

	orderService := service.OrderServiceInstance(repositories.OrderRepositoryInstance(db))

	order.OutTradeNo = helper.GenerateId32()
	order.Type = 0 //TODO 订单类型要改，先占个位，还没想好要怎么定义
	order.UserID = userId
	order.PrepayID = ""
	order.Platform = postData.BuyChannel
	order.Score = 0 //支付的积分

	err = orderService.Save(order)
	if err != nil {
		fmt.Println("err:::", err)
	}
	var orderProductList []model.OrderProduct
	temp := model.OrderProduct{}
	for _, p := range shoppingCartList {
		for _, postProduct := range postData.PayProductList {
			if postProduct.ShoppingCartId == p.ID {
				temp.ShopId = p.Product.ShopId
				temp.Money = p.Sku.Price * postProduct.Number
				temp.OrderId = order.ID
				temp.ProductId = postProduct.ProductId
				temp.Price = p.Sku.Price
				temp.Number = postProduct.Number //TODO 要完善
				//temp.CreatedAt = datetime.DateTime.CurrentTime()
				temp.SkuId = p.Sku.ID
				orderProductList = append(orderProductList, temp)
			}
		}
	}

	orderProductService := service.OrderProductServiceInstance(repositories.OrderProductRepositoryInstance(db))
	err = orderProductService.BatchInsert(orderProductList)
	if err != nil {
		fmt.Println("err in orderProductService BatchInsert:", err.Error())
	}

	db.Commit()
	orderReturn.Order = *order
	orderReturn.LogisticsAmount = 100
	orderReturn.IsNeedLogistics = true
	orderReturn.CouponList = nil

	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: helper.StatusText(helper.SaveStatusOK),
			Content: orderReturn,
		})
}
