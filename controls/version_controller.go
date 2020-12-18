package control

import (
	helper "bingomall/helpers"
	"github.com/gin-gonic/gin"

	"net/http"
)

// 版本校验
// @Summary 版本校验
// @Accept json
// @Produce json
// @Tags ScoreController
// @Success 200 {object} helpers.JsonObject
// @Router /api/version/check [get]
func Check(context *gin.Context) {
	data := make(map[string]interface{})
	versionIOS := make(map[string]interface{})
	versionAndroid := make(map[string]interface{})
	versionIOS["latest_version"] = "1.0.1"
	versionIOS["force_update"] = 0
	data["ios"] = versionIOS
	versionAndroid["latest_version"] = "1.0.1"
	versionAndroid["force_update"] = 0
	data["android"] = versionAndroid
	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: helper.StatusText(helper.GetDataOK),
			Content: data,
		})
}
