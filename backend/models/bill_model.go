package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type BillModel struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	CustomerNumber    string             `bson:"customer_number" json:"customer_number"`
	CustomerName      string             `bson:"customer_name" json:"customer_name"`
	Date              string             `bson:"date" json:"date"`
	BilledTime        string             `bson:"billed_time" json:"billed_time"`
	PaidTime          string             `bson:"paid_time" json:"paid_time"`
	PaymentStatus     bool               `bson:"payment_status" json:"payment_status"`
	PaymentType       string             `bson:"payment_type" json:"payment_type"`
	Discount          string             `bson:"discount" json:"discount"`
	TotalItemAmount   string             `bson:"total_item_amount" json:"total_item_amount"`
	TotalParcelAmount string             `bson:"total_parcel_amount" json:"total_parcel_amount"`
	NetTotal          string             `bson:"net_total" json:"net_total"`
	Type              string             `bson:"type" json:"type"`
	Items             []BillItemModel    `bson:"items" json:"items"`
	Parcels           []BillItemModel    `bson:"parcels" json:"parcels"`
}

type BillItemModel struct {
	ItemId     string `bson:"item_id" json:"item_id"`
	ItemName   string `bson:"item_name" json:"item_name"`
	ItemQty    string `bson:"item_qty" json:"item_qty"`
	Qty        string `bson:"qty" json:"qty"`
	Price      string `bson:"price" json:"price"`
	TotalPrice string `bson:"total_price" json:"total_price"`
}
