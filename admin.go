package main

import (
	"context"
	"encoding/json"
	"net/http"
	"shoppinglist/modals"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Admin struct {
	client *mongo.Client
}

func NewAdminController(c *mongo.Client) *Admin {
	return &Admin{c}
}

func (a Admin) allUserInfo() ([]modals.User, error) {

	collection := a.client.Database("shoppinglist").Collection("users")

	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		// Hata durumunda nil kullanıcı dizisi ve hatayı döndür
		return nil, err
	}
	defer cursor.Close(context.Background())

	var users []modals.User

	for cursor.Next(context.Background()) {
		var user modals.User
		if err := cursor.Decode(&user); err != nil {

			return nil, err
		}
		users = append(users, user)
	}
	if err := cursor.Err(); err != nil {

		return nil, err
	}

	return users, nil
}

func (a Admin) getAllUserInformation(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	if req.Method == http.MethodGet {
		username := p.ByName("username")

		if username != "Admin" {
			http.Error(w, "You are not authorized", http.StatusForbidden)
			return
		}

		users, err := a.allUserInfo()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		jsonUsers, err := json.Marshal(users)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonUsers)
	} else {
		http.Error(w, "Method Not Allowed", http.StatusForbidden)
		return
	}
}

func (a Admin) deleteUserAdmin(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	if req.Method == http.MethodDelete {
		// admin username , user username for delete.
		var admin modals.AdminDeleteUser
		err := json.NewDecoder(req.Body).Decode(&admin)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if admin.Username != "Admin" {
			http.Error(w, "You are not authorized", http.StatusForbidden)
			return
		}
		collection := a.client.Database("shoppinglist").Collection("users")
		filter := bson.M{"username": admin.Userusername}

		_, err = collection.DeleteOne(context.Background(), filter)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				http.Error(w, "User could not found", http.StatusNotFound)
				return
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		collection = a.client.Database("shoppinglist").Collection("orderlist")
		_, err = collection.DeleteOne(context.Background(), filter)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				http.Error(w, "User do not have an order list", http.StatusNotFound)
				return
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		users, err := a.allUserInfo()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonUser, err := json.Marshal(users)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonUser)
	} else {
		http.Error(w, "Method Not Allowed", http.StatusForbidden)
	}

}
