package enum

const (
	// CommodityDisabled 已禁止售卖
	CommodityDisabled = true
	// CommodityActive 已开启售卖
	CommodityActive = false

	// OrderTypeWallet 钱包类型订单
	OrderTypeWallet = "wallet"
	// OrderTypeCommodity 商品类型订单
	OrderTypeCommodity = "commodity"

	// OrderStatusCancelled 已取消
	OrderStatusCancelled = 0
	// OrderStatusUnpaid 未支付
	OrderStatusUnpaid = 1
	// OrderStatusPaid 已支付
	OrderStatusPaid = 2
	// OrderStatusFinish 已完成
	OrderStatusFinish = 3

	// UserDisabled 已禁用
	UserDisabled = true
	// UserActive 已启用
	UserActive = false
)
