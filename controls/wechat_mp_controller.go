package control

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"bingomall/constant"
	"bingomall/helpers"
	"bingomall/helpers/cache"
	"bingomall/helpers/convention"
	model "bingomall/models"
	"bingomall/repositories"
	service "bingomall/services"
	"github.com/chanxuehong/rand"
	"github.com/chanxuehong/sid"
	mpOauth2 "github.com/chanxuehong/wechat/mp/oauth2"
	"github.com/chanxuehong/wechat/util"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
)

//var (
//	mpOauth2Endpoint wxOauth2.Endpoint = mpOauth2.NewEndpoint(constant.MPWechatAppId, constant.MPWxAppSecret)
//)

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
// @Router /api/wechat/wechat_mp_register [post]
func WechatMPRegister(context *gin.Context) {
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
	levelUserID := context.Query("level_user_id")
	context.SetCookie("level_user_id", levelUserID, 3600*24, "/", constant.Domain, false, true)
	AuthCodeURL := mpOauth2.AuthCodeURL(constant.MPWechatAppId, constant.MPWechatCallback, constant.MPOauth2Scope, state)

	context.Redirect(http.StatusFound, AuthCodeURL)
}

func MerchantWechatMPRegisterToken(context *gin.Context) {
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

	client := cache.DefaultClient()
	token := string(rand.NewHex())
	err := client.Set(token, 1, 0).Err()
	if err != nil {
		panic(err)
	}
	data := map[string]string{}
	data["token"] = token
	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: "OK",
			Content: data,
		})
}

// wechat授权,建立必要的 session, 然后跳转到授权页面
// @Summary 业务员分享给商户注册
// @Tags UserController
// @Accept json
// @Produce json
// @Param user_id       query string true  "用户id"
// @Success 200 {object} helpers.JsonObject
// @Router /api/web/merchant/wechat_mp_register [post]
func MerchantWechatMPRegister(context *gin.Context) {
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

	phone := context.PostForm("phone")
	title := context.PostForm("title")
	description := context.PostForm("description")
	token := context.PostForm("token")

	if token == "" {
		context.JSON(http.StatusOK,
			&helper.JsonObject{
				Code:    2010,
				Message: "token错误",
			})
		return
	}

	client := cache.DefaultClient()
	err := client.Set(sId, state, 0).Err()
	if err != nil {
		fmt.Println(err)
	}
	val, err := client.Get(sId).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("token:", val)

	context.SetCookie("description", description, 3600*24, "/", constant.Domain, false, true)
	context.SetCookie("title", title, 3600*24, "/", constant.Domain, false, true)
	context.SetCookie("phone", phone, 3600*24, "/", constant.Domain, false, true)
	context.SetCookie("merchant_register", "1", 3600*24, "/", constant.Domain, false, true)
	AuthCodeURL := mpOauth2.AuthCodeURL(constant.MPWechatAppId, constant.MPWechatCallback, constant.MPOauth2Scope, state)

	data := map[string]string{}
	data["url"] = AuthCodeURL
	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: "ok",
			Content: data,
		})
}

func WechatMPConfigCallBack(context *gin.Context) {
	token := constant.MPWxAPPToken
	signature := context.Query("signature")
	timestamp := context.Query("timestamp")
	nonce := context.Query("nonce")
	echostr := context.Query("echostr")
	tmps := []string{token, timestamp, nonce}
	sort.Strings(tmps)
	tmpStr := tmps[0] + tmps[1] + tmps[2]
	tmp := str2sha1(tmpStr)
	if tmp == signature {
		n, err := fmt.Fprintf(context.Writer, echostr)
		if err != nil {
			fmt.Println(n, err)
		}
	}
}

func str2sha1(data string) string {
	t := sha1.New()
	n, err := io.WriteString(t, data)
	if err != nil {
		fmt.Println(n, err)
	}
	return fmt.Sprintf("%x", t.Sum(nil))
}

//微信公众号回调
// /api/web/mpcallback
func WechatMPCallBack(context *gin.Context) {
	log.Println(context)
	echostr := context.DefaultQuery("echostr", "")
	if echostr != "" {
		WechatMPConfigCallBack(context)
		return
	}

	sId, err := context.Cookie("sid")
	levelUserIDStr, _ := context.Cookie("level_user_id") //微信分享的userId
	levelUserID := convention.StringToUint64(levelUserIDStr)
	fmt.Println("levelUserID:", levelUserID)

	if err != nil {
		context.JSON(http.StatusOK,
			&helper.JsonObject{
				Code:    2003,
				Message: err.Error(),
			})
		context.Abort()
		return
	}

	CallbackSession, err := sessionStorage.Get(sId)
	var savedState string

	if CallbackSession == nil {
		client := cache.DefaultClient()
		savedState, err = client.Get(sId).Result()
		err = nil
	} else {
		savedState = CallbackSession.(string) // 一般是要序列化的, 这里保存在内存所以可以这么做
	}

	if err != nil {
		context.JSON(http.StatusOK,
			&helper.JsonObject{
				Code:    2004,
				Message: err.Error(),
				Content: savedState,
			})
		return
	}

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

	oauth2Client := mpOauth2.NewEndpoint(constant.MPWechatAppId, constant.MPWxAppSecret)
	//以下链接获取access_token
	accessTokenUrl := oauth2Client.ExchangeTokenURL(code)

	log.Printf("accessTokenUrl: %+v\r\n", accessTokenUrl)
	httpClient := util.DefaultHttpClient

	httpResp, err := httpClient.Get(accessTokenUrl)
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
	refreshTokenUrl := oauth2Client.RefreshTokenURL(data.RefreshToken)
	httpRefreshResp, err := httpClient.Get(refreshTokenUrl)
	if err != nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			_ = httpRefreshResp.Body.Close()
		}
	}()

	if httpRefreshResp.StatusCode != http.StatusOK {
		err = fmt.Errorf("http.Status: %s", httpRefreshResp.Status)
		return
	}

	strRefresh, _ := ioutil.ReadAll(httpRefreshResp.Body)
	var dataRefresh model.CallbackInfo
	if err := json.Unmarshal(strRefresh, &dataRefresh); err != nil {
		panic(err)
	}
	userService := service.UserServiceInstance(repositories.UserRepositoryInstance(helper.GetUserDB()))
	userInfo, _ := mpOauth2.GetUserInfo(dataRefresh.AccessToken, dataRefresh.OpenId, "", nil)
	user := userService.GetByUserUnionId(userInfo.UnionId)
	//h5获取用户信息

	fmt.Println("userInfo:", helper.Json(userInfo))
	if user == nil {
		phone, err := context.Cookie("phone")
		userTemp := &model.User{}
		userTemp.MPOpenId = userInfo.OpenId
		userTemp.UnionId = userInfo.UnionId
		userTemp.AccessToken = data.AccessToken
		userTemp.RefreshToken = data.RefreshToken
		userTemp.MultiLevel = levelUserID
		userTemp.Phone = phone
		userTemp.Type = constant.NormalUser
		userTemp.Nickname = data.Nickname
		merchantRegister, err := context.Cookie("merchant_register")
		if merchantRegister != "" {
			userTemp.Type = constant.MerchantUser
		}

		err = userService.SaveByApp(userTemp)

		//info := userService.GetByUserUnionId(userInfo.UnionId)
		//phone, err := context.Cookie("phone")
		//description, _ := context.Cookie("Description")
		//title, err := context.Cookie("title")
		//if phone != "" {
		//	merchant := &model.Merchant{}
		//	merchantService := service.MerchantServiceInstance(repositories.MerchantRepositoryInstance(helper.GetDBByName(constant.DBMerchant)))
		//	merchant.Phone = phone
		//	merchant.Title = title
		//	merchant.Description = description
		//	merchant.UserID = info.UserID
		//	merchant.MerchantId = uuid.NewV4().String()
		//	err := merchantService.SaveOrUpdate(merchant)
		//	fmt.Println(err)
		//}

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
		err = walletService.InitMyWallet(userTemp.ID)
		if err != nil {
			context.JSON(http.StatusOK,
				&helper.JsonObject{
					Code:    4004,
					Message: err.Error(),
				})
			return
		}
	} else {
		phone, err := context.Cookie("phone")
		user.MPOpenId = userInfo.OpenId
		user.UnionId = userInfo.UnionId
		user.AccessToken = data.AccessToken
		user.RefreshToken = data.RefreshToken
		user.Nickname = data.Nickname
		if user.MultiLevel == 0 {
			user.MultiLevel = levelUserID
		}
		if len(phone) > 0 {
			user.Phone = phone
		}
		user.Type = constant.NormalUser
		merchantRegister, err := context.Cookie("merchant_register")
		if merchantRegister != "" {
			user.Type = constant.MerchantUser
		}

		fmt.Println("user in WechatMPCallBack:", helper.Json(user))
		err = userService.UpdateByApp(user)
		fmt.Println(err)
	}

	context.Redirect(http.StatusMovedPermanently, constant.JumpUrl)
}

//func MerchantWechatMPCallBack(context *gin.Context) {
//	log.Println(context)
//	echostr := context.DefaultQuery("echostr", "")
//	if echostr != "" {
//		WechatMPConfigCallBack(context)
//		return
//	}
//
//	sId, err := context.Cookie("sid")
//	levelUserID, _ := context.Cookie("level_user_id") //微信分享的userId
//	fmt.Println("levelUserID:", levelUserID)
//
//	if err != nil {
//		context.JSON(http.StatusOK,
//			&helper.JsonObject{
//				Code:    "2003",
//				Message: err.Error(),
//			})
//		context.Abort()
//		return
//	}
//
//	CallbackSession, err := sessionStorage.Get(sId)
//
//	if err != nil {
//		context.JSON(http.StatusOK,
//			&helper.JsonObject{
//				Code:    "2004",
//				Message: err.Error(),
//			})
//		return
//	}
//
//	savedState := CallbackSession.(string) // 一般是要序列化的, 这里保存在内存所以可以这么做
//
//	code := context.Query("code")
//	state := context.Query("state")
//	if code == "" {
//		log.Println("用户禁止授权")
//		context.JSON(http.StatusOK,
//			&helper.JsonObject{
//				Code:    "2005",
//				Message: "用户禁止授权",
//			})
//		return
//	}
//
//	if state == "" {
//		log.Println("state 参数为空")
//		context.JSON(http.StatusOK,
//			&helper.JsonObject{
//				Code:    "2006",
//				Message: "state 参数为空",
//			})
//		return
//	}
//	if savedState != state {
//		str := fmt.Sprintf("state 不匹配, session 中的为 %q, url 传递过来的是 %q", savedState, helper.Json(context.Request))
//		context.JSON(http.StatusOK,
//			&helper.JsonObject{
//				Code:    "2007",
//				Message: str,
//			})
//		log.Println(str)
//		return
//	}
//
//	oauth2Client := mpOauth2.NewEndpoint(constant.MPWechatAppId, constant.MPWxAppSecret)
//	//以下链接获取access_token
//	accessTokenUrl := oauth2Client.ExchangeTokenURL(code)
//
//	log.Printf("accessTokenUrl: %+v\r\n", accessTokenUrl)
//	httpClient := util.DefaultHttpClient
//
//	httpResp, err := httpClient.Get(accessTokenUrl)
//	if err != nil {
//		return
//	}
//	defer func() {
//		if r := recover(); r != nil {
//			_ = httpResp.Body.Close()
//		}
//	}()
//
//	if httpResp.StatusCode != http.StatusOK {
//		err = fmt.Errorf("http.Status: %s", httpResp.Status)
//		return
//	}
//
//	str, _ := ioutil.ReadAll(httpResp.Body)
//	var data model.CallbackInfo
//	if err := json.Unmarshal(str, &data); err != nil {
//		fmt.Println("err in panic, err:", err)
//		panic(err)
//	}
//	refreshTokenUrl := oauth2Client.RefreshTokenURL(data.RefreshToken)
//	fmt.Println("refreshTokenUrl:", refreshTokenUrl)
//	httpRefreshResp, err := httpClient.Get(refreshTokenUrl)
//	if err != nil {
//		return
//	}
//	defer func() {
//		if r := recover(); r != nil {
//			_ = httpRefreshResp.Body.Close()
//		}
//	}()
//
//	if httpRefreshResp.StatusCode != http.StatusOK {
//		err = fmt.Errorf("http.Status: %s", httpRefreshResp.Status)
//		return
//	}
//
//	strRefresh, _ := ioutil.ReadAll(httpRefreshResp.Body)
//	var dataRefresh model.CallbackInfo
//	if err := json.Unmarshal(strRefresh, &dataRefresh); err != nil {
//		fmt.Println("err in panic, err:", err)
//		panic(err)
//	}
//	fmt.Println("dataRefresh:", helper.Json(dataRefresh))
//	userService := service.UserServiceInstance(repositories.UserRepositoryInstance(helper.GetUserDB()))
//	user := userService.GetByUserUnionId(dataRefresh.UnionId)
//	userInfo, _ := mpOauth2.GetUserInfo(dataRefresh.AccessToken, dataRefresh.OpenId, "", nil)
//
//	if user == nil {
//		user = &model.User{}
//		user.MPOpenId = userInfo.OpenId
//		user.UnionId = data.UnionId
//		user.AccessToken = data.AccessToken
//		user.RefreshToken = data.RefreshToken
//		user.MultiLevel = levelUserID
//		user.Type = constant.MerchantUser
//
//		err = userService.SaveByApp(user)
//		if err != nil {
//			context.JSON(http.StatusOK,
//				&helper.JsonObject{
//					Code:    "4100",
//					Message: err.Error(),
//				})
//			return
//		}
//	}
//
//	//TODO 跳转到商户添加页面
//	context.Redirect(200, "http://a.app.qq.com/o/simple.jsp?pkgname=immortal.lh2424.com.immortal&from=timeline&isappinstalled=0")
//}
