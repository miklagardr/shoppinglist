package controllers

import (
	"context"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"shoppinglist/modals"
	"strconv"
)

type ProductController struct {
	client *mongo.Client
}

func NewProductController(c *mongo.Client) *ProductController {
	return &ProductController{c}
}

func (pc ProductController) GetAllProduct(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Bütün productsları getir

	if req.Method == http.MethodGet {
		collection := pc.client.Database("shoppinglist").Collection("products")
		cursor, err := collection.Find(context.Background(), bson.D{}) 
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		defer cursor.Close(context.Background())

		var products []modals.Products
		err = cursor.All(context.Background(), &products)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		jsonProducts, err := json.Marshal(products)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonProducts)
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (pc ProductController) GetProduct(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	// id'ye product getir
	if req.Method == http.MethodGet {
		productID := p.ByName("id")
		intProductID, _ := strconv.Atoi(productID)
		collection := pc.client.Database("shoppinglist").Collection("products")
		var product modals.Products
		err := collection.FindOne(context.Background(), bson.M{"productID": intProductID}).Decode(&product)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				http.Error(w, "No Products", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
		singleProduct, err := json.Marshal(product)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Write(singleProduct)
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (pc ProductController) GetCategoryProduct(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	if req.Method == http.MethodGet {
		category := p.ByName("category")
		var products []modals.Products

		collection := pc.client.Database("shoppinglist").Collection("products")
		cursor, err := collection.Find(context.Background(), bson.M{"category": category})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer cursor.Close(context.Background())

		err = cursor.All(context.Background(), &products)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonProducts, err := json.Marshal(products)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonProducts)

	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
