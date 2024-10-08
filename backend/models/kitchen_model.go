package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type KitchenModel struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	KitchenName string             `bson:"kitchen_name" json:"kitchen_name"`
	Password    string             `bson:"password" json:"password"`
}
