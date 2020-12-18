package model

import (
	"bingomall/constant"
	helper "bingomall/helpers"
	"gorm.io/gorm"
)

type BaseShoppingCart struct {
	Model

	/** shoppingCart id */
	//ShoppingCartId string `gorm:"type:varchar(36);column:shopping_cart_id;not null;" json:"shoppingCartId" form:"shoppingCartId"`

	ProductId uint64 `gorm:"type:bigint;column:product_id;not null;" json:"productId" form:"productId"`

	Number uint64 `gorm:"type:int(10)" json:"number" form:"number"`

	Price uint64 `gorm:"type:varchar(30);" form:"price" json:"price"`

	SkuId uint64 `gorm:"type:bigint" form:"skuId" json:"skuId"` //反正要去sku里面拿库存，CombineId就没有必要存在了

	UserID uint64 `gorm:"type:bigint" form:"userId" json:"userId"`

	//Product *Product `gorm:"hasOne;ForeignKey:ProductId;AssociationForeignKey:ProductId" json:"product"`
	Product *Product `gorm:"foreignKey:ProductId" json:"product"`

	Sku *Sku `gorm:"foreignKey:SkuId" json:"sku"`
}

type ShoppingCart struct {
	BaseShoppingCart

	OptionList []*ProductOption `gorm:"foreignKey:ProductId;references:ProductId" json:"optionList"`

	CrudTime
}

type ShoppingCartPost struct {
	ProductId        string   `json:"productId"`
	Number           uint64   `json:"number"`
	ProductOptionIds []string `json:"productOptionIds"`
}

type ShoppingCartInfo struct {
	BaseShoppingCart
	OptionList map[uint8][]*ProductOption `json:"optionList"`
}

// 表结构初始化
func init() {
	// 创建或更新表结构
	_ = helper.GetDBByName(constant.DBMerchant).AutoMigrate(&ShoppingCart{})
}

// 插入前生成主键
func (sCart *ShoppingCart) BeforeCreate(db *gorm.DB) error {
	//id := uuid.NewV4()
	//db.Set("ID", &id)
	//sCart.ShoppingCartId = id.String()
	return nil
}

// 校验表单中提交的参数是否合法
func (sCart *ShoppingCart) Validator() error {
	return nil
}
