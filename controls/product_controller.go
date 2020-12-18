package control

import (
	"bingomall/constant"
	helper "bingomall/helpers"
	"bingomall/helpers/convention"
	model "bingomall/models"
	"bingomall/repositories"
	service "bingomall/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

// 获取 product 列表
// @Summary 获取 product 列表
// @Tags ProductController
// @Accept json
// @Produce json
// @Success 200 {object} model.ListsResponse
// @Router /api/shop/product_list [get]
func ProductList(context *gin.Context) {
	productService := service.ProductServiceInstance(repositories.ProductRepositoryInstance(helper.GetDBByName(constant.DBMerchant)))
	pageStr := context.DefaultQuery("page", "0")
	page, _ := strconv.Atoi(pageStr)
	pageSizeStr := context.DefaultQuery("page_size", constant.PageSize)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	shopId := context.Query("shop_id")
	productList := productService.GetListByShopId(page, pageSize, shopId)

	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: "ok",
			Content: productList,
		})
}

func BargainProduct(context *gin.Context) {
	productBargain := &model.ProductDetail{}

	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: "ok",
			Content: productBargain,
		})
}

func ProductListAll(context *gin.Context) {
	var product *model.Product
	productService := service.ProductServiceInstance(repositories.ProductRepositoryInstance(helper.GetDBByName(constant.DBMerchant)))
	pageStr := context.DefaultQuery("page", "0")
	page, _ := strconv.Atoi(pageStr)
	pageSizeStr := context.DefaultQuery("page_size", constant.PageSize)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	productList := productService.GetPage(page, pageSize, product)

	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: "ok",
			Content: productList,
		})
}

// 获取 product 详情
// @Summary 获取 product 详情
// @Tags ProductController
// @Accept json
// @Produce json
// @Success 200 {object} model.Product
// @Router /api/shop/product_detail [get]
func ProductDetail(context *gin.Context) {
	db := helper.GetDBByName(constant.DBMerchant)
	productService := service.ProductServiceInstance(repositories.ProductRepositoryInstance(db))
	productId := convention.StringToUint64(context.Query("productId"))
	var productAllDetail = &model.ProductAllDetail{}
	product := productService.GetByProductId(productId)
	productOptionService := service.ProductOptionServiceInstance(repositories.ProductOptionRepositoryInstance(db))
	optionList := productOptionService.GetAllGroupByType(productId)
	skuService := service.SkuServiceInstance(repositories.SkuRepositoryInstance(db))
	skuList := skuService.GetAll(productId)

	productAllDetail.Product = product
	productAllDetail.Sku = skuList
	productAllDetail.Option = optionList
	var imageList []string
	imageList = append(imageList, "https://ss1.bdstatic.com/70cFvXSh_Q1YnxGkpoWK1HF6hhy/it/u=1821767832,4151893847&fm=11&gp=0.jpg")
	imageList = append(imageList, "https://ss1.bdstatic.com/70cFvXSh_Q1YnxGkpoWK1HF6hhy/it/u=241281405,3821910790&fm=11&gp=0.jpg")
	productAllDetail.Image = imageList
	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: "ok",
			Content: productAllDetail,
		})
}

func ProductPrice(context *gin.Context) {
	db := helper.GetDBByName(constant.DBMerchant)
	//productOptionIds := context.Query("productOptionIds")
	productDetail := &model.ProductDetail{}

	if err := context.Bind(productDetail); err != nil {

	}
	productOptionIdsMap := productDetail.ProductOptionIds
	var productOptionIds []string
	for _, v := range productOptionIdsMap {
		productOptionIds = append(productOptionIds, v)
	}
	sort.Strings(productOptionIds)
	combineId := strings.Join(productOptionIds, "_")

	skuService := service.SkuServiceInstance(repositories.SkuRepositoryInstance(db))
	sku := skuService.GetByProductIdCombineId(productDetail.ProductId, combineId)
	productService := service.ProductServiceInstance(repositories.ProductRepositoryInstance(db))
	var product *model.Product
	product = productService.GetByProductId(productDetail.ProductId)
	var productPriceObj = &model.ProductPriceObj{}
	productPriceObj.Sku = sku
	productPriceObj.Product = product
	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: "ok",
			Content: productPriceObj,
		})
}

func GetProductListByOrderProductList(db *gorm.DB, orderProductList []*model.OrderProduct) []*model.Product {
	productIds := make([]uint64, 1)

	for _, v := range orderProductList {
		productIds = append(productIds, v.ProductId)
	}

	var productList []*model.Product
	productService := service.ProductServiceInstance(repositories.ProductRepositoryInstance(db))
	productList = productService.GetListByProductIds(productIds)

	return productList
}

func GetProductListByOrderId(db *gorm.DB, orderId uint64) []*model.Product {
	orderProductService := service.OrderProductServiceInstance(repositories.OrderProductRepositoryInstance(db))
	orderProductList := orderProductService.GetOrderItemsByOrderId(orderId)
	productIds := make([]uint64, 1)
	for _, v := range orderProductList {
		productIds = append(productIds, v.ProductId)
	}

	for _, v := range orderProductList {
		productIds = append(productIds, v.ProductId)
	}

	var productList []*model.Product
	productService := service.ProductServiceInstance(repositories.ProductRepositoryInstance(db))
	productList = productService.GetListByProductIds(productIds)

	return productList
}

// 编辑 product
// @Summary  product 编辑
// @Tags ProductController
// @Accept json
// @Produce json
// @Success 200 {object} model.Product
// @Router /api/product/save [post]
func SaveProduct(context *gin.Context) {
	product := &model.Product{}
	if err := context.Bind(product); err == nil {
		productService := service.ProductServiceInstance(repositories.ProductRepositoryInstance(helper.GetDBByName(constant.DBMerchant)))
		if int(product.Discount+product.DiscountLevel+product.DiscountUser+product.DiscountPlatform) != 100 {
			context.JSON(http.StatusOK,
				&helper.JsonObject{
					Code:    4209,
					Message: "折扣填写错误！折扣率加起来必须为100%!",
				})
			return
		}
		err := productService.SaveOrUpdate(product)
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
			Content: err.Error(),
		})
	}
}
