package model

import (
	"bingomall/constant"
	helper "bingomall/helpers"
	"gorm.io/gorm"
)

type Banner struct {
	Model

	/** banner id */
	//BannerID string `gorm:"type:varchar(36);column:banner_id;not null;" json:"bannerId" form:"banner_id"`

	/** 1是平台自己，2是商户 */
	Type int8 `gorm:"type:tinyint(3);" form:"type" json:"type"`

	/** 店铺id */
	ShopId string `gorm:"type:varchar(36)" json:"shopId" form:"shop_id" binding:"required"`

	/** banner url */
	ImageUrl string `gorm:"type:varchar(5000);" form:"image_url" binding:"required" json:"imageUrl"`

	/** banner 要跳转的地址 */
	ToUrl string `gorm:"type:varchar(255)" form:"to_url" json:"toUrl"`

	/** 跳转类型，1：web 2：app */
	ToType int8 `gorm:"type:tinyint(3);default:1" form:"to_type" json:"toType"`

	/** banner title */
	Title string `gorm:"type:varchar(255);" form:"title" json:"title"`

	/** banner 的描述 */
	Content string `gorm:"type:varchar(255)" json:"content"`

	Duration

	CrudTime
}

// 表结构初始化
func init() {
	// 创建或更新表结构
	_ = helper.GetDBByName(constant.DBMerchant).AutoMigrate(&Banner{})
}

// 插入前生成主键
func (banner *Banner) BeforeCreate(db *gorm.DB) error {
	//id := uuid.NewV4()
	//db.Set("ID", &id)
	//banner.BannerID = id.String()
	return nil
}

// 校验表单中提交的参数是否合法
func (banner *Banner) Validator() error {
	return nil
}
