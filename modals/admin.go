package modals

type Admin struct {
	Username   string `json:"username" bson:"username"`
}
type AdminDeleteUser struct{
	Username string `json:"username" bson:"username"`
	Userusername string `json:"userusername" bson:"userusername"` 
}