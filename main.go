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
	"shoppinglist/controllers"
	"shoppinglist/modals"
)

func init() {
	gob.Register(modals.User{})
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("loading env variables", err)
	}
}

// Databaseden ürünleri getirme işlemi
// Database sipariş listesini ekleme. Kullanıcı ismine göresss
func main() {
	r := httprouter.New()
	pc := controllers.NewProductController(getClient())
	uc := controllers.NewUserController(getClient())
	olc := controllers.NewOrderListController(getClient())
	o := controllers.NewOrderController(getClient())
	admin := NewAdminController(getClient())

	corsHandler := func(h http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", getOrigin(r))
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

	r.POST("/order/createOrder", o.CreateOrder)
	r.GET("/order/getOrders", o.GetOrders) 

	r.GET("/admin/getAllUser/:username" , admin.getAllUserInformation)
	r.DELETE("/admin/deleteUserByAdmin" , admin.deleteUserAdmin)  

	 port := os.Getenv("PORT")
	 if port == "" {
	 	port = "8080"
	 }
	http.ListenAndServe(":"+port, corsHandler(r))
}

func getOrigin(req *http.Request) string {
	origin := req.Header.Get("Origin")
	if origin == "https://mf-shoppinglist.vercel.app" {
		return "https://mf-shoppinglist.vercel.app"
	}
	return "http://localhost:3000"

	// Bir tane origin kabul ediliyor. Veya hepsi. Ama hepsi güvenlik açısından riskli
	// Bu yüzden bu fonksiyon ile sadece belirli origin kabul ediliyor.
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
