package control

import (
	"bingomall/constant"
	helper "bingomall/helpers"
	"bingomall/helpers/convention"
	model "bingomall/models"
	"bingomall/repositories"
	service "bingomall/services"
	"github.com/gin-gonic/gin"

	"net/http"
)

// 用户
// @Summary 用户邮寄地址
// @Accept json
// @Produce json
// @Tags ScoreController
// @Param Gin-Access-Token header string true "令牌"
// @Success 200 {object} helpers.JsonObject
// @Router /api/user/address_list [get]
func AddressList(context *gin.Context) {
	userId := helper.GetUserID(context)
	db := helper.GetDBByName(constant.DBMerchant)
	addressService := service.AddressServiceInstance(repositories.AddressRepositoryInstance(db))
	addressList := addressService.GetByUserID(userId)
	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: helper.StatusText(helper.GetDataOK),
			Content: addressList,
		})
}

// 获取 address 详情
// @Summary 获取 address 详情
// @Tags AddressController
// @Accept json
// @Produce json
// @Success 200 {object} model.Address
// @Router /api/app/address_detail [get]
func AddressDetail(context *gin.Context) {
	addressService := service.AddressServiceInstance(repositories.AddressRepositoryInstance(helper.GetDBByName(constant.DBMerchant)))
	addressId := context.Query("addressId")
	userId := helper.GetUserID(context)
	var address *model.Address
	address = addressService.GetAddressByAddressId(convention.StringToUint64(addressId), userId)

	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: "ok",
			Content: address,
		})
}

// 获取默认 address 详情
// @Summary 获取默认 address 详情
// @Tags AddressController
// @Accept json
// @Produce json
// @Success 200 {object} model.Address
// @Router /api/app/address_default [get]
func AddressDefault(context *gin.Context) {
	addressService := service.AddressServiceInstance(repositories.AddressRepositoryInstance(helper.GetDBByName(constant.DBMerchant)))
	userId := helper.GetUserID(context)
	address := addressService.GetDefaultAddressByUserID(userId)

	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: "ok",
			Content: address,
		})
}

// 添加收货地址
// @Summary 添加收货地址
// @Tags AddressController
// @Accept json
// @Produce json
// @Success 200 {object} model.Address
// @Router /api/address/save [post]
func SaveAddress(context *gin.Context) {
	address := &model.Address{}
	if err := context.Bind(address); err == nil {
		addressService := service.AddressServiceInstance(repositories.AddressRepositoryInstance(helper.GetDBByName(constant.DBMerchant)))
		userId := helper.GetUserID(context)
		address.UserID = userId
		err := addressService.SaveOrUpdate(address)
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
					Code:    4204,
					Message: err.Error(),
				})
			return
		}
	} else {
		context.JSON(http.StatusUnprocessableEntity, helper.JsonObject{
			Code:    4203,
			Message: helper.StatusText(helper.BindModelErr),
			Content: err,
		})
	}
}

// 把收货地址置为默认的，前端这个写得不好，先保留了
// @Summary 添加收货地址
// @Tags AddressController
// @Accept json
// @Produce json
// @Success 200 {object} model.Address
// @Router /api/set_default_address [post]
func SetDefaultAddress(context *gin.Context) {
	addressService := service.AddressServiceInstance(repositories.AddressRepositoryInstance(helper.GetDBByName(constant.DBMerchant)))
	userId := helper.GetUserID(context)
	address := &model.Address{}
	if err := context.Bind(address); err != nil {
		context.JSON(http.StatusUnprocessableEntity, helper.JsonObject{
			Code:    4204,
			Message: helper.StatusText(helper.BindModelErr),
			Content: err,
		})
	}
	_ = addressService.SetDefaultAddress(userId, address.ID)

	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: helper.StatusText(helper.SaveStatusOK),
		})
}
