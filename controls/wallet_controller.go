package control

import (
	"bingomall/constant"
	helper "bingomall/helpers"
	"bingomall/models"
	"bingomall/repositories"
	service "bingomall/services"
	"github.com/gin-gonic/gin"

	"net/http"
)

// 用户钱包
// @Summary 获取用户积分
// @Accept json
// @Produce json
// @Tags ScoreController
// @Param Gin-Access-Token header string true "令牌"
// @Success 200 {object} helpers.JsonObject
// @Router /api/user/score [get]
func Score(context *gin.Context) {
	userId := helper.GetUserID(context)
	db := helper.GetDBByName(constant.DBMerchant)
	walletService := service.WalletServiceInstance(repositories.WalletRepositoryInstance(db))
	score := walletService.GetScoreByUserID(userId)
	data := make(map[string]uint64)
	data["score"] = score
	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: helper.StatusText(helper.GetDataOK),
			Content: data,
		})
}

// @Summary 获取用户积分
// @Accept json
// @Produce json
// @Tags ScoreController
// @Param Gin-Access-Token header string true "令牌"
// @Success 200 {object} helpers.JsonObject
// @Router /api/user/score [get]
func Wallet(context *gin.Context) {
	userId := helper.GetUserID(context)
	db := helper.GetDBByName(constant.DBMerchant)
	walletService := service.WalletServiceInstance(repositories.WalletRepositoryInstance(db))
	var wallet *model.Wallet
	wallet = walletService.GetByUserID(userId)

	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: helper.StatusText(helper.GetDataOK),
			Content: wallet,
		})
}
