package resp

type Commodity struct {
	Hash     string  `json:"hash"`
	Image    string  `json:"image"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Amount   int     `json:"amount"`
	OrderNum string  `json:"orderNum"`
}
type Commodities struct {
	Commodities []*Commodity
}
