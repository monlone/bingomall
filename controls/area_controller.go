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
// @Router /api/common/area_list [get]
func AreaList(context *gin.Context) {
	cityCode := context.DefaultQuery("cityCode", "")
	db := helper.GetDBByName(constant.DBMerchant)
	areaService := service.AreaServiceInstance(repositories.AreaRepositoryInstance(db))
	list := areaService.GetAll(convention.StringToUint64(cityCode))

	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: helper.StatusText(helper.GetDataOK),
			Content: list,
		})
}
