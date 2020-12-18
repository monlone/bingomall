package model

import (
	"bingomall/constant"
	helper "bingomall/helpers"
	"gorm.io/gorm"
)

// 产品分类结构体
type Address struct {
	Model

	/** 主键 id */
	//AddressID string `gorm:"type:varchar(36);" form:"addressId" json:"addressId"`

	/** 添加者的userId*/
	UserID uint64 `gorm:"type:bigint;" form:"userId" json:"-"`

	Phone string `gorm:"type:varchar(20);" form:"phone" json:"phone"`

	Contact string `gorm:"type:varchar(20);" form:"contact" json:"contact"`

	/** 排序 */
	Order uint `gorm:"type:int;" form:"order"`

	/** 名称 */
	Province string `gorm:"type:varchar(36);" form:"province" json:"province"`
	
	ProvinceCode string `gorm:"type:varchar(36);" form:"provinceCode" json:"provinceCode"`

	City     string `gorm:"type:varchar(36);" form:"city" json:"city"`
	CityCode string `gorm:"type:varchar(36);" form:"cityCode" json:"cityCode"`

	Area     string `gorm:"type:varchar(36);" form:"area" json:"area"`
	AreaCode string `gorm:"type:varchar(36);" form:"areaCode" json:"areaCode"`

	AddressDetail string `gorm:"type:varchar(40);" form:"address" json:"address"`

	/** 描述 */
	Desc string `gorm:"type:varchar(36);" form:"desc"`

	IsDefault bool `gorm:"type:tinyint" form:"isDefault" json:"isDefault"`

	CrudTime
}

// 插入前生成主键
func (address *Address) BeforeCreate(db *gorm.DB) error {
	//id := uuid.NewV4()
	//db.Set("ID", &id)
	//address.ID = id.String()
	return nil
}

func init() {
	// 创建或更新表结构
	_ = helper.GetDBByName(constant.DBMerchant).AutoMigrate(&Address{})
}
