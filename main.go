package main

import (
	"context"
	"encoding/gob"
	"os"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"

	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	//"os"
	"shoppinglist/controllers"
	"shoppinglist/modals"
)

func init() {
	gob.Register(modals.User{})
}

// Databaseden ürünleri getirme işlemi
// Database sipariş listesini ekleme. Kullanıcı ismine göresss
func main() {
	r := httprouter.New()
	pc := controllers.NewProductController(getClient())
	uc := controllers.NewUserController(getClient())
	olc := controllers.NewOrderListController(getClient())

	

	corsHandler := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			h.ServeHTTP(w, r)
		})
	}

	// Uniqe username'e göre kullanıcı oluşturma ve getirme işlemi

	r.POST("/user/login", uc.LogInUser)
	r.POST("/user/createUser", uc.CreateUser)
	r.POST("/user/logout", uc.LogOutUser)
	r.GET("/user/getUser", uc.GetUser)
	r.PUT("/user/editEmail", uc.UpdateEmail)
	r.PUT("/user/editPassword", uc.UpdatePassword)
	r.DELETE("/user/deleteUser", uc.DeleteUser)

	r.GET("/products", pc.GetAllProduct)  // Bütün productsları getir
	r.GET("/products/:id", pc.GetProduct) // Id'ye göre product getir

	r.GET("/category/:category", pc.GetCategoryProduct)

	r.POST("/orderlist/create", olc.CreateOrderList)        // Kullanıcı id'sine göre order list ekle
	r.PUT("/orderlist/addproduct", olc.AddProductOrderList) // Id'ye göre order list sil
	r.DELETE("/orderlist/deleteproduct", olc.DeleteOrderList)
	r.GET("/orderlist/fetch/:username", olc.GetOrderList)

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("loading env variables", err)
	}

	 port := os.Getenv("PORT")
	 if port == "" {
	 	port = "8080"
	 }
	http.ListenAndServe(":"+port, corsHandler(r))
}
func getClient() *mongo.Client {
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGODB_URI"))  // Bağlantı adresi
	client, err := mongo.Connect(context.Background(), clientOptions) // Bağlantı kurma
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.Background(), nil) // Bağlantıyı kontrol etme
	if err != nil {
		log.Fatal(err)
	}
	return client
}
