package model

import (
	"bingomall/constant"
	helper "bingomall/helpers"
	"gorm.io/gorm"
)

// 产品分类结构体
type OrderProduct struct {
	Model

	/** 定单id */
	OrderId uint64 `gorm:"type:bigint;" form:"orderId"`

	/** 商店id */
	ShopId uint64 `gorm:"type:bigint;" form:"shopId" json:"shopId"`

	/** 产品id */
	ProductId uint64 `gorm:"type:bigint;" form:"productId" json:"productId"`

	/** 下单时的产品价格*/
	Price uint64 `gorm:"type:bigint(20)" json:"price"`

	/**积分*/
	Score uint64 `gorm:"type:bigint(20)" json:"score"`

	/** 购买个数*/
	Number uint64 `gorm:"type:bigint(20)" json:"number"`

	/** 商品的价格*/
	Money uint64 `gorm:"type:bigint(20)" json:"money" form:"money"`

	/**0:待支付，1：已支付，待回调确认，2：支付成功，已回调确认，3：线下商品商户已经核销， 4：已经发货，5：已收货，6：纠纷中，7：已完成*/
	Status int8 `gorm:"type:tinyint(3);" json:"status"`

	SkuId uint64 `gorm:"type:bigint;" json:"skuId"`

	Type uint8 `gorm:"type:tinyint(3);" json:"type"`

	Shop *Shop `gorm:"ForeignKey:ShopId;" json:"shop"`

	Product *Product `gorm:"ForeignKey:ProductId;" json:"product"`

	CrudTime
}

type OrderDetail struct {
	OrderId uint64 `gorm:"type:bigint(20);" form:"orderId"`
	Number  uint64 `gorm:"type:bigint(20)" json:"number"`
}

func (OrderDetail) TableName() string {
	tableName := helper.GetDBByName(constant.DBMerchant).Model(&OrderProduct{}).Name()
	return tableName
}

// 插入前生成主键
func (orderProduct *OrderProduct) BeforeCreate(db *gorm.DB) error {
	//id := uuid.NewV4()
	//db.Set("ID", &id)
	//orderProduct.OrderProductId = id.String()
	return nil
}

func init() {
	// 创建或更新表结构
	_ = helper.GetDBByName(constant.DBMerchant).AutoMigrate(&OrderProduct{})
}
