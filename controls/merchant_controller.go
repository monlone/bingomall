package control

import (
	"bingomall/constant"
	helper "bingomall/helpers"
	"bingomall/helpers/convention"
	model "bingomall/models"
	"bingomall/repositories"
	service "bingomall/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// 获取 merchant 列表
// @Summary 获取 merchant 列表
// @Tags MerchantController
// @Accept json
// @Produce json
// @Success 200 {object} model.Merchant
// @Router /api/merchant_list [get]
func MerchantList(context *gin.Context) {
	var merchant *model.Merchant
	merchantService := service.MerchantServiceInstance(repositories.MerchantRepositoryInstance(helper.GetDBByName(constant.DBMerchant)))
	pageStr := context.DefaultQuery("page", "0")
	page, _ := strconv.Atoi(pageStr)
	pageSizeStr := context.DefaultQuery("page_size", constant.PageSize)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	merchantList := merchantService.GetPage(page, pageSize, merchant)

	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: "ok",
			Content: merchantList,
		})
}

// 添加 merchant
// @Summary 添加 merchant
// @Tags MerchantController
// @Accept json
// @Produce json
// @Success 200 {object} model.Merchant
// @Router /api/merchant/save [post]
func SaveMerchant(context *gin.Context) {
	merchant := &model.Merchant{}
	if err := context.Bind(merchant); err == nil {
		merchantService := service.MerchantServiceInstance(repositories.MerchantRepositoryInstance(helper.GetDBByName(constant.DBMerchant)))
		err := merchantService.SaveOrUpdate(merchant)
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

// 获取 merchant（商户） 详情
// @Summary  获取 merchant（商户） 详情
// @Tags MerchantController
// @Accept json
// @Produce json
// @Success 200 {object} model.Product
// @Router /api/merchant_detail [get]
func MerchantDetail(context *gin.Context) {
	merchantService := service.MerchantServiceInstance(repositories.MerchantRepositoryInstance(helper.GetDBByName(constant.DBMerchant)))
	merchantId := context.Query("merchant_id")
	merchant := merchantService.GetByMerchantId(convention.StringToUint64(merchantId))

	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: "ok",
			Content: merchant,
		})
}

// 获取 product 列表
// @Summary 获取 product 列表
// @Tags MerchantController
// @Accept json
// @Produce json
// @Success 200 {object} model.Product
// @Router /api/merchant_detail_with_shop [get]
func MerchantDetailWithShop(context *gin.Context) {
	merchantService := service.MerchantServiceInstance(repositories.MerchantRepositoryInstance(helper.GetDBByName(constant.DBMerchant)))
	pageStr := context.DefaultQuery("page", "0")
	page, _ := strconv.Atoi(pageStr)
	pageSizeStr := context.DefaultQuery("page_size", constant.PageSize)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	merchantId := context.Query("merchant_id")
	productList := merchantService.ShopList(page, pageSize, merchantId)

	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: "ok",
			Content: productList,
		})
}
