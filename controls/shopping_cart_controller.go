package control

import (
	"bingomall/constant"
	helper "bingomall/helpers"
	"bingomall/helpers/convention"
	model "bingomall/models"
	"bingomall/repositories"
	service "bingomall/services"
	"github.com/gin-gonic/gin"
	"sort"
	"strings"

	"net/http"
)

// 购物车
// @Summary 添加商品到购物车
// @Accept json
// @Produce json
// @Tags ScoreController
// @Param Gin-Access-Token header string true "令牌"
// @Success 200 {object} helpers.JsonObject
// @Router /api/user/AddToCart [get]
func AddToCart(context *gin.Context) {
	db := helper.GetDBByName(constant.DBMerchant)
	skuService := service.SkuServiceInstance(repositories.SkuRepositoryInstance(db))
	shoppingCartPost := &model.ShoppingCartPost{}
	shoppingCart := &model.ShoppingCart{}
	userId := helper.GetUserID(context)

	if err := context.Bind(shoppingCartPost); err != nil {

	}
	productOptionIdsMap := shoppingCartPost.ProductOptionIds
	var productOptionIds []string
	for _, v := range productOptionIdsMap {
		productOptionIds = append(productOptionIds, v)
	}
	sort.Strings(productOptionIds)
	combineId := strings.Join(productOptionIds, "_")
	sku := skuService.GetByProductIdCombineId(convention.StringToUint64(shoppingCartPost.ProductId), combineId)
	if combineId == "" {
		context.JSON(http.StatusOK,
			&helper.JsonObject{
				Code:    0,
				Message: helper.StatusText(helper.UpdateObjIsNil),
				Content: nil,
			})
		return
	}

	shoppingCartService := service.ShoppingCartServiceInstance(repositories.ShoppingCartRepositoryInstance(db))
	shoppingCart.ProductId = sku.ProductId
	shoppingCart.Number = shoppingCartPost.Number
	shoppingCart.Price = sku.Price
	shoppingCart.UserID = userId
	//shoppingCart.CombineId = combineId
	shoppingCart.SkuId = sku.ID

	list := shoppingCartService.SaveOrUpdate(shoppingCart)

	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: helper.StatusText(helper.SaveStatusOK),
			Content: list,
		})
}

// @Summary 获取购物车信息
// @Accept json
// @Produce json
// @Tags ScoreController
// @Param Gin-Access-Token header string true "令牌"
// @Success 200 {object} helpers.JsonObject
// @Router /api/user/ShoppingCartInfo [get]
func ShoppingCartInfo(context *gin.Context) {
	db := helper.GetDBByName(constant.DBMerchant)
	shoppingCartService := service.ShoppingCartServiceInstance(repositories.ShoppingCartRepositoryInstance(db))
	var shoppingCartList []*model.ShoppingCart
	var shoppingCartInfo []*model.ShoppingCartInfo
	userId := helper.GetUserID(context)

	shoppingCartList = shoppingCartService.GetShoppingCartByUserID(userId)
	if shoppingCartList != nil {
		for _, shoppingCart := range shoppingCartList {
			optionList := make(map[uint8][]*model.ProductOption)
			for _, v := range shoppingCart.OptionList {
				if _, ok := optionList[v.Type]; ok {
					temp := optionList[v.Type]
					temp = append(temp, v)
					optionList[v.Type] = temp
				} else {
					var temp []*model.ProductOption
					temp = append(temp, v)
					optionList[v.Type] = temp
				}
			}
			cartInfoTemp := &model.ShoppingCartInfo{}
			cartInfoTemp.BaseShoppingCart = shoppingCart.BaseShoppingCart
			cartInfoTemp.OptionList = optionList
			shoppingCartInfo = append(shoppingCartInfo, cartInfoTemp)
		}
	}

	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: helper.StatusText(helper.GetDataOK),
			Content: shoppingCartInfo,
		})

}
