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
// @Router /api/common/province_list [get]
func ProvinceList(context *gin.Context) {
	provinceId := context.DefaultQuery("provinceId", "")
	db := helper.GetDBByName(constant.DBMerchant)
	provinceService := service.ProvinceServiceInstance(repositories.ProvinceRepositoryInstance(db))
	list := provinceService.GetAll(convention.StringToUint64(provinceId))

	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: helper.StatusText(helper.GetDataOK),
			Content: list,
		})
}
