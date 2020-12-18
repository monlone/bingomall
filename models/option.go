package model

import (
	"bingomall/constant"
	helper "bingomall/helpers"
	"gorm.io/gorm"
)

// 结构体
type Option struct {
	/** 主键id*/
	Model

	/** id */
	//OptionId string `gorm:"type:varchar(36);" form:"optionId" json:"optionId"`

	/** option描述 */
	Desc string `gorm:"type:varchar(36);" form:"desc" json:"desc"`

	/** option的图片地址*/
	ImageUrl string `gorm:"type:varchar(255);" form:"imageUrl" json:"imageUrl"`

	/** option的类别 1：颜色，2：尺寸 3：其他*/
	Type uint8 `gorm:"type:tinyint" form:"type" json:"type"`

	/** 排序 */
	Order uint `gorm:"type:int;" form:"order" json:"order"`

	/** 添加者的userId，为了好追踪 */
	UserID uint64 `gorm:"type:bigint;" form:"userId" json:"userId"`

	CrudTime
}

// 插入前生成主键
func (option *Option) BeforeCreate(db *gorm.DB) error {
	//id := uuid.NewV4()
	//db.Set("ID", &id)
	//option.OptionId = id.String()
	return nil
}

func init() {
	// 创建或更新表结构
	_ = helper.GetDBByName(constant.DBMerchant).AutoMigrate(&Option{})
}
