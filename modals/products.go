package modals

import "go.mongodb.org/mongo-driver/bson/primitive"

type Products struct {
	ID                 primitive.ObjectID `json:"_id,omitempty" bson:"_id"`
	Title              string             `json:"title" bson:"title"`
	Description        string             `json:"description" bson:"description"`
	Price              int                `json:"price" bson:"price"`
	DiscountPercentage float64            `json:"discountPercentage" bson:"discountPercentage"`
	Rating             float64            `json:"rating" bson:"rating"`
	Stock              int                `json:"stock" bson:"stock"`
	Brand              string             `json:"brand" bson:"brand"`
	Category           string             `json:"category" bson:"category"`
	Thumbnail          string             `json:"thumbnail" bson:"thumbnail"`
	Images             []string           `json:"images" bson:"images"`
	ProductID          int                `json:"productID" bson:"productID"`
	Quantity           int                `json:"quantity" bson:"quantity"`
}
