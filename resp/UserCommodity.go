package resp

type UserCommodity struct {
	Hash  string `json:"hash"`
	Image string `json:"image"`
	Name  string `json:"name"`
}
type UserCommodities struct {
	UserCommodities []*UserCommodity
}
