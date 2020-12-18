package model

import (
	"bingomall/constant"
	helper "bingomall/helpers"
	"gorm.io/gorm"
)

// 功能菜单结构体
type Wallet struct {
	Model

	/** 主键 id */
	//WalletID string `gorm:"type:varchar(36);" form:"walletId"`

	/** 用户id */
	UserID uint64 `gorm:"type:bigint;" form:"userId" json:"userId"`

	/** 用户积分 */
	Score uint64 `gorm:"type:bigint;" form:"score" json:"score"`

	/** 用户总金额*/
	Money uint64 `gorm:"type:bigint;" form:"money" json:"money"`

	/** 用户余额,总额度-冻结金额 */
	Balance uint64 `gorm:"type:bigint;" form:"balance" json:"balance"`

	/** 用户冻结金额*/
	Freeze uint64 `gorm:"type:bigint;" form:"freeze" json:"freeze"`

	/** 用户成长值*/
	Growth uint64 `gorm:"type:bigint;" form:"money" json:"growth"`

	ExpiredTime
	CrudTime
}

// 插入前生成主键
func (wallet *Wallet) BeforeCreate(db *gorm.DB) error {
	//id := uuid.NewV4()
	//db.Set("ID", &id)
	//wallet.WalletID = id.String()
	return nil
}

func init() {
	// 创建或更新表结构
	_ = helper.GetDBByName(constant.DBMerchant).AutoMigrate(&Wallet{})
}
