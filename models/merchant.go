package model

import (
	"bingomall/constant"
	helper "bingomall/helpers"
	"gorm.io/gorm"
)

type Merchant struct {
	Model

	/** merchant id */
	//MerchantId string `gorm:"type:varchar(36);column:merchant_id;" json:"merchant_id" form:"merchant_id"`

	UserID uint64 `gorm:"type:bigint(20);column:user_id;" json:"userId" form:"userId"`

	/** merchant title */
	Title string `gorm:"type:varchar(255);" form:"title" json:"title"`

	/** 商户电话 */
	Phone string `gorm:"type:varchar(20);" form:"phone" binding:"required" json:"phone"`

	/** 商户描述 */
	Description string `gorm:"type:varchar(1000);" form:"description" json:"description"`

	Creator string `gorm:"type:varchar(20);" form:"creator" json:"-"`

	/** 商户状态，0：不可用，1：可用*/
	Status int8 `gorm:"type:tinyint(3);default:0" form:"status" json:"status"`

	ShopList []Shop `gorm:"ForeignKey:MerchantId;" json:"shopList"`

	CrudTime
}

type MerchantSummary struct {
	MerchantId string `gorm:"type:varchar(36);column:merchant_id;" json:"merchantId"`
	Title      string `gorm:"type:varchar(255); column:title;" json:"title"`
	Phone      string `gorm:"type:varchar(20); column:phone" json:"phone"`
	Logo       string `gorm:"type:varchar(500); column:logo;" json:"logo"`
}

func (MerchantSummary) TableName() string {
	tableName := helper.GetDBByName(constant.DBMerchant).Model(&MerchantSummary{}).Name()
	return tableName
}

// 表结构初始化
func init() {
	// 创建或更新表结构
	_ = helper.GetDBByName(constant.DBMerchant).AutoMigrate(&Merchant{})
}

// 插入前生成主键
func (merchant *Merchant) BeforeCreate(db *gorm.DB) error {
	//id := uuid.NewV4()
	//db.Set("ID", &id)
	//merchant.MerchantId = id.String()
	return nil
}

// 校验表单中提交的参数是否合法
func (merchant *Merchant) Validator() error {
	return nil
}
