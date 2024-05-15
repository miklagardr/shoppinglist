package controllers

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"log"

	"net/http"

	"shoppinglist/modals"

	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var store *sessions.CookieStore

func init() {

	store = sessions.NewCookieStore(generateSessionKey())

	store.Options = &sessions.Options{
		Path:     "/",                   
		MaxAge:   31536000,                 
		SameSite: http.SameSiteNoneMode, 
		Secure:   true,                    
	}
}

type UserController struct {
	client *mongo.Client
}

func NewUserController(c *mongo.Client) *UserController {
	return &UserController{c}
}

func generateSessionKey() []byte {
	key := make([]byte, 32) // 32 byte'lık rastgele bir anahtar oluşturun
	_, err := rand.Read(key)
	if err != nil {
		log.Fatal("Error generating session key: ", err)
	}
	return key
}

func (uc UserController) GetUser(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	session, err := store.Get(req, "user-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if authenticated, ok := session.Values["authenticated"].(bool); ok && authenticated {
		user, ok := session.Values["user"].(modals.User) // Cookie'deki user bilgilerini alıyoruz.
		if !ok {
			http.Error(w, "Invalid user data in the session", http.StatusInternalServerError)
			return
		}
		jsonUser, _ := json.Marshal(user)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonUser)

	} else {
		http.Error(w, "User is not logged in", http.StatusUnauthorized)
		return
	}
}

func (uc UserController) LogInUser(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	if req.Method == http.MethodPost {

		session, err := store.Get(req, "user-session")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if authenticated, ok := session.Values["authenticated"].(bool); ok && authenticated {
			http.Error(w, "User is  logged in", http.StatusUnauthorized)
			return
		}

		var user modals.User
		err = json.NewDecoder(req.Body).Decode(&user)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var existingUser modals.User
		collection := uc.client.Database("shoppinglist").Collection("users")
		err = collection.FindOne(context.TODO(), bson.M{"username": user.Username}).Decode(&existingUser)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				http.Error(w, "There is no such a user", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
		if user.Password != existingUser.Password {
			http.Error(w, "Invalid Password", http.StatusUnauthorized)
			return
		}

		session.Values["user"] = modals.User{
			Username:   existingUser.Username,
			Email:      existingUser.Email,
			Membership: existingUser.Membership,
		}
		session.Values["authenticated"] = true

		err = session.Save(req, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonResponse := modals.User{
			Username:   existingUser.Username,
			Email:      existingUser.Email,
			Membership: existingUser.Membership,
		}
		// Kullanıcı bilgilerini bastırdık. Ama password veya id bastırılmadı. Güvenlik sebeplerinden ötürü.

		w.Header().Set("Content-Type", "application/json")
		jsonUser, _ := json.Marshal(jsonResponse)
		w.Write(jsonUser) // Kullanıcı bilgilerini bastır.

	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}

}

func (uc UserController) LogOutUser(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	if req.Method == http.MethodPost {

		session, err := store.Get(req, "user-session")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Check if the user is authenticated before attempting to log out
		if authenticated, ok := session.Values["authenticated"].(bool); ok && authenticated {
			// Clear the sessions data
			for key := range session.Values {
				delete(session.Values, key)
			}

			session.Values["authenticated"] = false
			err = session.Save(req, w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("User logged out successfully"))
		} else {
			http.Error(w, "User is not logged in", http.StatusUnauthorized)
			return
		}
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
}

// Bir veriye daha önce bir şey atadıysan TODO , Atamadıysan context.Background() kullanılır.

func (uc UserController) CreateUser(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	if req.Method == http.MethodPost {
		var user modals.User

		err := json.NewDecoder(req.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		collection := uc.client.Database("shoppinglist").Collection("users")
		var existingUser modals.User
		err = collection.FindOne(context.TODO(), bson.M{"$or": []bson.M{
			{"username": user.Username},
			{"email": user.Email},
		}}).Decode(&existingUser)
		// Kullanıcı adını kontrol ettik. Eğer hata vermezse kullanıcı var demek.
		if err == nil {
			http.Error(w, "Username or Email already exists", http.StatusConflict)
			return
		} else if err != mongo.ErrNoDocuments {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = collection.InsertOne(context.Background(), user) // Result alıcaz. Ama almıyoruz çünkü user'ı direk bastırıyoruz.
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("User created successfully.Redirect to login page..")) // User bilgilerini geri döndermek riskli.

	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (uc UserController) DeleteUser(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	if req.Method == http.MethodDelete {

		var user modals.User
		err := json.NewDecoder(req.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		session, err := store.Get(req, "user-session")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for key := range session.Values {
			delete(session.Values, key)
		}

		session.Values["authenticated"] = false
		err = session.Save(req, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		collection := uc.client.Database("shoppinglist").Collection("users")
		_, err = collection.DeleteOne(context.Background(), bson.M{"username": user.Username})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte("User deleted successfully"))

	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (uc UserController) UpdateEmail(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	if req.Method == http.MethodPut {
		var user modals.EditUser
		err := json.NewDecoder(req.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		session, err := store.Get(req, "user-session")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if authenticated, ok := session.Values["authenticated"].(bool); ok && authenticated {
			existingUser, ok := session.Values["user"].(modals.User)
			if !ok {
				http.Error(w, "Invalid user data in the session", http.StatusInternalServerError)
				return
			}

			collection := uc.client.Database("shoppinglist").Collection("users")
			if existingUser.Email == user.Email {
				w.Header().Set("Content-Type", "application/json")
				jsonUser, _ := json.Marshal("Email is already same")
				w.Write(jsonUser)
				return
			}
			_, err = collection.UpdateOne(context.Background(), bson.M{"username": existingUser.Username}, bson.M{"$set": bson.M{"email": user.Email}})

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			existingUser.Email = user.Email
			session.Values["user"] = existingUser
			err = session.Save(req, w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			jsonUser, _ := json.Marshal("Email updated successfully")
			w.Write(jsonUser)
		}

	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (uc UserController) UpdatePassword(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	// Email'de bilgileri session'dan aldık. Fakat password için session'dan alamayız. Çünkü password session'da tutulmaz.
	// Bu yüzden body'den alıcaz.
	// Aynı zamanda body içerisinden username'i de alıcaz. Çünkü password değiştirirken username'e ihtiyacımız var.
	// Hemde eski password'u da alıcaz. Çünkü eski password'u kontrol etmemiz gerekiyor.

	if req.Method == http.MethodPut {
		var user modals.EditUser
		err := json.NewDecoder(req.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		collection := uc.client.Database("shoppinglist").Collection("users")
		var existingUser modals.User
		err = collection.FindOne(context.TODO(), bson.M{"username": user.Username}).Decode(&existingUser)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				http.Error(w, "There is no such a user", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
		if user.Password != existingUser.Password {
			http.Error(w, "Invalid old Password", http.StatusUnauthorized)
			return
		} else {
			_, err = collection.UpdateOne(context.Background(), bson.M{"username": user.Username}, bson.M{"$set": bson.M{"password": user.NewPassword}})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			jsonUser, _ := json.Marshal("Password updated successfully")
			w.Write(jsonUser)
		}

	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

}
