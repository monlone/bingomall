package model

import (
	"bingomall/constant"
	helper "bingomall/helpers"
	"gorm.io/gorm"
)

type MerchantOrderDetail struct {
	Model

	OrderId uint64 `gorm:"type:bigint;column:order_id;not null;" json:"orderId"`

	MerchantId uint64 `gorm:"type:bigint;column:merchant_id;not null;" json:"merchantId"`

	Number uint64 `gorm:"type:int(10)" json:"number"`

	/**商户的收入*/
	Money uint64 `gorm:"type:bigint(20)" json:"money"`

	ShopId uint64 `gorm:"type:bigint;column:shop_id" json:"shopId"`

	/**status  1:待核销，2:已核销*/
	Status uint8 `gorm:"type:tinyint(3);" json:"status"`

	Platform string `gorm:"type:varchar(10)" json:"platform"`

	/** 购买者userId*/
	UserID uint64 `gorm:"type:bigint;column:user_id;" form:"user_id" json:"userId" weChat:"userId"`

	/** 核销的店员用户id*/
	CheckUserID uint64 `gorm:"type:bigint;column:check_user_id;" form:"checkUserID" json:"checkUserID"`

	ProductId uint64 `gorm:"type:bigint;column:product_id;" json:"productId"`
	Describe  string `gorm:"type:varchar(255)" json:"describe"`

	/** 1:用户核销(贷入)2：和平台结算(贷出) */
	Type uint8 `gorm:"type:tinyint(1)" json:"type"`

	/** 1:实物，2：代金券 对应order的type*/
	GoodsType uint8 `gorm:"type:tinyint(1)" json:"goodsType"`

	//ShopDetail    *ShopDetail    `gorm:"ForeignKey:ShopId;AssociationForeignKey:ShopId" json:"shop"`
	ShopDetail Shop `gorm:"foreignKey:ShopId" json:"shop"`
	//ProductDetail *ProductDetail `gorm:"ForeignKey:ProductId;AssociationForeignKey:ProductId" json:"product"`
	ProductDetail Product `gorm:"foreignKey:ProductId" json:"product"`

	CrudTime
}

type TotalResult struct {
	TotalMoney uint64 `json:"total_money"`
}

type MerchantOrderDetailPage struct {
	Status uint8  `from:"status" json:"status"`
	Month  string `from:"status" json:"status"`
	PageObject
}

type MerchantOrderListsResponse struct {
	helper.JsonObject
	Data []MerchantOrderDetail
}

// 表结构初始化
func init() {
	// 创建或更新表结构
	_ = helper.GetDBByName(constant.DBMerchant).AutoMigrate(&MerchantOrderDetail{})
	// 生成外键约束
	//helper.SQL.Model(&MerchantOrderDetail{}).AddForeignKey("role_id", "role(id)", "no action", "no action")
}

// 插入前生成主键
func (order *MerchantOrderDetail) BeforeCreate(db *gorm.DB) error {
	//id := uuid.NewV4()
	//db.Set("ID", &id)
	//order.ID = id.String()
	return nil
}

// 校验表单中提交的参数是否合法
func (order *MerchantOrderDetail) Validator() error {
	return nil
}
