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
	"net/http"
	"strconv"
)

// 获取 shop 列表
// @Summary 获取 shop 列表
// @Tags ShopController
// @Accept json
// @Produce json
// @Success 200 {object} model.Shop
// @Router /api/app/shop_list [get]
func ShopList(context *gin.Context) {
	shop := &model.Shop{}
	shopService := service.ShopServiceInstance(repositories.ShopRepositoryInstance(helper.GetDBByName(constant.DBMerchant)))
	pageStr := context.DefaultQuery("page", "0")
	page, _ := strconv.Atoi(pageStr)
	pageSizeStr := context.DefaultQuery("page_size", constant.PageSize)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	shop.Title = context.DefaultQuery("title", "")
	shopList := shopService.GetPage(page, pageSize, shop)

	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: "ok",
			Content: shopList,
		})
}

// 获取附近 shop 列表
// @Summary 获取附近 shop 列表
// @Tags ShopController
// @Accept json
// @Produce json
// @Success 200 {object} model.Shop
// @Router /api/app/shop_list_nearby [get]
func ShopListNearby(context *gin.Context) {
	shop := &model.ShopDetailDistance{}
	if err := context.Bind(shop); err == nil {
		shopService := service.ShopServiceInstance(repositories.ShopRepositoryInstance(helper.GetDBByName(constant.DBMerchant)))
		pageStr := context.DefaultQuery("page", "0")
		page, _ := strconv.Atoi(pageStr)
		pageSizeStr := context.DefaultQuery("page_size", constant.PageSize)
		pageSize, _ := strconv.Atoi(pageSizeStr)
		shopList := shopService.ShopListNearby(page, pageSize, shop)

		context.JSON(http.StatusOK,
			&helper.JsonObject{
				Code:    0,
				Message: "ok",
				Content: shopList,
			})
	} else {
		context.JSON(http.StatusOK,
			&helper.JsonObject{
				Code:    4403,
				Message: "bind shop error",
				Content: nil,
			})
	}

}

// 编辑 shop
// @Summary  shop 编辑
// @Tags ShopController
// @Accept json
// @Produce json
// @Success 200 {object} model.Shop
// @Router /api/shop/save [post]
func SaveShop(context *gin.Context) {
	shop := &model.Shop{}
	if err := context.Bind(shop); err == nil {
		shopService := service.ShopServiceInstance(repositories.ShopRepositoryInstance(helper.GetDBByName(constant.DBMerchant)))
		err := shopService.SaveOrUpdate(shop)
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
		fmt.Println("fff:", helper.Json(shop), "err.Error():", err.Error())
		context.JSON(http.StatusUnprocessableEntity, helper.JsonObject{
			Code:    helper.BindModelErr,
			Message: helper.StatusText(helper.BindModelErr),
			Content: err.Error(),
		})
	}
}

// 获取 shop（店铺） 详情
// @Summary  获取 shop（店铺） 详情
// @Tags ShopController
// @Accept json
// @Produce json
// @Param shop_id query string true "商户列表的shop_id"
// @Success 200 {object} model.Product
// @Router /api/shop_detail [get]
func ShopDetail(context *gin.Context) {
	shopService := service.ShopServiceInstance(repositories.ShopRepositoryInstance(helper.GetDBByName(constant.DBMerchant)))
	shopId := context.Query("shop_id")
	shop := shopService.GetByShopId(convention.StringToUint64(shopId))
	_ = shopService.AddAttention(shop)
	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: "ok",
			Content: shop,
		})
}

// 获取带 product 的店铺详情
// @Summary 获取带 product 的店铺详情
// @Tags ShopController
// @Accept json
// @Produce json
// @Param shop_id query string true "商户列表的shop_id"
// @Success 200 {object} model.Product
// @Router /api/shop_detail_with_product [get]
func ShopDetailWithProduct(context *gin.Context) {
	shopService := service.ShopServiceInstance(repositories.ShopRepositoryInstance(helper.GetDBByName(constant.DBMerchant)))
	pageStr := context.DefaultQuery("page", "0")
	page, _ := strconv.Atoi(pageStr)
	pageSizeStr := context.DefaultQuery("page_size", constant.PageSize)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	shopId := context.Query("shop_id")
	shop := &model.Shop{}
	shop.ID = convention.StringToUint64(shopId)
	_ = shopService.AddAttention(shop)
	productList := shopService.ProductList(page, pageSize, shopId)

	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: "ok",
			Content: productList,
		})
}
