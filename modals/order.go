package modals

import "time"

type Order struct{
	OrderId string `json:"orderid" bson:"orderid"` 
	OrderList OrderList `json:"orderlist" bson:"orderlist"` 
	Date time.Time `json:"date" bson:"date"`
}