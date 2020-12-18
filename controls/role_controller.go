package control

import (
	helper "bingomall/helpers"
	"bingomall/models"
	"bingomall/repositories"
	service "bingomall/services"
	"github.com/gin-gonic/gin"

	"net/http"
)

// 添加修改角色
// @Summary 添加修改角色
// @Tags RoleController
// @Accept json
// @Produce json
// @Param id             query string false "角色id,新增时id为空"
// @Param role_name      query string true  "角色名称"
// @Param role_key       query string true  "角色类别标识"
// @Param description    query string true  "角色描述信息"
// @Success 200 {object} helpers.JsonObject
// @Router /api/save_role [post]
func SaveRole(context *gin.Context) {
	var role model.Role
	if err := context.Bind(&role); err == nil {
		roleService := service.RoleServiceInstance(repositories.RoleRepositoryInstance(helper.GetUserDB()))
		err := roleService.SaveOrUpdate(&role)
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
					Code:    3001,
					Message: helper.StatusText(helper.SaveStatusErr),
					Content: err.Error(),
				})
			return
		}
	} else {
		context.JSON(http.StatusUnprocessableEntity, helper.JsonObject{
			Code:    3002,
			Message: helper.StatusText(helper.BindModelErr),
			Content: err,
		})
		return
	}
}
