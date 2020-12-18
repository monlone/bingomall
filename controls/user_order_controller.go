package control

import (
	"bingomall/constant"
	helper "bingomall/helpers"
	"bingomall/helpers/convention"
	model "bingomall/models"
	"bingomall/repositories"
	service "bingomall/services"
	"bingomall/system"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// 获取 userOrderDetail 列表
// @Summary 获取 userOrderDetail 列表
// @Tags UserOrderDetailController
// @Accept json
// @Produce json
// @Success 200 {object} model.UserOrderListsResponse
// @Router /api/app/order/user_list [get]
func UserOrderDetailList(context *gin.Context) {
	claims, ok := context.Get("claims")
	if !ok {
		context.JSON(http.StatusOK, helper.JsonObject{
			Code:    1105,
			Message: "token错误",
		})
		context.Abort()
	}
	userOrderDetailService := service.UserOrderDetailServiceInstance(repositories.UserOrderDetailRepositoryInstance(helper.GetDBByName(constant.DBMerchant)))
	pageObj := model.PageObject{}
	_ = context.Bind(&pageObj)
	userInfo := claims.(*system.CustomClaims)
	userOrderDetail := &model.UserOrderDetail{}
	userOrderDetail.UserID = userInfo.ID
	month := context.DefaultQuery("month", "")
	if pageObj.PageSize == 0 {
		pageSize, _ := strconv.Atoi(constant.PageSize)
		pageObj.PageSize = pageSize
	}
	userOrderDetailList := userOrderDetailService.GetPageByMonth(pageObj.Page, pageObj.PageSize, userOrderDetail, month)

	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: "ok",
			Content: userOrderDetailList,
		})
}

func UserOrderDetailListAll(context *gin.Context) {
	var userOrderDetail *model.UserOrderDetail
	userOrderDetailService := service.UserOrderDetailServiceInstance(repositories.UserOrderDetailRepositoryInstance(helper.GetDBByName(constant.DBMerchant)))
	pageStr := context.DefaultQuery("page", "0")
	page, _ := strconv.Atoi(pageStr)
	pageSizeStr := context.DefaultQuery("page_size", constant.PageSize)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	userOrderDetailList := userOrderDetailService.GetPage(page, pageSize, userOrderDetail)

	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: "ok",
			Content: userOrderDetailList,
		})
}

// 获取 userOrderDetail 详情
// @Summary 获取 userOrderDetail 详情
// @Tags UserOrderDetailController
// @Accept json
// @Produce json
// @Success 200 {object} model.UserOrderDetail
// @Router /api/app/userOrderDetail_detail [get]
func UserOrderDetailDetail(context *gin.Context) {
	userOrderDetailService := service.UserOrderDetailServiceInstance(repositories.UserOrderDetailRepositoryInstance(helper.GetDBByName(constant.DBMerchant)))
	userOrderDetailId := context.Query("userOrderDetail_id")
	userOrderDetail := userOrderDetailService.GetByUserOrderDetailID(convention.StringToUint64(userOrderDetailId))

	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: "ok",
			Content: userOrderDetail,
		})
}

// 编辑 userOrderDetail
// @Summary  userOrderDetail 编辑
// @Tags UserOrderDetailController
// @Accept json
// @Produce json
// @Success 200 {object} model.UserOrderDetail
// @Router /api/userOrderDetail/verification [post]
func SaveUserOrderDetail(context *gin.Context) {
	userOrderDetail := &model.UserOrderDetail{}
	if err := context.Bind(userOrderDetail); err == nil {
		userOrderDetailService := service.UserOrderDetailServiceInstance(repositories.UserOrderDetailRepositoryInstance(helper.GetDBByName(constant.DBMerchant)))
		err := userOrderDetailService.SaveOrUpdate(userOrderDetail)
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
