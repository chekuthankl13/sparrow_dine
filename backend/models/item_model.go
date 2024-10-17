package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ItemModel struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	ImageUrl string             `bson:"image_url" json:"image_url"`
	ItemName string             `bson:"item_name" json:"item_name"`
	Price    string             `bson:"price" json:"price"`
	Qty      string             `bson:"qty" json:"qty"`
	SubQty   []SubItem          `bson:"sub_qty" json:"sub_qty"`
}

type SubItem struct {
	Qty   string `json:"qty" form:"qty"`
	Price string `json:"price" form:"price"`
}
