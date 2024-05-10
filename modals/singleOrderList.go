package modals

type SingleOrderList struct {
	Username      string   `json:"username"`
	Products      Products `json:"products"`
	ProductsPrice int      `json:"productsprice"`
}
