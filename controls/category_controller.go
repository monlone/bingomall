package control

import (
	"bingomall/constant"
	helper "bingomall/helpers"
	"bingomall/helpers/convention"
	"bingomall/models"
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
// @Router /api/shop/category/list [get]
func CategoryList(context *gin.Context) {
	shopId := context.DefaultQuery("shopId", "")
	db := helper.GetDBByName(constant.DBMerchant)
	categoryService := service.CategoryServiceInstance(repositories.CategoryRepositoryInstance(db))
	var categoryList []*model.Category
	categoryList = categoryService.GetAll(convention.StringToUint64(shopId))

	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: helper.StatusText(helper.GetDataOK),
			Content: categoryList,
		})
}
