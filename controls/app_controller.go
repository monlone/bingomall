package control

import (
	"encoding/json"
	"fmt"
	"bingomall/constant"
	"bingomall/helpers"
	"bingomall/helpers/ecode"
	model "bingomall/models"
	"bingomall/repositories"
	"bingomall/services"
	mpOauth2 "github.com/chanxuehong/wechat/mp/oauth2"
	"github.com/chanxuehong/wechat/util"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// 授权后回调页面
// app授权
// @Summary app授权
// @Tags AppController
// @Accept json
// @Produce json
// @Param code       query string true  "token"
// @Success 200 {object} helpers.JsonObject
// @Router /api/app/exchange_token [post]

func AppExchangeToken(context *gin.Context) {
	code := context.PostForm("code")
	var token string
	httpClient := util.DefaultHttpClient

	_urlParams := "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
	_url := fmt.Sprintf(_urlParams, constant.WxAppId, constant.WxAppSecret, code)

	httpResp, err := httpClient.Get(_url)
	if err != nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			_ = httpResp.Body.Close()
		}
	}()

	if httpResp.StatusCode != http.StatusOK {
		err = fmt.Errorf("http.Status: %s", httpResp.Status)
		return
	}

	str, _ := ioutil.ReadAll(httpResp.Body)
	var data model.CallbackInfo
	if err := json.Unmarshal(str, &data); err != nil {
		fmt.Println("err in panic, err:", err)
		panic(err)
	}

	helper.ServiceLogger.Println("wechat mini program return:", helper.Json(data))

	if "" == data.UnionId {
		context.JSON(http.StatusOK, helper.JsonObject{
			Code:    40029,
			Message: helper.StatusText(helper.WechatCodeErr),
			Content: err,
		})
		return
	}

	userService := service.UserServiceInstance(repositories.UserRepositoryInstance(helper.GetUserDB()))
	user := userService.GetByUserUnionId(data.UnionId)
	var weChatInfo model.WechatUserInfo
	userInfo, err := mpOauth2.GetUserInfo(data.AccessToken, data.OpenId, "", nil)
	//weChatInfo.OpenId = data.OpenId
	if userInfo == nil {
		context.JSON(http.StatusOK,
			&helper.JsonObject{
				Code:    4101,
				Message: "获取用户信息错误",
			})
		return
	}

	if user != nil {
		user.LogonCount += 1
		user.LoginTime = time.Now()
		user.OpenId = data.OpenId
		user.Nickname = userInfo.Nickname
		user.AccessToken = data.AccessToken
		user.RefreshToken = data.RefreshToken
		user.Avatar = data.HeadImageURL
		err := userService.UpdateByApp(user)
		if err != nil {
			context.JSON(http.StatusOK, helper.JsonObject{
				Code:    4003,
				Message: helper.StatusText(helper.LoginStatusSQLErr),
				Content: err,
			})
			return
		}
		token, err = generateAppToken(user)
		if err != nil {
			context.JSON(http.StatusOK, helper.JsonObject{
				Code:    40029,
				Message: helper.StatusText(helper.WechatCodeErr),
				Content: err,
			})
			return
		}
	} else {
		user = &model.User{}
		user.OpenId = data.OpenId
		user.UnionId = data.UnionId
		user.AccessToken = data.AccessToken
		user.RefreshToken = data.RefreshToken
		user.Nickname = userInfo.Nickname
		user.Avatar = data.HeadImageURL
		user.LoginTime = time.Now()
		err = userService.SaveByApp(user)
		if err != nil {
			context.JSON(http.StatusOK,
				&helper.JsonObject{
					Code:    4100,
					Message: err.Error(),
				})
			return
		}
		user := userService.GetByUserUnionId(data.UnionId)
		token, err = generateAppToken(user)
	}

	weChatInfo.Nickname = userInfo.Nickname
	weChatInfo.HeadImageURL = userInfo.HeadImageURL
	weChatInfo.AccessToken = token
	weChatInfo.UserType = user.Type
	weChatInfo.UserID = strconv.FormatUint(user.ID, 10)

	e := ecode.Cause(err)

	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    e.Code(),
			Message: e.Message(),
			Content: weChatInfo,
		})

	return
}
