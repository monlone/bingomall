package control

import (
	"bingomall/constant"
	"bingomall/helpers"
	"bingomall/models"
	"bingomall/repositories"
	"bingomall/services"
	"bingomall/system"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

// 用户登陆接口
// @Summary 用户登陆接口
// @Tags LoginController
// @Accept json
// @Produce json
// @Param username query string true "用户名"
// @Param password query string true "密码"
// @Success 200 {object} helpers.JsonObject
// @Router /api/login [post]
func Login(context *gin.Context) {
	//defer func() {
	//	if r := recover(); r != nil {
	//		fmt.Println(r)
	//	}
	//}()
	//err := errors.New("test")
	//if err.Error() != "" {
	//	panic(err)
	//}

	params := &helper.LoginParams{}
	if err := context.Bind(params); err == nil {
		userService := service.UserServiceInstance(repositories.UserRepositoryInstance(helper.GetUserDB()))
		user := userService.GetByUsername(params.Username)
		if user != nil && user.Password == helper.SHA256(params.Password) {
			user.LogonCount += 1
			user.LoginTime = time.Now()
			err := userService.SaveOrUpdate(user)
			if err == nil {
				generateToken(context, user)
			} else {
				context.JSON(http.StatusOK, helper.JsonObject{
					Code:    constant.SaveUserError,
					Message: helper.StatusText(helper.LoginStatusSQLErr),
					Content: err,
				})
				context.Abort() //退出
			}
		} else {
			context.JSON(http.StatusOK, helper.JsonObject{
				Code:    constant.UserNotExistOrPWDError,
				Message: helper.StatusText(helper.LoginStatusErr),
			})
		}
	} else {
		context.JSON(http.StatusUnprocessableEntity, helper.JsonObject{
			Code:    1004,
			Message: helper.StatusText(helper.BindModelErr),
			Content: err,
		})
	}
}

// 生成令牌
func generateToken(context *gin.Context, user *model.User) {
	j := system.NewJWT()
	claims := system.CustomClaims{
		ID:    user.ID,
		Name:  user.Username,
		Phone: user.Phone,
		StandardClaims: jwt.StandardClaims{
			NotBefore: int64(time.Now().Unix() + system.GetTokenConfig().ActiveTime),       // 签名生效时间
			ExpiresAt: int64(time.Now().Unix() + system.GetTokenConfig().ExpiredTime*3600), // 过期时间
			Issuer:    system.GetTokenConfig().Issuer,
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		context.JSON(http.StatusOK, helper.JsonObject{
			Code:    1005,
			Message: err.Error(),
		})
		context.Abort()
	}
	context.JSON(http.StatusOK, helper.JsonObject{
		Code:    0,
		Message: helper.StatusText(helper.LoginStatusOK),
		Content: gin.H{"accessToken": token, "user": user},
	})
}

func generateAppToken(user *model.User) (token string, err error) { //生成token
	j := system.NewJWT()
	claims := system.CustomClaims{
		ID:     user.ID,
		Name:   user.Username,
		Openid: user.OpenId,

		StandardClaims: jwt.StandardClaims{
			NotBefore: int64(time.Now().Unix() + system.GetTokenConfig().ActiveTime),       // 签名生效时间
			ExpiresAt: int64(time.Now().Unix() + system.GetTokenConfig().ExpiredTime*3600), // 过期时间
			Issuer:    system.GetTokenConfig().Issuer,
		},
	}
	token, err = j.CreateToken(claims)

	return
}

func init() {
	// 先读取Token配置文件
	err := system.LoadTokenConfig("./conf/token-config.yml")
	if err != nil {
		helper.ErrorLogger.Errorln("读取Token配置错误：", err)
		helper.WorkLogger.Errorln("读取Token配置错误：", err)
	}
	if len(strings.TrimSpace(system.GetTokenConfig().SignKey)) > 0 {
		system.SetSignKey(system.GetTokenConfig().SignKey)
	}
}
