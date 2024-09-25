package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type AdminCredModel struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	UserName string             `json:"user_name" bson:"user_name"`
	Password string             `json:"password" bson:"password"`
}
