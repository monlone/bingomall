package model

import (
	"bingomall/constant"
	helper "bingomall/helpers"
	"gorm.io/gorm"
)

/*资金、积分明细表*/
type MoneyLog struct {
	Model

	//MoneyLogID string `gorm:"type:varchar(36);column:money_log_id;not null;" json:"money_log_id"`

	/**微信支付订单号*/
	TransactionID string `gorm:"type:varchar(32);column:transaction_id" json:"transaction_id"`

	/** 用户消费的金额*/
	Cost uint64 `gorm:"type:bigint(20)" json:"cost" form:"cost"`

	/** 本次上线返现，商户的收益金额,或者自己的返现*/
	Income uint64 `gorm:"type:bigint(20)" json:"income" form:"income"`

	/** 用户花费的积分*/
	CostScore uint64 `gorm:"type:bigint(20)" json:"costScore" form:"costScore"`

	/** 用户赚到的积分*/
	Score uint64 `gorm:"type:bigint(20)" json:"score" form:"score"`

	/** 这条记录的userId*/
	UserID uint64 `gorm:"type:bigint;column:user_id;" form:"user_id" json:"user_id"`

	/**相关用户id，merchant的userId或者上线的userId,或者merchantId*/
	RelationUserID uint64 `gorm:"type:bigint;column:relation_user_id;" form:"relation_user_id" json:"relation_user_id"`

	OrderId  uint64 `gorm:"type:bigint;column:order_id;" json:"orderId"`
	Describe string `gorm:"type:varchar(255)" json:"describe"`
	/** 1:商户代核销入账，2：下线消费入账, 3:自己消费返现, 4:用积分购买商品，5:商户自卖产品入账，6:商户提现，7：自己购买产品 */
	Type uint8 `gorm:"type:tinyint(1)" json:"type"`

	/**1:等待线下支付，2:已经线下支付*/
	//Status uint8 `gorm:"type:tinyint(1);default:1" form:"status" json:"status"`

	Merchant Merchant `gorm:"ForeignKey:RelationUserID" json:"merchant"`
	CrudTime
}

// 表结构初始化
func init() {
	// 创建或更新表结构
	_ = helper.GetDBByName(constant.DBMerchant).AutoMigrate(&MoneyLog{})
}

// 插入前生成主键
func (moneyLog *MoneyLog) BeforeCreate(db *gorm.DB) error {
	//id := uuid.NewV4()
	//db.Set("ID", &id)
	//moneyLog.ID = id.String()
	return nil
}

// 校验表单中提交的参数是否合法
func (moneyLog *MoneyLog) Validator() error {
	return nil
}
