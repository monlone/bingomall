package helper

import (
	"bingomall/system"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetClaims(context *gin.Context) (data *system.CustomClaims) {
	info, _ := context.Get("claims")
	data = info.(*system.CustomClaims)

	if data == nil {
		context.JSON(http.StatusOK, JsonObject{
			Code:    42001,
			Content: "获取claims错误",
		})
		context.Abort()
	}
	return
}

func GetTokenInfo(context *gin.Context) (info *system.CustomClaims) {
	return GetClaims(context)
}

func GetUserID(context *gin.Context) (userId uint64) {
	tokenInfo := GetClaims(context)
	if nil == tokenInfo {
		return
	}

	return tokenInfo.ID
}

func GetOpenId(context *gin.Context) (openId string) {
	tokenInfo := GetClaims(context)
	if nil == tokenInfo {
		return
	}

	return tokenInfo.Openid
}
