package model

import (
	"bingomall/constant"
	helper "bingomall/helpers"
	"gorm.io/gorm"
)

// 产品分类结构体
type Sku struct {
	Model

	//SkuId string `gorm:"type:varchar(36);" form:"skuId" json:"skuId"`

	/** 产品id */
	ProductId uint64 `gorm:"type:bigint;not null;" form:"productId" json:"productId"`

	/** 产品option id */
	//ProductOptionId uint `gorm:"type:int;" form:"productOptionId" json:"productOptionId"`

	/** option 自定义描述 */
	Desc string `gorm:"type:varchar(36);" form:"desc" json:"desc"`

	Stock uint64 `gorm:"type:bigint" form:"stock" json:"stock"`

	/** 一个sku对应的产品价格*/
	Price uint64 `gorm:"type:bigint" form:"price" json:"price"`

	/** 以option的primary id 按从小到在的方式组合，加上productId一起变成一个唯一值，用来给C端用户查价格和库存*/
	CombineId string `gorm:"type:varchar(255)" form:"combineId" json:"combineId"`

	/** 排序 */
	Order uint `gorm:"type:int;" form:"order" json:"order"`

	/** 添加者的userId，为了好追踪 */
	UserID uint64 `gorm:"type:bigint;" form:"userId" json:"userId"`

	//Product Product `gorm:"foreignKey:ProductId;" json:"product"`

	//OptionList []*ProductOption `gorm:"ForeignKey:ProductId;AssociationForeignKey:ProductId" json:"OptionList"`
	CrudTime
}

// 插入前生成主键
func (sku *Sku) BeforeCreate(db *gorm.DB) error {
	//id := uuid.NewV4()
	//db.Set("ID", &id)
	//sku.SkuId = id.String()
	return nil
}

func init() {
	// 创建或更新表结构
	_ = helper.GetDBByName(constant.DBMerchant).AutoMigrate(&Sku{})
}
