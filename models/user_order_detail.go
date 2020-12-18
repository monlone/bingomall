package model

import (
	"bingomall/constant"
	helper "bingomall/helpers"
	"gorm.io/gorm"
)

type UserOrderDetail struct {
	Model
	OrderId    uint64 `gorm:"type:bigint;column:order_id;not null;" json:"order_id"`
	MerchantId uint64 `gorm:"type:bigint;column:merchant_id;not null;" json:"merchant_id"`
	Number     uint64 `gorm:"type:int(10)" json:"number"`
	Money      uint64 `gorm:"type:bigint(20)" json:"money"`
	ShopId     uint64 `gorm:"type:bigint;column:shop_id" json:"shop_id"`
	/**1: */
	Status   uint8  `gorm:"type:tinyint(3);" json:"status"`
	Platform string `gorm:"type:varchar(10)" json:"platform"`

	/** 购买者userId，或者推荐者的userId*/
	UserID uint64 `gorm:"type:bigint;column:user_id;" form:"user_id" json:"user_id" weChat:"user_id"`

	/** 核销店员的userId，或者下线的userId*/
	CheckUserID uint64 `gorm:"type:bigint;column:check_user_id;" form:"check_user_id" json:"check_user_id"`

	ProductId uint64 `gorm:"type:bigint;column:product_id;" json:"product_id"`
	Describe  string `gorm:"type:varchar(255)" json:"describe"`

	/** 1:用户核销(贷入)2：和平台结算(贷出) */
	Type uint8 `gorm:"type:tinyint(1)" json:"type"`

	/** 1:实物，2：代金券 3.下线消费提成 对应order的type*/
	GoodsType uint8 `gorm:"type:tinyint(1)" json:"goods_type"`

	Shop    *Shop    `gorm:"ForeignKey:ShopId;association_foreignkey:ShopId" json:"shop"`
	Product *Product `gorm:"ForeignKey:ProductId;association_foreignkey:ProductId" json:"product"`

	CrudTime
}
type UserOrderListsResponse struct {
	helper.JsonObject
	Data []UserOrderDetail
}

// 表结构初始化
func init() {
	// 创建或更新表结构
	_ = helper.GetDBByName(constant.DBMerchant).AutoMigrate(&UserOrderDetail{})
	// 生成外键约束
	//helper.SQL.Model(&UserOrderDetail{}).AddForeignKey("role_id", "role(id)", "no action", "no action")
}

// 插入前生成主键
func (order *UserOrderDetail) BeforeCreate(db *gorm.DB) error {
	//id := uuid.NewV4()
	//db.Set("ID", &id)
	//order.ID = id.String()
	return nil
}

// 校验表单中提交的参数是否合法
func (order *UserOrderDetail) Validator() error {
	return nil
}
