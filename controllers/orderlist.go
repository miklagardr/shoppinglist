package controllers

import (
	"context"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"shoppinglist/modals"
)

type OrderListController struct {
	client *mongo.Client
}

func NewOrderListController(c *mongo.Client) *OrderListController {
	return &OrderListController{c}
}
func (olc OrderListController) CreateOrderList(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	if req.Method == http.MethodPost {
		var orderList modals.OrderList
		err := json.NewDecoder(req.Body).Decode(&orderList)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		collection := olc.client.Database("shoppinglist").Collection("orderlist")
		_, err = collection.InsertOne(context.Background(), orderList)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonOrderList, err := json.Marshal(orderList)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonOrderList)

	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
}
func (olc OrderListController) AddProductOrderList(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	if req.Method == http.MethodPut {
		var newOrderList modals.SingleOrderList // Bize gelen değer. İçinde hangi kullanıcının hangi ürünü eklediği ve ürünün fiyatı yazıyor.
		err := json.NewDecoder(req.Body).Decode(&newOrderList)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		collection := olc.client.Database("shoppinglist").Collection("orderlist")
		var existingModel modals.OrderList
		filter := bson.M{"username": newOrderList.Username}
		// filterProduct := bson.M{"productID": newOrderList.Products.ProductID}

		err = collection.FindOne(context.TODO(), filter).Decode(&existingModel)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				http.Error(w, "No order list or login", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}

		for i, product := range existingModel.Products {
			if product.ProductID == newOrderList.Products.ProductID {
				existingModel.OrderTotalPrice += newOrderList.Products.Price
				existingModel.Products[i].Quantity += 1
				update := bson.M{"$set": bson.M{
					"products":        existingModel.Products,
					"ordertotalprice": existingModel.OrderTotalPrice,
				}}
				_, err = collection.UpdateOne(context.Background(), filter, update)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				jsonResult, err := json.Marshal(existingModel)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write(jsonResult)
				return
			}
		}

		existingModel.OrderTotalPrice += newOrderList.ProductsPrice                    // Kullanıcının gönderdiği ürünün fiyatını existingModel'in fiyatına ekleme
		existingModel.Products = append(existingModel.Products, newOrderList.Products) // Kullanıcının gönderdiği ürünü existingModel'e ekliyoruz.

		update := bson.M{"$set": bson.M{
			"products":        existingModel.Products,
			"ordertotalprice": existingModel.OrderTotalPrice,
		}}

		_, err = collection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		jsonResult, err := json.Marshal(existingModel)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResult)

	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (olc OrderListController) DeleteOrderList(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	if req.Method == http.MethodDelete {
		var newOrderList modals.SingleOrderList // Burada sadece silinecek ürün bilgisini alıyoruz
		err := json.NewDecoder(req.Body).Decode(&newOrderList)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		collection := olc.client.Database("shoppinglist").Collection("orderlist")
		filter := bson.M{"username": newOrderList.Username}

		var existingModel modals.OrderList
		err = collection.FindOne(context.Background(), filter).Decode(&existingModel)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				http.Error(w, "No order list or login", http.StatusNotFound)
				return
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		// Silinecek ürünü order listesinden çıkar
		for i, product := range existingModel.Products {
			if product.ProductID == newOrderList.Products.ProductID {
				if existingModel.Products[i].Quantity == 1 {
					existingModel.OrderTotalPrice -= newOrderList.Products.Price
					existingModel.Products = append(existingModel.Products[:i], existingModel.Products[i+1:]...)
					// 0'dan i'ye kadar olan ürünleri al, i'den sonrasını al ve birleştir.
					// i dahil olmuyor ve silinmiş oluyor.
					break
				} else {
					existingModel.OrderTotalPrice -= newOrderList.Products.Price
					existingModel.Products[i].Quantity -= 1
					break
				}

			}
		}

		// Eğer ürün kalmadıysa tamamı silinsin
		if len(existingModel.Products) == 0 {
			_, err = collection.DeleteOne(context.Background(), filter)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			// Sipariş listesini güncelle
			update := bson.M{"$set": bson.M{
				"products":        existingModel.Products,
				"ordertotalprice": existingModel.OrderTotalPrice,
			}}

			_, err := collection.UpdateOne(context.Background(), filter, update)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		jsonResp, _ := json.Marshal(existingModel)

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResp)
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (olc OrderListController) GetOrderList(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	if req.Method == http.MethodGet {
		username := p.ByName("username")
		filter := bson.M{"username": username}
		var orderlist modals.OrderList
		collection := olc.client.Database("shoppinglist").Collection("orderlist")
		err := collection.FindOne(context.Background(), filter).Decode(&orderlist)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				http.Error(w, "No order list or login", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
		jsonOrder, _ := json.Marshal(orderlist)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonOrder)

	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
}
