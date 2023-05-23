package resp

type User struct {
	UserId     string           `json:"userId"`
	UserName   string           `json:"userName" `
	Email      string           `json:"email"`
	Icon       string           `json:"icon"`
	Collection []*UserCommodity `json:"collection"`
}
