package modals

type User struct {
	Username   string `json:"username" bson:"username"`
	Email      string `json:"email" bson:"email"`
	Password   string `json:"password" bson:"password"`
	Membership bool   `json:"membership" bson:"membership"`
}
type EditUser struct {
	Username   string `json:"username" bson:"username"`
	Email      string `json:"email" bson:"email"`
	Password   string `json:"password" bson:"password"`
	NewPassword   string `json:"newpassword" bson:"newpassword"`
}

