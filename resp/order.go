package resp

type Order struct {
	OrderNum    string      `json:"orderNum"`
	SellerId    string      `json:"sellerId"`
	BuyerId     string      `json:"buyerId"`
	Commodities Commodities `json:"commodities"`
	OrderAmount float32     `json:"orderAmount"`
}

func (o *Order) CreateReOrder(order *Order, commodities Commodities) {
	*o = Order{
		OrderNum:    order.OrderNum,
		SellerId:    order.SellerId,
		BuyerId:     order.BuyerId,
		Commodities: commodities,
		OrderAmount: order.OrderAmount,
	}
}
