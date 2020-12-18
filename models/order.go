package model

import (
	"bingomall/constant"
	helper "bingomall/helpers"
	"gorm.io/gorm"
)

type Order struct {
	Model

	//OrderId    string `gorm:"type:varchar(36);column:order_id;" json:"orderId" form:"orderId"`
	OutTradeNo string `gorm:"type:varchar(32);unique" form:"out_trade_no" json:"-"`

	/**微信支付订单号*/
	TransactionID string `gorm:"type:varchar(32);column:transaction_id" json:"transactionID"`

	/** 用户购买的代金券的张数，要是大于0就要用orderId去代金券表里查*/
	VoucherNumber uint `gorm:"type:int(10)" json:"voucherNumber" form:"voucherNumber"`

	/** 总金额*/
	TotalAmount uint64 `gorm:"type:bigint(20)" json:"totalAmount" form:"totalAmount"`

	/** 运费*/
	LogisticsAmount uint64 `gorm:"type:bigint(20)" json:"logisticsAmount" form:"logisticsAmount"`

	/**使用的积分*/
	Score uint64 `gorm:"type:bigint(20)" json:"totalScoreToPay"`

	/** 用户使用代金券后商品的应该支付的金额*/
	Pay uint64 `gorm:"type:bigint(20)" json:"pay"`

	PrepayID string `gorm:"type:varchar(36);column:prepay_id" json:"prepayId"`

	//0:待支付，1：已支付，待回调确认，2：支付成功，已回调确认，3:线下商品商户已经核销，4：已经发货，5：已收货，6：申请退款，7:已退款，8：纠纷中，9：已关闭，100：已完成
	Status int8 `gorm:"type:tinyint(3);" json:"status"`

	Platform string `gorm:"type:varchar(10)" json:"platform"`

	/** 购买者userId，或者推荐者的userId*/
	UserID uint64 `gorm:"type:bigint;column:user_id;" form:"userId" json:"userId" weChat:"user_id"`

	Describe string `gorm:"type:varchar(255)" json:"describe"`

	/** 定单类型，与商品的价格类型要保持一致 1-普通价格商品，2-秒杀商品，3-砍价商品，4-拼团商品*/
	Type uint8 `gorm:"type:tinyint(1)" json:"type"`

	//Product Product `gorm:"foreignKey:ID;" json:"product"`

	OrderProduct []*OrderProduct `gorm:"foreignKey:OrderId" json:"orderProduct"`

	Total int `json:"total"`

	CrudTime
}

type Cost struct {
	OrderId uint64 `form:"orderId"`
	/** 用户消费的代金券张数*/
	Number uint64 `form:"number"`

	/** 用户剩下的代金券金额或者剩下商品的总金额*/
	Money uint64 `form:"money"`

	ShopId   uint64 `form:"shopId"`
	Platform string `form:"platform"`
	/** 购买者userId*/
	UserID uint64 `form:"userId"`
	/**核销店员user_id*/
	CheckUserID uint64 `form:"checkUserID"`
	ProductId   string `form:"productId"`
}

type Pay struct {
	Score      uint64   `form:"score" json:"score"`
	ProductIds []string `form:"productIds" json:"productIds"`
	OrderId    uint64   `form:"orderId" json:"orderId"`
	CouponList []string `form:"couponList" json:"couponList"`
}

type OrderListObject struct {
	OrderId uint64 `form:"orderId" json:"orderId"`
	ShopId  uint64 `form:"shopId" json:"shopId"`
	PageObject
}

type OrderReturn struct {
	Order
	IsNeedLogistics bool        `json:"isNeedLogistics"`
	CouponList      interface{} `json:"couponList"`
}

type ProductOrder struct {
	ProductId      uint64 `json:"productId" form:"productId"`
	Number         uint64 `json:"number" form:"number"` //为了更方便计算用了个大的，本来应该用不到uint64
	LogisticsType  string `json:"logisticsType" form:"logisticsType"`
	InviterId      uint64 `json:"inviterId" form:"inviterId"`
	ShoppingCartId uint64 `json:"shoppingCartId"`
}

// 表结构初始化
func init() {
	// 创建或更新表结构
	_ = helper.GetDBByName(constant.DBMerchant).AutoMigrate(&Order{})
	// 生成外键约束
	//helper.SQL.Model(&Order{}).AddForeignKey("role_id", "role(id)", "no action", "no action")
}

// 插入前生成主键
func (order *Order) BeforeCreate(db *gorm.DB) error {
	//id := uuid.NewV4()
	//db.Set("ID", &id)
	//order.ID = id.String()
	return nil
}

// 校验表单中提交的参数是否合法
func (order *Order) Validator() error {
	return nil
}
