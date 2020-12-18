package control

import (
	"fmt"
	"bingomall/constant"
	helper "bingomall/helpers"
	model "bingomall/models"
	"bingomall/repositories"
	service "bingomall/services"
	"github.com/gin-gonic/gin"

	"net/http"
)

// 获取全部
// @Summary 获取全部商户结算列表
// @Tags MoneyLogController
// @Accept json
// @Produce json
// @Success 200 {object} helpers.JsonObject
// @Router /api/money_log/list_all [get]
func MerchantMoneyLog(context *gin.Context) {
	pageObject := &model.PageObject{}
	err := context.Bind(pageObject)
	if err != nil {
		context.JSON(http.StatusOK,
			&helper.JsonObject{
				Code:    4401,
				Message: "绑定参数错误",
			})
		return
	}

	moneyLogService := service.MoneyLogServiceInstance(repositories.MoneyLogRepositoryInstance(helper.GetDBByName(constant.DBMerchant)))
	moneyLog := &model.MoneyLog{}
	fmt.Println("moneyLog:", helper.Json(moneyLog), "pageObject:", helper.Json(pageObject))
	dataList := moneyLogService.GetPageByMonthWithMerchant(pageObject, moneyLog)
	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: "ok",
			Content: dataList,
		})
}

// 获取当前商户的全部结算列表
// @Summary 获取当前商户结算列表
// @Tags MoneyLogController
// @Accept json
// @Produce json
// @Success 200 {object} helpers.JsonObject
// @Router /api/app/money_log/list [get]
func MoneyLogList(context *gin.Context) {
	userId := helper.GetUserID(context)
	if userId == 0 {
		context.JSON(http.StatusOK,
			&helper.JsonObject{
				Code:    4402,
				Message: "获取用户id错误",
			})
		return
	}
	pageObject := &model.PageObject{}
	err := context.ShouldBindQuery(pageObject)
	if err != nil {
		context.JSON(http.StatusOK,
			&helper.JsonObject{
				Code:    4401,
				Message: "绑定参数错误",
			})
		return
	}

	moneyLogService := service.MoneyLogServiceInstance(repositories.MoneyLogRepositoryInstance(helper.GetDBByName(constant.DBMerchant)))
	moneyLog := &model.MoneyLog{}
	moneyLog.UserID = userId
	dataList := moneyLogService.GetPageByMonth(pageObject, moneyLog)
	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: "ok",
			Content: dataList,
		})
}

// 更新商户结算状态
// @Summary 更新商户结算状态
// @Tags MoneyLogController
// @Accept json
// @Produce json
// @Success 200 {object} helpers.JsonObject
// @Router /api/money_log/update [post]
func UpdateMoneyLog(context *gin.Context) {
	moneyLogData := &model.MoneyLog{}
	err := context.Bind(moneyLogData)
	if err != nil {
		context.JSON(http.StatusOK,
			&helper.JsonObject{
				Code:    4401,
				Message: "绑定参数错误",
			})
		return
	}

	moneyLogService := service.MoneyLogServiceInstance(repositories.MoneyLogRepositoryInstance(helper.GetDBByName(constant.DBMerchant)))
	//moneyLogData.Status = 2
	err = moneyLogService.Update(moneyLogData)
	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: "ok",
			Content: "",
		})
}
