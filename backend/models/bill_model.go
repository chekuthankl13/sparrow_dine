package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type BillModel struct {
	ID            primitive.ObjectID `bson:"_id,omitemity" json:"id,omitempty"`
	Date          string             `bson:"date" json:"date"`
	BilledTime    string             `bson:"billed_time" json:"billed_time"`
	PayedTime     string             `bson:"payed_time" json:"payed_time"`
	PaymentStatus bool               `bson:"payment_status" json:"payment_status"`
	PaymentType   string             `bson:"payment_type" json:"payment_type"`
	Discount      string             `bson:"discount" json:"discount"`
	TotalAmout    string             `bson:"total_amount" json:"total_amount"`
	NetTotal      string             `bson:"net_total" json:"net_total"`
	Type          string             `bson:"type" json:"type"`
	Items         []BillItemModel    `bson:"items" json:"items"`
}

type BillItemModel struct {
	ItemId     string `bson:"item_id" json:"item_id"`
	ItemQty    string `bson:"item_qty" json:"item_qty"`
	Qty        string `bson:"qty" json:"qty"`
	Price      string `bson:"price" json:"price"`
	TotalPrice string `bson:"total_price" json:"total_price"`
}
