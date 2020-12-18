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

// 获取 banner 轮播图
// @Summary 获取 banner 轮播图,传 banner_id 按 banner_id 查，不传就查平台的
// @Tags BannerController
// @Accept json
// @Produce json
// @Param shop_id        query string false "产品记录id"
// @Success 200 {object} model.Banner
// @Router /api/banner/list [get]
func BannerList(context *gin.Context) {
	var banner []*model.Banner
	shopId := context.DefaultPostForm("shop_id", "1")
	bannerService := service.BannerServiceInstance(repositories.BannerRepositoryInstance(helper.GetDBByName(constant.DBMerchant)))
	banner = bannerService.GetAllByShopId(convention.StringToUint64(shopId))

	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: "ok",
			Content: banner,
		})
}

// 获取 banner 列表
// @Summary 获取 banner 列表
// @Tags BannerController
// @Accept json
// @Produce json
// @Success 200 {object} model.ListsResponse
// @Router /api/banner_list_all [get]
func BannerListAll(context *gin.Context) {
	var banner *model.Banner
	bannerService := service.BannerServiceInstance(repositories.BannerRepositoryInstance(helper.GetDBByName(constant.DBMerchant)))
	pageStr := context.DefaultQuery("page", "0")
	page, _ := strconv.Atoi(pageStr)
	pageSizeStr := context.DefaultQuery("page_size", constant.PageSize)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	bannerList := bannerService.GetPage(page, pageSize, banner)

	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: "ok",
			Content: bannerList,
		})
}

// 获取 banner 详情
// @Summary  获取 banner 详情
// @Tags BannerController
// @Accept json
// @Produce json
// @Param banner_id query string true "列表的banner_id"
// @Success 200 {object} model.Banner
// @Router /api/app/banner_detail [get]
func BannerDetail(context *gin.Context) {
	bannerService := service.BannerServiceInstance(repositories.BannerRepositoryInstance(helper.GetDBByName(constant.DBMerchant)))
	bannerId := context.Query("banner_id")
	banner := bannerService.GetByBannerID(convention.StringToUint64(bannerId))

	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: "ok",
			Content: banner,
		})
}

// 编辑 banner
// @Summary  banner 编辑
// @Tags BannerController
// @Accept json
// @Produce json
// @Success 200 {object} model.Banner
// @Router /api/banner/save [post]
func SaveBanner(context *gin.Context) {
	banner := &model.Banner{}
	if err := context.Bind(banner); err == nil {
		bannerService := service.BannerServiceInstance(repositories.BannerRepositoryInstance(helper.GetDBByName(constant.DBMerchant)))
		err := bannerService.SaveOrUpdate(banner)
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

// 删除 banner
// @Summary 删除 banner
// @Tags BannerController
// @Accept json
// @Produce json
// @Success 200 {object} model.Banner
// @Router /api/banner/remove [delete]
func RemoveBanner(context *gin.Context) {
	bannerId := context.Query("banner_id")
	bannerService := service.BannerServiceInstance(repositories.BannerRepositoryInstance(helper.GetDBByName(constant.DBMerchant)))
	err := bannerService.DeleteByBannerID(convention.StringToUint64(bannerId))
	if err == nil {
		context.JSON(http.StatusOK,
			&helper.JsonObject{
				Code:    0,
				Message: helper.StatusText(helper.SaveStatusOK),
			})
	} else {
		context.JSON(http.StatusOK,
			&helper.JsonObject{
				Code:    4208,
				Message: err.Error(),
			})
	}
}
