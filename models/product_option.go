package model

import (
	"bingomall/constant"
	helper "bingomall/helpers"
	"gorm.io/gorm"
)

// 产品分类结构体
type ProductOption struct {
	/** 主键id*/
	Model

	/** id */
	//ProductOptionId string `gorm:"type:varchar(36);" form:"productOptionId" json:"productOptionId"`

	/** 产品id */
	ProductId uint64 `gorm:"type:bigint;not null;" form:"productId" json:"productId"`

	/** optionId 与option表关联，一个产品对应多个option，比如：红色，黄色，大码，小码*/
	OptionId uint64 `gorm:"type:bigint;" form:"optionId" json:"optionId"`

	/** option中的类别 1：颜色，2：尺寸 3：其他，这个是冗余存的，不想连表查了*/
	Type uint8 `gorm:"type:tinyint" form:"type" json:"type"`

	/** option 自定义描述 */
	Desc string `gorm:"type:varchar(36);" form:"desc" json:"desc"`

	/** 自定义option的图片地址*/
	ImageUrl string `gorm:"type:varchar(255);" form:"imageUrl" json:"imageUrl"`

	/** 排序 */
	Order uint `gorm:"type:int;" form:"order" json:"order"`

	/** 添加者的userId，为了好追踪 */
	UserID uint64 `gorm:"type:bigint;" form:"userId" json:"userId"`

	Option *Option `gorm:"foreignKey:OptionId;" json:"option"`

	CrudTime
}

// 插入前生成主键
func (productOption *ProductOption) BeforeCreate(db *gorm.DB) error {
	//id := uuid.NewV4()
	//db.Set("ID", &id)
	//productOption.ProductOptionId = id.String()
	return nil
}

func init() {
	// 创建或更新表结构
	_ = helper.GetDBByName(constant.DBMerchant).AutoMigrate(&ProductOption{})
}
