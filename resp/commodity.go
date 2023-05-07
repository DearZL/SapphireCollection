package resp

type Commodity struct {
	Hash    []byte  `json:"hash"`
	Image   string  `json:"image"`
	Name    string  `json:"name"`
	Price   float32 `json:"price"`
	OrderID string  `json:"orderID"`
}
type Commodities struct {
	Commodities []*Commodity
}
