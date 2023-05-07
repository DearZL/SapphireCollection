package resp

type EntityA struct {
	Code  int         `json:"code"`
	Msg   string      `json:"msg"`
	Token string      `json:"token"`
	Data  interface{} `json:"data"`
}
type EntityB struct {
	Total     int `json:"total"`
	TotalPage int `json:"totalPage"`
	*EntityA
}

func (e *EntityA) SetMsg(s string) {
	e.Msg = s
}
func (e *EntityA) SetCode(i int) {
	e.Code = i
}
func (e *EntityB) SetMsg(s string) {
	e.SetMsg(s)
}
func (e *EntityB) SetCode(i int) {
	e.SetCode(i)
}
func (e *EntityB) SetTotal(Total int, TotalPage int) {
	e.Total = Total
	e.TotalPage = TotalPage
}
