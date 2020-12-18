package model

import (
	"bingomall/constant"
	helper "bingomall/helpers"
	"gorm.io/gorm"
)

// 省份结构体
type City struct {
	Model

	/**  */
	Code uint64 `gorm:"type:bigint;" form:"code" json:"code"`

	/** 名称 */
	Name string `gorm:"type:varchar(36);" form:"name" json:"name"`

	ProvinceCode string `gorm:"type:varchar(36);" form:"provinceCode" json:"provinceCode"`
}

// 插入前生成主键
func (city *City) BeforeCreate(db *gorm.DB) error {
	//id := uuid.NewV4()
	//db.Set("ID", &id)
	//city.CityID = id.String()
	return nil
}

func init() {
	// 创建或更新表结构
	_ = helper.GetDBByName(constant.DBMerchant).AutoMigrate(&City{})
}
