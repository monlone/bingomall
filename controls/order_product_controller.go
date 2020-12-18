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

// 产品分类
// @Summary 用户积分
// @Accept json
// @Produce json
// @Tags ScoreController
// @Param Gin-Access-Token header string true "令牌"
// @Success 200 {object} helpers.JsonObject
// @Router /api/user/OrderProduct_list [get]
func OrderProductList(context *gin.Context) {
	shopId := context.DefaultQuery("shopId", "")
	db := helper.GetDBByName(constant.DBMerchant)
	orderProductService := service.OrderProductServiceInstance(repositories.OrderProductRepositoryInstance(db))
	list := orderProductService.GetAll(convention.StringToUint64(shopId))
	data := make(map[string]interface{})
	//fmt.Println(reflect.TypeOf(data).Kind())
	data["list"] = list
	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: helper.StatusText(helper.GetDataOK),
			Content: data,
		})
}
