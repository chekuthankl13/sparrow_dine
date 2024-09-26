package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type StaffModel struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string             `json:"staff_name" bson:"staff_name"`
	PhoneNumber string             `json:"phone_number" bson:"phone_number"`
	Age         string             `json:"age" `
	Password    string             `json:"password"`
}
