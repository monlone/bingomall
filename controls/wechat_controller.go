package control

import (
	"encoding/json"
	"fmt"
	"bingomall/constant"
	"bingomall/helpers"
	"bingomall/helpers/wechat"
	model "bingomall/models"
	"bingomall/repositories"
	service "bingomall/services"
	"log"
	"net/http"
	"time"

	"github.com/chanxuehong/rand"
	"github.com/chanxuehong/session"
	"github.com/chanxuehong/sid"
	mpOauth2 "github.com/chanxuehong/wechat/mp/oauth2"
	"github.com/chanxuehong/wechat/oauth2"
	"github.com/gin-gonic/gin"
)

var (
	sessionStorage                 = session.New(20*60, 60*60)
	oauth2Endpoint oauth2.Endpoint = mpOauth2.NewEndpoint(constant.WxAppId, constant.WxAppSecret)
)

type WXLoginResp struct {
	OpenId     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionId    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
	Code       string `json:"code"`
}

type RegisterForm struct {
	WechatCode    string `form:"code" binding:"required" json:"code"`
	EncryptedData string `form:"encryptedData" binding:"required" json:"encryptedData"`
	Iv            string `form:"iv" binding:"required" json:"iv"`
	Referrer      string `form:"referrer" json:"referrer"`
}

type LoginForm struct {
	WechatCode string `form:"code" binding:"required" json:"code"`
}

// wechat授权,建立必要的 session, 然后跳转到授权页面
// @Summary wechat授权
// @Tags UserController
// @Accept json
// @Produce json
// @Param username       query string true  "用户名"
// @Param password       query string true  "密码"
// @Param phone          query string true  "电话号码"
// @Param email          query string true  "邮件"
// @Success 200 {object} helpers.JsonObject
// @Router /api/wechat/auth [post]
func Auth(context *gin.Context) {
	sId := sid.New()
	state := string(rand.NewHex())

	if err := sessionStorage.Add(sId, state); err != nil {
		context.JSON(http.StatusOK,
			&helper.JsonObject{
				Code:    2010,
				Message: err.Error(),
			})
		log.Println(err)
		return
	}

	context.SetCookie("sid", sId, 3600*24*30, "/", constant.Domain, false, true)
	AuthCodeURL := mpOauth2.AuthCodeURL(constant.WxAppId, constant.Oauth2RedirectURI, constant.Oauth2Scope, state)
	log.Println("AuthCodeURL:", AuthCodeURL)

	context.Redirect(http.StatusFound, AuthCodeURL)
}

// 授权后回调页面
func ExchangeToken(context *gin.Context) {
	log.Println(context)

	sId, err := context.Cookie("sid")

	if err != nil {
		context.JSON(http.StatusOK,
			&helper.JsonObject{
				Code:    2003,
				Message: err.Error(),
			})
		return
	}

	CallbackSession, err := sessionStorage.Get(sId)

	if err != nil {
		context.JSON(http.StatusOK,
			&helper.JsonObject{
				Code:    2004,
				Message: err.Error(),
			})
		return
	}

	savedState := CallbackSession.(string) // 一般是要序列化的, 这里保存在内存所以可以这么做

	code := context.Query("code")
	state := context.Query("state")
	if code == "" {
		log.Println("用户禁止授权")
		context.JSON(http.StatusOK,
			&helper.JsonObject{
				Code:    2005,
				Message: "用户禁止授权",
			})
		return
	}

	if state == "" {
		log.Println("state 参数为空")
		context.JSON(http.StatusOK,
			&helper.JsonObject{
				Code:    2006,
				Message: "state 参数为空",
			})
		return
	}
	if savedState != state {
		str := fmt.Sprintf("state 不匹配, session 中的为 %q, url 传递过来的是 %q", savedState, helper.Json(context.Request))
		context.JSON(http.StatusOK,
			&helper.JsonObject{
				Code:    2007,
				Message: str,
			})
		log.Println(str)
		return
	}

	oauth2Client := oauth2.Client{
		Endpoint: oauth2Endpoint,
	}
	token, err := oauth2Client.ExchangeToken(code)
	if err != nil {
		context.JSON(http.StatusOK,
			&helper.JsonObject{
				Code:    2008,
				Message: err.Error(),
			})
		log.Println(err)
		return
	}
	log.Printf("token: %+v\r\n", token)
	fmt.Println("token:", token)

	userInfo, err := mpOauth2.GetUserInfo(token.AccessToken, token.OpenId, "", nil)
	if err != nil {
		context.JSON(http.StatusOK,
			&helper.JsonObject{
				Code:    2009,
				Message: err.Error(),
			})
		log.Println(err)
		return
	}
	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: "ok",
			Content: userInfo,
		})

	log.Printf("userinfo: %+v\r\n", userInfo)
	return
}

// 这个函数以 code 作为输入, 返回调用微信接口得到的对象指针和异常情况
func RegisterByWeChat(context *gin.Context) {
	url := constant.WechatUserInfoUrl
	form := &RegisterForm{}
	//code := context.PostForm("code")  // 'Content-Type': 'application/x-www-form-urlencoded' 用这个方法取，
	// 'Content-Type': 'application/json'用context.BindJSON(&form)的方法取
	if context.BindJSON(&form) != nil {
		context.JSON(http.StatusOK,
			&helper.JsonObject{
				Code:    500,
				Message: "请求参数错误",
			})
		return
	}

	url = fmt.Sprintf(url, constant.WechatAppId, constant.WechatAppSecret, form.WechatCode)

	resp, err := http.Get(url)
	if err != nil {
		context.JSON(http.StatusOK,
			&helper.JsonObject{
				Code:    500,
				Message: "请求微信服务器错误",
			})
		return
	}
	defer resp.Body.Close()

	// 解析http请求中body 数据到我们定义的结构体中
	wxResp := WXLoginResp{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&wxResp); err != nil {
		return
	}

	pc := wechat.WxBizDataCrypt{AppID: constant.WechatAppId, SessionKey: wxResp.SessionKey}
	userInfo, err := pc.Decrypt(form.EncryptedData, form.Iv, true) //第三个参数解释：需要返回JSON数据类型时使用true, 需要返回map数据类型时使用false
	if err != nil {
		fmt.Println(err)
	}

	person := helper.BaseUserInfo{}
	_ = json.Unmarshal([]byte(userInfo.(string)), &person) //把json中多个数据取部分数据出来

	// 判断微信接口返回的是否是一个异常情况
	if wxResp.ErrCode != 0 {
		context.JSON(http.StatusOK,
			&helper.JsonObject{
				Code:    wxResp.ErrCode,
				Message: wxResp.ErrMsg,
			})
		return
	}

	userService := service.UserServiceInstance(repositories.UserRepositoryInstance(helper.GetUserDB()))
	user := userService.GetByUserUnionId(wxResp.UnionId)
	if user != nil {
		user.LogonCount += 1
		user.LoginTime = time.Now()
		user.OpenId = wxResp.OpenId
		user.Nickname = person.Nickname
		user.Avatar = person.AvatarUrl
		err := userService.UpdateByApp(user)
		if err != nil {
			context.JSON(http.StatusOK, helper.JsonObject{
				Code:    4003,
				Message: helper.StatusText(helper.LoginStatusSQLErr),
				Content: err,
			})
			return
		}
	} else {
		user = &model.User{}
		user.OpenId = wxResp.OpenId
		user.UnionId = wxResp.UnionId
		user.AccessToken = wxResp.SessionKey
		user.Nickname = person.Nickname
		user.Avatar = person.AvatarUrl
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
		db := helper.GetDBByName(constant.DBMerchant)
		walletService := service.WalletServiceInstance(repositories.WalletRepositoryInstance(db))
		err = walletService.InitMyWallet(user.ID)
		if err != nil {
			context.JSON(http.StatusOK,
				&helper.JsonObject{
					Code:    4004,
					Message: err.Error(),
				})
			return
		}
	}

	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: "ok",
			Content: wxResp,
		})
}

func LoginByWeChat(context *gin.Context) {
	form := &LoginForm{}
	if context.BindJSON(&form) != nil {
		context.JSON(http.StatusOK,
			&helper.JsonObject{
				Code:    500,
				Message: "请求参数错误",
			})
		return
	}

	//code := context.PostForm("code")  // 'Content-Type': 'application/x-www-form-urlencoded' 用这个方法取，
	// 'Content-Type': 'application/json'用context.BindJSON(&form)的方法取

	url := fmt.Sprintf(constant.WechatUserInfoUrl, constant.WechatAppId, constant.WechatAppSecret, form.WechatCode)

	resp, err := http.Get(url)
	if err != nil {
		context.JSON(http.StatusOK,
			&helper.JsonObject{
				Code:    500,
				Message: "请求微信服务器错误",
			})
		return
	}
	defer resp.Body.Close()

	// 解析http请求中body 数据到我们定义的结构体中
	wxResp := WXLoginResp{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&wxResp); err != nil {
		context.JSON(http.StatusOK,
			&helper.JsonObject{
				Code:    wxResp.ErrCode,
				Message: err.Error(),
			})
		return
	}

	if wxResp.ErrCode != 0 {
		context.JSON(http.StatusOK,
			&helper.JsonObject{
				Code:    wxResp.ErrCode,
				Message: wxResp.ErrMsg,
			})
		return
	}

	userService := service.UserServiceInstance(repositories.UserRepositoryInstance(helper.GetUserDB()))
	var user *model.User
	if wxResp.UnionId != "" {
		user = userService.GetByUserUnionId(wxResp.UnionId)
	} else {
		user = userService.GetByUserOpenId(wxResp.OpenId) //没有绑定公众号，没有unionid，这个临时这样子解决
	}

	person := helper.BaseUserInfo{}
	token, err := generateAppToken(user)
	person.Token = token
	person.UserID = user.ID

	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: "ok",
			Content: person,
		})
}
func CheckTokenWeChat(context *gin.Context) {
	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: "ok",
			Content: "",
		})
}
