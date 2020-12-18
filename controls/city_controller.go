package control

import (
	"bingomall/constant"
	helper "bingomall/helpers"
	"bingomall/helpers/convention"
	"bingomall/repositories"
	service "bingomall/services"
	"github.com/gin-gonic/gin"

	"net/http"
)

// 省份列表
// @Summary 省份列表
// @Accept json
// @Produce json
// @Tags ScoreController
// @Param Gin-Access-Token header string true "令牌"
// @Success 200 {object} helpers.JsonObject
// @Router /api/common/city_list [get]
func CityList(context *gin.Context) {
	provinceCode := context.DefaultQuery("provinceCode", "")
	db := helper.GetDBByName(constant.DBMerchant)
	cityService := service.CityServiceInstance(repositories.CityRepositoryInstance(db))
	list := cityService.GetAll(convention.StringToUint64(provinceCode))

	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: helper.StatusText(helper.GetDataOK),
			Content: list,
		})
}
