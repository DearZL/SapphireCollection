package resp

type Order struct {
	OrderNum    string      `json:"orderNum"`
	SellerId    string      `json:"sellerId"`
	BuyerId     string      `json:"buyerId"`
	Commodities Commodities `json:"commodities"`
	OrderAmount float32     `json:"orderAmount"`
	OrderTime   string      `json:"orderTime"`
}
