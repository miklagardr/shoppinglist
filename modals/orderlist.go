package modals

type OrderList struct {
	Username        string     `json:"username" bson:"username"`
	Products        []Products `json:"products" bson:"products"`
	OrderTotalPrice int        `json:"ordertotalprice" bson:"ordertotalprice"`
}
