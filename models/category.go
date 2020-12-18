package model

import (
	"bingomall/constant"
	helper "bingomall/helpers"
	"gorm.io/gorm"
)

// 产品分类结构体
type Category struct {
	Model

	/** 主键 id */
	//CategoryId string `gorm:"type:varchar(36);" form:"categoryId" json:"categoryId"`

	/** 商店id */
	ShopId uint64 `gorm:"type:bigint(20);" form:"shopId" json:"shopId"`

	/** 分类排序 */
	Order uint `gorm:"type:int;" form:"order" json:"order"`

	/** 分类名称 */
	Name string `gorm:"type:varchar(36);" form:"name" json:"name"`

	ImageUrl string `gorm:"type:varchar(500);" form:"imageUrl" json:"imageUrl"`

	/** 分类层级*/
	Level uint8 `gorm:"type:smallint" form:"level" json:"level"`

	/** 分类描述 */
	Desc string `gorm:"type:varchar(36);" form:"desc" json:"desc"`

	/** 添加者的userId，为了好追踪 */
	UserID uint64 `gorm:"type:bigint;" form:"userId" json:"-"`

	ExpiredTime
	CrudTime
}

// 插入前生成主键
func (category *Category) BeforeCreate(db *gorm.DB) error {
	//id := uuid.NewV4()
	//db.Set("ID", &id)
	//category.CategoryId = id.String()
	return nil
}

func init() {
	// 创建或更新表结构
	_ = helper.GetDBByName(constant.DBMerchant).AutoMigrate(&Category{})
}
