package control

import (
	"fmt"
	"bingomall/constant"
	helper "bingomall/helpers"
	"bingomall/helpers/convention"
	model "bingomall/models"
	"bingomall/repositories"
	service "bingomall/services"
	"bingomall/system"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

var merchantDB *gorm.DB

func init() {
	merchantDB = helper.GetDBByName(constant.DBMerchant)
}

// 获取 merchantOrderDetail 列表
// @Summary 获取 merchantOrderDetail 列表
// @Tags MerchantOrderDetailController
// @Accept json
// @Produce json
// @Success 200 {object} model.MerchantOrderListsResponse
// @Router /api/app/order/merchant_list [get]
func MerchantOrderDetailList(context *gin.Context) {
	claims, ok := context.Get("claims")
	if !ok {
		context.JSON(http.StatusOK, helper.JsonObject{
			Code:    1105,
			Message: "token错误",
		})
		context.Abort()
	}
	merchantOrderDetailService := service.MerchantOrderDetailServiceInstance(repositories.MerchantOrderDetailRepositoryInstance(merchantDB))
	userInfo := claims.(*system.CustomClaims)
	merchantOrderDetail := &model.MerchantOrderDetail{}
	MerchantOrderDetailPage := &model.MerchantOrderDetailPage{}
	_ = context.Bind(MerchantOrderDetailPage)
	merchantOrderDetail.Status = uint8(MerchantOrderDetailPage.Status)

	//通过user_shop表取shopIds lee
	//userShopService := service.UserShopServiceInstance(repositories.UserShopRepositoryInstance(merchantDB))
	//shop := userShopService.GetByUserID(userInfo.ID)
	//if len(shop) == 0 {
	//	context.JSON(http.StatusOK, helper.JsonObject{
	//		Code:    "1106",
	//		Message: "商店不存在",
	//	})
	//	return
	//}
	//
	//shopIds := helper.ModelObjectToSlice(shop, "ShopId")

	//通过merchant表取shopIds
	merchantService := service.MerchantServiceInstance(repositories.MerchantRepositoryInstance(merchantDB))
	merchant := merchantService.GetByUserID(userInfo.ID)
	MerchantIds := helper.ModelObjectToSlice(merchant, "MerchantId")
	shopService := service.ShopServiceInstance(repositories.ShopRepositoryInstance(merchantDB))
	shopList := shopService.GetShopsByMerchantIds(MerchantIds)
	shopIds := helper.ModelObjectToSlice(shopList, "ShopId")
	//通过merchant表取shopIds结束

	if MerchantOrderDetailPage.PageSize == 0 {
		pageSize, _ := strconv.Atoi(constant.PageSize)
		MerchantOrderDetailPage.PageSize = pageSize
	}

	merchantOrderDetailList := merchantOrderDetailService.GetPageByShopIds(MerchantOrderDetailPage, merchantOrderDetail, shopIds)

	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: "ok",
			Content: merchantOrderDetailList,
		})
}

// 获取 merchantOrderDetail 列表
// @Summary 获取 merchantOrderDetail 列表
// @Tags MerchantOrderDetailController
// @Accept json
// @Produce json
// @Success 200 {object} model.MerchantOrderListsResponse
// @Router /api/app/order/liquidate [post]
func Liquidate(context *gin.Context) {
	claims, ok := context.Get("claims")
	if !ok {
		context.JSON(http.StatusOK, helper.JsonObject{
			Code:    1105,
			Message: "token错误",
		})
		context.Abort()
	}
	dbTransaction := merchantDB.Begin()
	merchantOrderDetailService := service.MerchantOrderDetailServiceInstance(repositories.MerchantOrderDetailRepositoryInstance(dbTransaction))
	userInfo := claims.(*system.CustomClaims)
	merchantOrderDetail := &model.MerchantOrderDetail{}
	MerchantOrderDetailPage := &model.MerchantOrderDetailPage{}
	_ = context.Bind(MerchantOrderDetailPage)
	merchantOrderDetail.Status = uint8(MerchantOrderDetailPage.Status)
	//userShopService := service.UserShopServiceInstance(repositories.UserShopRepositoryInstance(dbTransaction))
	//shop := userShopService.GetByUserID(userInfo.ID)
	//if len(shop) == 0 {
	//	context.JSON(http.StatusOK, helper.JsonObject{
	//		Code:    "1106",
	//		Message: "商店不存在",
	//	})
	//	return
	//}
	//
	//shopIds := helper.ModelObjectToSlice(shop, "ShopId")

	//TODO 要判断最近7天，小于7天不能结算
	merchantService := service.MerchantServiceInstance(repositories.MerchantRepositoryInstance(dbTransaction))
	merchant := merchantService.GetByUserID(userInfo.ID)

	MerchantIds := helper.ModelObjectToSlice(merchant, "MerchantId")
	shopService := service.ShopServiceInstance(repositories.ShopRepositoryInstance(merchantDB))
	shopList := shopService.GetShopsByMerchantIds(MerchantIds)
	shopIds := helper.ModelObjectToSlice(shopList, "ShopId")

	merchantIds := helper.ModelObjectToSlice(merchant, "MerchantId")
	liquidateTotal, _ := merchantOrderDetailService.LiquidateTotal(merchantIds, shopIds)
	err := merchantOrderDetailService.Liquidate(merchantIds, shopIds)

	if err != nil {
		dbTransaction.Rollback()
		fmt.Println("1107, 结算失败, err:", err)
		context.JSON(http.StatusOK, helper.JsonObject{
			Code:    1107,
			Message: "结算失败！",
		})
		return
	}

	moneyLog := service.MoneyLogServiceInstance(repositories.MoneyLogRepositoryInstance(dbTransaction))
	moneyLogData := &model.MoneyLog{}
	moneyLogData.UserID = userInfo.ID
	moneyLogData.Type = constant.MoneyLogMerchant
	moneyLogData.RelationUserID = merchantIds[0]
	moneyLogData.Cost = liquidateTotal.TotalMoney
	moneyLogData.Describe = "商户结算"

	err = moneyLog.Save(moneyLogData)
	if err != nil {
		dbTransaction.Rollback()
	} else {
		dbTransaction.Commit()
	}

	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: "ok",
			Content: "结算成功",
		})
}

// 获取 merchantOrderDetail 列表
// @Summary 获取 merchantOrderDetail 列表
// @Tags MerchantOrderDetailController
// @Accept json
// @Produce json
// @Success 200 {object} model.MerchantOrderListsResponse
// @Router /api/app/order/liquidate_total [get]
func LiquidateTotal(context *gin.Context) {
	claims, ok := context.Get("claims")
	if !ok {
		context.JSON(http.StatusOK, helper.JsonObject{
			Code:    1105,
			Message: "token错误",
		})
		context.Abort()
	}
	merchantOrderDetailService := service.MerchantOrderDetailServiceInstance(repositories.MerchantOrderDetailRepositoryInstance(merchantDB))
	userInfo := claims.(*system.CustomClaims)
	//userShopService := service.UserShopServiceInstance(repositories.UserShopRepositoryInstance(merchantDB))
	//shop := userShopService.GetByUserID(userInfo.ID)
	//if len(shop) == 0 {
	//	context.JSON(http.StatusOK, helper.JsonObject{
	//		Code:    "1106",
	//		Message: "商店不存在",
	//	})
	//	return
	//}
	//shopIds := helper.ModelObjectToSlice(shop, "ShopId")
	merchantService := service.MerchantServiceInstance(repositories.MerchantRepositoryInstance(helper.GetDBByName(constant.DBMerchant)))
	merchant := merchantService.GetByUserID(userInfo.ID)
	MerchantIds := helper.ModelObjectToSlice(merchant, "MerchantId")
	shopService := service.ShopServiceInstance(repositories.ShopRepositoryInstance(helper.GetDBByName(constant.DBMerchant)))
	shopList := shopService.GetShopsByMerchantIds(MerchantIds)
	shopIds := helper.ModelObjectToSlice(shopList, "ShopId")

	merchantIds := helper.ModelObjectToSlice(merchant, "MerchantId")
	liquidateTotal, _ := merchantOrderDetailService.LiquidateTotal(merchantIds, shopIds)

	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: "ok",
			Content: liquidateTotal,
		})
}

// 获取 merchantEmployeeOrderDetail 列表
// @Summary 获取 merchantEmployeeOrderDetail 列表
// @Tags MerchantOrderDetailController
// @Accept json
// @Produce json
// @Success 200 {object} model.MerchantOrderListsResponse
// @Router /api/app/order/merchant_employee_list [get]
func MerchantEmployeeOrderDetailList(context *gin.Context) {
	claims, ok := context.Get("claims")
	if !ok {
		context.JSON(http.StatusOK, helper.JsonObject{
			Code:    1105,
			Message: "token错误",
		})
		context.Abort()
	}
	merchantOrderDetailService := service.MerchantOrderDetailServiceInstance(repositories.MerchantOrderDetailRepositoryInstance(helper.GetDBByName(constant.DBMerchant)))
	pageStr := context.DefaultQuery("page", "0")
	page, _ := strconv.Atoi(pageStr)
	pageSizeStr := context.DefaultQuery("page_size", constant.PageSize)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	userInfo := claims.(*system.CustomClaims)
	merchantOrderDetail := &model.MerchantOrderDetail{}
	merchantOrderDetail.CheckUserID = userInfo.ID
	merchantOrderDetail.ShopId = convention.StringToUint64(context.DefaultQuery("shop_id", ""))
	merchantOrderDetailList := merchantOrderDetailService.GetPageByCheckUserID(page, pageSize, merchantOrderDetail)

	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: "ok",
			Content: merchantOrderDetailList,
		})
}

func MerchantOrderDetailListAll(context *gin.Context) {
	var merchantOrderDetail *model.MerchantOrderDetail
	merchantOrderDetailService := service.MerchantOrderDetailServiceInstance(repositories.MerchantOrderDetailRepositoryInstance(helper.GetDBByName(constant.DBMerchant)))
	pageStr := context.DefaultQuery("page", "0")
	page, _ := strconv.Atoi(pageStr)
	pageSizeStr := context.DefaultQuery("page_size", constant.PageSize)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	merchantOrderDetailList := merchantOrderDetailService.GetPage(page, pageSize, merchantOrderDetail)

	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: "ok",
			Content: merchantOrderDetailList,
		})
}

// 获取 merchantOrderDetail 详情
// @Summary 获取 merchantOrderDetail 详情
// @Tags MerchantOrderDetailController
// @Accept json
// @Produce json
// @Success 200 {object} model.MerchantOrderDetail
// @Router /api/app/merchantOrderDetail_detail [get]
func MerchantOrderDetailDetail(context *gin.Context) {
	merchantOrderDetailService := service.MerchantOrderDetailServiceInstance(repositories.MerchantOrderDetailRepositoryInstance(helper.GetDBByName(constant.DBMerchant)))
	merchantOrderDetailId := context.Query("merchantOrderDetail_id")
	merchantOrderDetail := merchantOrderDetailService.GetByMerchantOrderDetailID(convention.StringToUint64(merchantOrderDetailId))

	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: "ok",
			Content: merchantOrderDetail,
		})
}

// 编辑 merchantOrderDetail
// @Summary  merchantOrderDetail 编辑
// @Tags MerchantOrderDetailController
// @Accept json
// @Produce json
// @Success 200 {object} model.MerchantOrderDetail
// @Router /api/merchantOrderDetail/verification [post]
func SaveMerchantOrderDetail(context *gin.Context) {
	merchantOrderDetail := &model.MerchantOrderDetail{}
	if err := context.Bind(merchantOrderDetail); err == nil {
		merchantOrderDetailService := service.MerchantOrderDetailServiceInstance(repositories.MerchantOrderDetailRepositoryInstance(helper.GetDBByName(constant.DBMerchant)))
		err := merchantOrderDetailService.SaveOrUpdate(merchantOrderDetail)
		if err == nil {
			context.JSON(http.StatusOK,
				&helper.JsonObject{
					Code:    0,
					Message: helper.StatusText(helper.SaveStatusOK),
				})
			return
		} else {
			context.JSON(http.StatusOK,
				&helper.JsonObject{
					Code:    4204,
					Message: err.Error(),
				})
			return
		}
	} else {
		context.JSON(http.StatusUnprocessableEntity, helper.JsonObject{
			Code:    4203,
			Message: helper.StatusText(helper.BindModelErr),
			Content: err,
		})
	}
}
