package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"shoppinglist/modals"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/matoous/go-nanoid/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
)


type OrderController struct {
	client *mongo.Client 
}

func NewOrderController(c *mongo.Client) *OrderController{
	return &OrderController{c}
}

func (o OrderController) CreateOrder(w http.ResponseWriter , req *http.Request , _ httprouter.Params) {

	if req.Method == http.MethodPost{
		
		var orderlist modals.OrderList
		err := json.NewDecoder(req.Body).Decode(&orderlist) 
		if err != nil {
			http.Error(w,err.Error(),http.StatusInternalServerError) 
			return
		}

		var order modals.Order 
		order.OrderList = orderlist 
		order.OrderId = generateOrderID()
		loc, _ := time.LoadLocation("Africa/Cairo")
		now := time.Now().In(loc)
		order.Date = now

		collection := o.client.Database("shoppinglist").Collection("orders"); 
		_ , err = collection.InsertOne(context.Background() , order)
		if err != nil {
			http.Error(w,err.Error(),http.StatusInternalServerError)
			return
		}

		collection = o.client.Database("shoppinglist").Collection("orderlist")
		filter := bson.M{"username" : orderlist.Username}
		_, err = collection.DeleteOne(context.Background(),filter) 
		if err != nil{
			if err == mongo.ErrNoDocuments{
				http.Error(w,"No orderlist" , http.StatusNotFound)
				return
			}else{
				http.Error(w,err.Error(),http.StatusInternalServerError)
				return
			}
		}
	



	}else{
		http.Error(w,"Method not allowed" , http.StatusForbidden)
		return
	}

}

func (o OrderController) GetOrders(w http.ResponseWriter , req *http.Request , _ httprouter.Params){
	if req.Method == http.MethodGet{

		collection := o.client.Database("shoppinglist").Collection("orders")
		cursor, err := collection.Find(context.Background(), bson.D{}) // Bütün productları getir
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer cursor.Close(context.Background())

		var orders []modals.Order
		err = cursor.All(context.Background() , &orders)
		if err != nil{
			http.Error(w,err.Error(),http.StatusInternalServerError)
			return
		}
		jsonResp, err := json.Marshal(orders) 
		if err != nil {
			http.Error(w,err.Error(),http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResp)


	}else{
		http.Error(w,"Method Not Allowed" , http.StatusForbidden)
	}
} 

func generateOrderID() string {
	id, err := gonanoid.New(6)
	if err != nil {
		panic(err)
	}
	return id; 
}