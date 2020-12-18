package constant

const (
	OrderWaitedForPay   = 0
	OrderWaitedForCheck = 1 //已调起过支付，等支付平台待回调
	OrderPaySuccess     = 2 //已回调，支付成功
	OrderVerified       = 3 //商户已核销
	OrderShipped        = 4 //已发货
	OrderReceived       = 5 //已收货

	MerchantCheckUserPay = 1 //1:用户核销(贷入)2：和平台结算(贷出)

	MerchantOrderWaitedForLiquidate = 1 //未结算
	MerchantOrderLiquidated         = 2 //已结算

	UserCost  = 2 //用户消费
	LevelCost = 1 //下线消费

	OrderLevel = 3 //下线消费的提成

	OrderRealGoods    = 1 /**商品类型 1-实物商品要核销, 2-虚拟商品要核销，3-虚拟商品不用核销，4-实物商品要邮寄 5-实物商品不要邮寄 */
	OrderVirtualGoods = 2

	MoneyLogMerchant = 1
	MoneyLogUser     = 2
	MoneyLogSelf     = 3
	MoneyPay         = 4 //1:商户代核销入账，2：下线消费入账, 3:自己消费返现, 4:用积分购买商品，5:商户自卖产品入账，6:商户提现

	FirstLevelProfit  = 0.6 //一级上级分discount_level的60%
	SecondLevelProfit = 0.4 //二级上级分discount_level的40%

	//https://pay.weixin.qq.com/wiki/doc/api/allocation_sl.php?chapter=25_6&index=2
	MultiProfitSharing  = "https://api.mch.weixin.qq.com/secapi/pay/multiprofitsharing"
	ProfitSharingFinish = "https://api.mch.weixin.qq.com/secapi/pay/profitsharingfinish"
)
