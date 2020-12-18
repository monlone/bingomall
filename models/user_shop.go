package model

import (
	"bingomall/constant"
	helper "bingomall/helpers"
	"gorm.io/gorm"
)

// 这个用于添加商户和官方店的关联，商户的分销从这个表来找，然后分走指定的提成
type UserShop struct {
	Model
	//UserShopId string `gorm:"type:varchar(36);column:user_shop_id;not null;" json:"user_shop_id"`
	UserID uint64 `gorm:"type:bigint;column:user_id;" form:"user_id" json:"user_id" weChat:"user_id"`
	ShopId uint64 `gorm:"type:bigint;column:shop_id;not null;" json:"shop_id"`
	CrudTime
}
type UserShopListsResponse struct {
	helper.JsonObject
	Data []UserShop
}

// 表结构初始化
func init() {
	// 创建或更新表结构
	_ = helper.GetDBByName(constant.DBMerchant).AutoMigrate(&UserShop{})
	// 生成外键约束
	//helper.SQL.Model(&UserShop{}).AddForeignKey("role_id", "role(id)", "no action", "no action")
}

// 插入前生成主键
func (userShop *UserShop) BeforeCreate(db *gorm.DB) error {
	//id := uuid.NewV4()
	//db.Set("ID", &id)
	//userShop.UserShopId = id.String()
	return nil
}

// 校验表单中提交的参数是否合法
func (userShop *UserShop) Validator() error {
	return nil
}
