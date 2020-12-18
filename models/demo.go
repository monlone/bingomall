package model

import (
	"bingomall/constant"
	helper "bingomall/helpers"
	"gorm.io/gorm"
)

// 产品分类结构体
type Demo struct {
	/** 主键id*/
	Model

	/** id */
	//DemoId string `gorm:"type:varchar(36);" form:"demoId"`

	/** 商店id */
	ShopId uint64 `gorm:"type:bigint(20);" form:"shopId"`

	/** 排序 */
	Order uint `gorm:"type:int;" form:"order"`

	/** 分类名称 */
	Name string `gorm:"type:varchar(36);" form:"name"`

	/** 分类描述 */
	Desc string `gorm:"type:varchar(36);" form:"desc"`

	/** 添加者的userId，为了好追踪 */
	UserID uint64 `gorm:"type:bigint;" form:"userId"`

	ExpiredTime
	CrudTime
}

// 插入前生成主键
func (demo *Demo) BeforeCreate(db *gorm.DB) error {
	//id := uuid.NewV4()
	//db.Set("ID", &id)
	//demo.DemoId = id.String()
	return nil
}

func init() {
	// 创建或更新表结构
	_ = helper.GetDBByName(constant.DBMerchant).AutoMigrate(&Demo{})
}
