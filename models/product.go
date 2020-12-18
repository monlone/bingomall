package model

import (
	"bingomall/constant"
	helper "bingomall/helpers"
	"gorm.io/gorm"
)

type Product struct {
	Model

	/** product id */
	//ProductId uint64 `gorm:"type:bigint;not null;" json:"productId" form:"productId"`

	/** 店铺id */
	ShopId uint64 `gorm:"type:bigint;" binding:"required" json:"shopId" form:"shopId"`

	/** product title */
	Title string `gorm:"type:varchar(255);" form:"title" json:"title"`

	/** 因为可以设置不同的sku不同的价格，所以有产品最高价格，这个价格只做展示用，以sku表里的价格为准，数据库存的是分*/
	Price uint64 `gorm:"type:bigint;" form:"price" json:"price"`

	/** 因为可以设置不同的sku不同的价格，所以有产品的最低价格，这个价格只做展示用，以sku表里的价格为准，数据库存的是分*/
	MinPrice uint64 `gorm:"type:int;" form:"minPrice" json:"minPrice"`

	/** 产品团购价格，数据库存的是分*/
	GroupPrice uint64 `gorm:"type:int;" form:"groupPrice" json:"groupPrice"` //TODO 以后要改，放到sku表里去

	/** 产品原始价格,这个价格只做展示用，以sku表里的价格为准，数据库存的是分 */
	OriginalPrice uint64 `gorm:"type:int;" form:"originalPrice" json:"originalPrice"`

	/**商品类型 1-实物商品要核销, 2-虚拟商品要核销，3-虚拟商品不用核销，4-实物商品要邮寄 5-实物商品不要邮寄 */
	Type uint8 `gorm:"type:tinyint(3);" form:"type" json:"type"`

	/** 商品的价格类型 1-普通价格商品，2-秒杀商品，3-砍价商品，4-拼团商品*/
	PriceType uint8 `gorm:"type:tinyint(3);" form:"priceType" json:"priceType"`

	/** 商品折扣（平台和商家洽谈的折扣）*/
	Discount float64 `gorm:"type:varchar(6);" form:"discount" json:"discount"`

	/** 推广者的上级提成，从推广者身上扣钱（配置一级分成比例在const里面）*/
	DiscountLevel float64 `gorm:"type:int;" form:"discountLevel" json:"discountLevel"`

	/** 平台自己的收益，比如产品100元，DiscountPlatform为10，则平台可以得10元*/
	DiscountPlatform float64 `gorm:"type:bigint;" form:"discountPlatform" json:"discountPlatform"`

	/** 购买者返现*/
	DiscountUser float64 `gorm:"type:int;" form:"discountUser" json:"discountUser"`

	/** 商户代平台核销的分成 */
	DiscountMerchant float64 `gorm:"type:int;default:0" form:"discountMerchant" json:"discountMerchant"`

	/** 产品主图 */
	ImageUrl string `gorm:"type:varchar(255);" form:"imageUrl" binding:"required" json:"imageUrl"`

	/** 商品描述 */
	Description string `gorm:"type:varchar(255);" form:"description" json:"description"`

	/** 购买须知 */
	Describe string `gorm:"type:varchar(5000);" form:"describe" json:"describe"`

	/** 商品图文详情 */
	DetailDescribe string `gorm:"type:varchar(5000);" form:"detailDescribe" json:"detailDescribe"`

	///** 商品库存 */
	//Stock int32 `gorm:"type:int;" form:"stories" json:"stock"`

	SoldNumber uint32 `gorm:"type:bigint" form:"soldNumber" json:"soldNumber"`

	/** 运费模板 */
	LogisticsTemplateId string `gorm:"type:varchar(36)" json:"logisticsTemplateId"`

	CrudTime
}

type ListsResponse struct {
	helper.JsonObject
	Data []Product
}

type ProductDetail struct {
	ProductId        uint64            `gorm:"type:bigint(20);column:product_id;not null;" json:"productId"`
	Title            string            `gorm:"type:varchar(255);" json:"title"`
	OriginalPrice    string            `gorm:"type:varchar(255);" json:"price"`
	ImageUrl         string            `gorm:"type:varchar(500);" json:"imageUrl"`
	ProductOptionIds map[string]string `json:"productOptionIds"`
	ShopId           uint64            `gorm:"type:bigint(20)" json:"shopId"`
	//OptionList       []*ProductOption  `gorm:"ForeignKey:ProductId;AssociationForeignKey:ProductId" json:"OptionList"`
}

//这个是包括option和sku的
type ProductAllDetail struct {
	Product *Product                   `json:"product"`
	Option  map[uint8][]*ProductOption `json:"optionList"`
	Sku     []*Sku                     `json:"skuList"`
	Image   []string                   `json:"imageList"`
}

type ProductPriceObj struct {
	Product *Product `json:"product"`
	Sku     *Sku     `json:"sku"`
}

func (ProductDetail) TableName() string {
	//return constant.ProductDetailTable
	tableName := helper.GetDBByName(constant.DBMerchant).Model(&ProductDetail{}).Name()
	return tableName
}

// 表结构初始化
func init() {
	// 创建或更新表结构
	_ = helper.GetDBByName(constant.DBMerchant).AutoMigrate(&Product{})
}

// 插入前生成主键
func (product *Product) BeforeCreate(db *gorm.DB) error {
	//id := uuid.NewV4()
	//db.Set("ID", &id)
	//product.ID = id.String()
	return nil
}

// 校验表单中提交的参数是否合法
func (product *Product) Validator() error {
	return nil
}
