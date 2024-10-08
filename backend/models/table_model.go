package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type TableModel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	TableName string             `bson:"table_name" json:"table_name"`
	Status    string             `bson:"status" json:"status"`
}
