package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/chekuthankl13/sparrow_dine/helpers"
	"github.com/chekuthankl13/sparrow_dine/middlewares"
	"github.com/chekuthankl13/sparrow_dine/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/////////

type adminCred struct {
	Username string `form:"user_name" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type staff struct {
	StaffName   string `form:"staff_name" binding:"required"`
	PhoneNumber string `form:"phone_number" binding:"required"`
	Age         string `form:"age" binding:"required"`
	Password    string `form:"password" binding:"required"`
}

type staffUpdate struct {
	StaffName   string `form:"staff_name,omitempty"`
	PhoneNumber string `form:"phone_number,omitempty"`
	Age         string `form:"age,omitempty"`
	Password    string `form:"password,omitempty" `
}

type kitchen struct {
	KitchenName string `form:"kitchen_name" binding:"required"`
	Password    string `form:"password" binding:"required"`
}

type table struct {
	TableName string `form:"table_name" binding:"required"`
	Status    string `form:"status"`
}

type item struct {
	ItemName string `form:"item_name" binding:"required"`
	Price    string `form:"price" binding:"required"`
	Qty      string `form:"qty" binding:"required"`
	SubQty   string `form:"sub_qty"`
	Addon    string `form:"addons"`
}

type itemUpdate struct {
	ItemName string `form:"item_name,omitempty"`
	Price    string `form:"price,omitempty"`
	Qty      string `form:"qty,omitempty"`
	SubQty   string `form:"sub_qty,omitempty"`
	Addons   string `form:"addons,omitempty"`
}

type Bill struct {
	CustomerNumber    string `form:"customer_number" binding:"required"`
	CustomerName      string `form:"customer_name" binding:"required"`
	Date              string `form:"date" binding:"required"`
	BilledTime        string `form:"billed_time" binding:"required"`
	PaidTime          string `form:"paid_time"`
	PaymentStatus     bool   `form:"payment_status"`
	PaymentType       string `form:"payment_type"`
	Discount          string `form:"discount"`
	TotalItemAmount   string `form:"total_item_amount" binding:"required"`
	TotalParcelAmount string `form:"total_parcel_amount"`
	NetTotal          string `form:"net_total" binding:"required"`
	Type              string `form:"type" binding:"required"`
	Items             string `form:"items"`
	Parcels           string `form:"parcels"`
}

/////////////

func AdminLogin(c *gin.Context) {
	key := os.Getenv("SECRET_KEY")
	collection := helpers.DB.Collection("admin_cred")

	var inputCred adminCred
	if err := c.Bind(&inputCred); err != nil {
		helpers.BadResponse(c, err.Error())
		return
	}

	var result models.AdminCredModel

	if err := collection.FindOne(context.Background(), bson.M{"user_name": inputCred.Username}).Decode(&result); err != nil {
		helpers.BadResponse(c, "invalid credtional")
		return
	}

	isVerify := middlewares.CheckHashPsw(inputCred.Password, result.Password)

	if !isVerify {
		helpers.BadResponse(c, "password does not match !!")
		return
	}

	//// jwt generation

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{

		/*
			"data": map[string]string{"userId": "5f48d1a5f29cea66f634f2ec",
						"userType":  "default",
						"userTable": "customerTable"},
					"exp": time.Now().AddDate(1, 0, 0).Unix(),
					"iat": time.Now().Unix(),
		*/
		"id":  result.ID.Hex(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := claims.SignedString([]byte(key))
	if err != nil {
		helpers.BadResponse(c, err.Error())
		return
	}
	/////////
	loginres := map[string]string{
		"user_id":  result.ID.Hex(),
		"username": result.UserName, "token": token, "password": inputCred.Password}
	helpers.SuccessResponse(c, "login success !", loginres)

}

func AdminRegister(c *gin.Context) {
	collection := helpers.DB.Collection("admin_cred")
	fmt.Println("*****1*******")
	var inputCred adminCred
	if err := c.Bind(&inputCred); err != nil {
		helpers.BadResponse(c, err.Error())
		return
	}
	fmt.Println("*****2*******")

	count, _ := collection.CountDocuments(context.TODO(), bson.M{"user_name": inputCred.Username})
	fmt.Println("count -", count)
	if count >= 1 {
		helpers.BadResponse(c, "username already exist !!")
		return
	}
	fmt.Println("*****3*******")

	hashPsw, err := middlewares.HashPassword(inputCred.Password)
	if err != nil {
		helpers.BadResponse(c, err.Error())
		return
	}
	fmt.Println("*****4*******")

	cred := models.AdminCredModel{UserName: inputCred.Username, Password: hashPsw}

	result, err := collection.InsertOne(context.TODO(), &cred)
	if err != nil {
		helpers.BadResponse(c, err.Error())
		return
	}
	fmt.Println("*****5*******")

	helpers.SuccessResponse(c, "admin registered successfully", result.InsertedID)

}

func StaffCreate(c *gin.Context) {
	collection := helpers.DB.Collection("staff")

	var input staff

	if err := c.Bind(&input); err != nil {
		helpers.BadResponse(c, err.Error())
		return
	}

	count, err := collection.CountDocuments(context.Background(), bson.M{"phone_number": input.PhoneNumber})
	if err != nil {
		helpers.BadResponse(c, err.Error())
		return
	}

	if count >= 1 {
		helpers.BadResponse(c, "staff with the phone number already exist !!")
		return
	}

	hashPsw, err := middlewares.HashPassword(input.Password)
	if err != nil {
		helpers.BadResponse(c, err.Error())
		return
	}

	staff := models.StaffModel{Name: input.StaffName, PhoneNumber: input.PhoneNumber, Age: input.Age, Password: string(hashPsw)}

	result, err := collection.InsertOne(context.Background(), &staff)
	if err != nil {
		helpers.BadResponse(c, err.Error())
		return
	}

	helpers.SuccessResponse(c, "staff successfully created !", result.InsertedID)
}

func GetStaffs(c *gin.Context) {

	colllection := helpers.DB.Collection("staff")

	var data []models.StaffModel

	cur, err := colllection.Find(context.Background(), bson.M{})
	if err != nil {
		helpers.BadResponse(c, err.Error())
		return
	}

	for cur.Next(context.TODO()) {
		var staff models.StaffModel
		err := cur.Decode(&staff)
		if err != nil {
			helpers.BadResponse(c, err.Error())
			return
		}
		data = append(data, staff)
	}

	defer cur.Close(context.TODO())

	if data == nil {
		data = []models.StaffModel{}
	}

	helpers.SuccessResponse(c, "staffs", data)

}

func DeleteStaff(c *gin.Context) {
	collection := helpers.DB.Collection("staff")
	param := c.Param("id")
	id, err := primitive.ObjectIDFromHex(param)

	if err != nil {
		helpers.BadResponse(c, err.Error())
		return
	}

	count, _ := collection.CountDocuments(context.Background(), bson.M{"_id": id})

	if count == 0 {
		helpers.BadResponse(c, "Staff with the id not exist !!")
		return
	}

	result, err := collection.DeleteOne(context.Background(), bson.D{{Key: "_id", Value: id}})

	if err != nil {
		helpers.BadResponse(c, err.Error())
		return
	}

	helpers.SuccessResponse(c, "staff successfully deleted !!", result.DeletedCount)

}

func EditStaff(c *gin.Context) {
	collection := helpers.DB.Collection("staff")
	param := c.Param("id")
	id, err := primitive.ObjectIDFromHex(param)
	var input staffUpdate

	if err != nil {
		helpers.BadResponse(c, err.Error())
		return
	}

	count, _ := collection.CountDocuments(context.Background(), bson.M{"_id": id})
	if count == 0 {
		helpers.BadResponse(c, "Staff with the id not exist !!")
		return
	}

	if err := c.Bind(&input); err != nil {
		helpers.BadResponse(c, err.Error())
		return
	}
	fmt.Println("*******1******")
	var updatedField bson.D

	if input.Age != "" {
		updatedField = append(updatedField, bson.E{Key: "age", Value: input.Age})
	}
	if input.PhoneNumber != "" {
		updatedField = append(updatedField, bson.E{Key: "phone_number", Value: input.PhoneNumber})
	}

	if input.StaffName != "" {
		updatedField = append(updatedField, bson.E{Key: "staff_name", Value: input.StaffName})
	}

	if input.Password != "" {
		hashPsw, err := middlewares.HashPassword(input.Password)
		if err != nil {
			helpers.BadResponse(c, err.Error())
			return
		}
		updatedField = append(updatedField, bson.E{Key: "password", Value: hashPsw})

	}
	fmt.Println("*******2******")

	update := bson.D{{Key: "$set", Value: updatedField}}
	filter := bson.D{{Key: "_id", Value: id}}
	fmt.Println("*******3******")
	fmt.Println(updatedField)
	fmt.Println("*******3.5******")

	res, err := collection.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		helpers.BadResponse(c, err.Error())
		return
	}
	helpers.SuccessResponse(c, "staff updated successfully !", res.UpsertedID)

}

func CreateKitchen(c *gin.Context) {
	collection := helpers.DB.Collection("kitchen")
	var input kitchen
	if err := c.Bind(&input); err != nil {
		helpers.BadResponse(c, err.Error())
		return
	}

	hashPsw, err := middlewares.HashPassword(input.Password)
	if err != nil {
		helpers.BadResponse(c, err.Error())
		return
	}

	kitchen := models.KitchenModel{KitchenName: input.KitchenName, Password: hashPsw}

	res, err := collection.InsertOne(context.TODO(), &kitchen)
	if err != nil {
		helpers.BadResponse(c, err.Error())
		return
	}

	helpers.SuccessResponse(c, "kitchen created successfully", res.InsertedID)
}

func GetKitchen(c *gin.Context) {
	collection := helpers.DB.Collection("kitchen")

	var data []models.KitchenModel

	cur, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		helpers.BadResponse(c, err.Error())
		return
	}
	for cur.Next(context.TODO()) {
		var kitchen models.KitchenModel
		err := cur.Decode(&kitchen)
		if err != nil {
			helpers.BadResponse(c, err.Error())
			return
		}
		data = append(data, kitchen)
	}
	defer cur.Close(context.TODO())

	if data == nil {
		data = []models.KitchenModel{}
	}

	helpers.SuccessResponse(c, "kitchen list", data)
}

func DeleteKitchen(c *gin.Context) {
	collection := helpers.DB.Collection("kitchen")
	params := c.Param("id")

	id, err := primitive.ObjectIDFromHex(params)
	if err != nil {
		helpers.BadResponse(c, err.Error())
		return
	}

	count, _ := collection.CountDocuments(context.Background(), bson.M{"_id": id})
	if count == 0 {
		helpers.BadResponse(c, "kitchen not found")
		return
	}

	res, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		helpers.BadResponse(c, err.Error())
		return
	}

	helpers.SuccessResponse(c, "kitchen successfully deleted !!", res.DeletedCount)
}

func CreateTable(c *gin.Context) {
	collection := helpers.DB.Collection("table")
	var input table
	if err := c.Bind(&input); err != nil {
		helpers.BadResponse(c, err.Error())
		return
	}

	input.Status = "available"

	table := models.TableModel{TableName: input.TableName, Status: input.Status}

	res, err := collection.InsertOne(context.TODO(), &table)
	if err != nil {
		helpers.BadResponse(c, err.Error())
		return
	}
	helpers.SuccessResponse(c, "table created successfully", res.InsertedID)

}

func GetTable(c *gin.Context) {
	collection := helpers.DB.Collection("table")

	var data []models.TableModel

	cur, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		helpers.BadResponse(c, err.Error())
		return
	}

	for cur.Next(context.TODO()) {
		var table models.TableModel
		err := cur.Decode(&table)
		if err != nil {
			helpers.BadResponse(c, err.Error())
			return
		}

		data = append(data, table)
	}
	defer cur.Close(context.TODO())

	if data == nil {
		data = []models.TableModel{}
	}
	helpers.SuccessResponse(c, "table list", data)
}

func DeleteTable(c *gin.Context) {
	collection := helpers.DB.Collection("table")
	params := c.Param("id")

	id, err := primitive.ObjectIDFromHex(params)
	if err != nil {
		helpers.BadResponse(c, err.Error())
		return
	}

	count, _ := collection.CountDocuments(context.Background(), bson.M{"_id": id})

	if count == 0 {
		helpers.BadResponse(c, "table not found")
		return
	}

	res, err := collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		helpers.BadResponse(c, err.Error())
		return
	}

	helpers.SuccessResponse(c, "table seleted successfully", res.DeletedCount)
}

func CreateItem(c *gin.Context) {
	collection := helpers.DB.Collection("item")
	var input item
	fmt.Println("*********1********")

	if err := c.Bind(&input); err != nil {
		helpers.BadResponse(c, err.Error())
		return
	}
	fmt.Println("*********2********")

	file, err := c.FormFile("image")
	if err != nil {
		helpers.BadResponse(c, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	storage, err := helpers.FirebaseApp.Storage(context.Background())
	if err != nil {
		helpers.BadResponse(c, err.Error())
		return
	}

	bucket, err := storage.Bucket(os.Getenv("BUCKET_NAME"))
	if err != nil {
		helpers.BadResponse(c, err.Error())
		return
	}

	f, err := file.Open()
	if err != nil {
		helpers.BadResponse(c, err.Error())
		return
	}
	defer f.Close()

	obj := bucket.Object(os.Getenv("BUCKET_FOLDER") + "/" + file.Filename)

	writer := obj.NewWriter(ctx)
	id := uuid.New().String()
	writer.ObjectAttrs.Metadata = map[string]string{"firebaseStorageDownloadTokens": id}

	if _, err := io.Copy(writer, f); err != nil {
		helpers.BadResponse(c, err.Error())
		return
	}

	if err := writer.Close(); err != nil {
		helpers.BadResponse(c, err.Error())
		return
	}

	imageUrl := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/%s/o/%s%s%s?alt=media&token=%s", os.Getenv("BUCKET_NAME"), os.Getenv("BUCKET_FOLDER"), "%2f", file.Filename, id)

	var j []models.SubItem

	if input.SubQty != "" {
		err = json.Unmarshal([]byte(input.SubQty), &j)
		if err != nil {
			helpers.BadResponse(c, err.Error())
			return
		}
	}

	var addons []string

	if input.Addon != "" {
		if err := json.Unmarshal([]byte(input.SubQty), &addons); err != nil {
			helpers.BadResponse(c, err.Error())
			return
		}
	}

	var data = models.ItemModel{ImageUrl: imageUrl, Price: input.Price, ItemName: input.ItemName, Qty: input.Qty, SubQty: j, Addons: addons}
	res, err := collection.InsertOne(context.Background(), &data)
	if err != nil {
		helpers.BadResponse(c, err.Error())
		return
	}

	helpers.SuccessResponse(c, "item created successfully !!", res.InsertedID)

}

func GetItems(c *gin.Context) {
	collection := helpers.DB.Collection("item")

	var data []models.ItemModel

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		helpers.BadResponse(c, err.Error())
		return
	}

	for cursor.Next(context.Background()) {
		var item models.ItemModel
		err := cursor.Decode(&item)
		if err != nil {
			helpers.BadResponse(c, err.Error())
			return
		}

		data = append(data, item)
	}

	helpers.SuccessResponse(c, "item list", data)

}

func DeleteItem(c *gin.Context) {
	collection := helpers.DB.Collection("item")
	params := c.Param("id")
	id, err := primitive.ObjectIDFromHex(params)
	if err != nil {
		helpers.BadResponse(c, err.Error())
		return
	}

	count, _ := collection.CountDocuments(context.TODO(), bson.M{"_id": id})

	if count == 0 {
		helpers.BadResponse(c, "document not found !!")
		return
	}

	var res1 models.ItemModel

	if err := collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&res1); err != nil {
		helpers.BadResponse(c, err.Error())
		return
	}

	base := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/%s/o/%s%s", os.Getenv("BUCKET_NAME"), os.Getenv("BUCKET_FOLDER"), "%2f")

	// startUrl := len(base)
	// endUrl := len(res1.ImageUrl) - len("?alt=media")
	// fmt.Println(endUrl)
	url := strings.Split(strings.Join(strings.Split(res1.ImageUrl, base), ""), "?")[0]
	fmt.Println(url)
	storage, err := helpers.FirebaseApp.Storage(context.Background())
	if err != nil {
		helpers.BadResponse(c, err.Error())
		return
	}

	bucket, err := storage.Bucket(os.Getenv("BUCKET_NAME"))
	if err != nil {
		helpers.BadResponse(c, err.Error())
		return
	}

	obj := bucket.Object("demo/" + url)

	if err := obj.Delete(context.TODO()); err != nil {
		helpers.BadResponse(c, err.Error())
		return
	}

	res2, err := collection.DeleteOne(context.TODO(), bson.M{"_id": id})

	if err != nil {
		helpers.BadResponse(c, err.Error())
		return
	}

	if res2.DeletedCount == 0 {
		helpers.BadResponse(c, "not deleted !!")
		return
	}

	helpers.SuccessResponse(c, "item deleted successfully", res2.DeletedCount)

}

func UpdateItem(c *gin.Context) {
	collection := helpers.DB.Collection("item")
	var params string = c.Param("id")
	id, err := primitive.ObjectIDFromHex(params)

	if err != nil {
		helpers.BadResponse(c, err.Error())
		return
	}

	if count, _ := collection.CountDocuments(context.TODO(), bson.M{"_id": id}); count == 0 {
		helpers.BadResponse(c, "document not found !!")
		return
	}
	var res1 models.ItemModel

	if err := collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&res1); err != nil {
		helpers.BadResponse(c, err.Error())
		return
	}

	var input itemUpdate

	if err := c.Bind(&input); err != nil {
		helpers.BadResponse(c, err.Error())
		return
	}

	image, _ := c.FormFile("image")

	var updateFields bson.D

	if input.ItemName != "" {
		updateFields = append(updateFields, bson.E{Key: "item_name", Value: input.ItemName})
	}

	if input.Price != "" {
		updateFields = append(updateFields, bson.E{Key: "price", Value: input.Price})
	}

	if input.Qty != "" {
		updateFields = append(updateFields, bson.E{Key: "qty", Value: input.Qty})
	}

	if input.SubQty != "" {
		var subQty []models.SubItem
		err := json.Unmarshal([]byte(input.SubQty), &subQty)
		if err != nil {
			helpers.BadResponse(c, err.Error())
			return
		}
		updateFields = append(updateFields, bson.E{Key: "sub_qty", Value: subQty})
	}

	if input.Addons != "" {

		var i []string
		if err := json.Unmarshal([]byte(input.Addons), &i); err != nil {
			helpers.BadResponse(c, err.Error())
			return
		}

		updateFields = append(updateFields, bson.E{Key: "addons", Value: i})
	}

	if image != nil {
		base := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/%s/o/%s%s", os.Getenv("BUCKET_NAME"), os.Getenv("BUCKET_FOLDER"), "%2f")
		imageUrl := strings.Split(strings.Join(strings.Split(res1.ImageUrl, base), ""), "?alt=media")[0]

		storage, err := helpers.FirebaseApp.Storage(context.Background())
		if err != nil {
			helpers.BadResponse(c, err.Error())
			return
		}
		bucket, err := storage.Bucket(os.Getenv("BUCKET_NAME"))
		if err != nil {
			helpers.BadResponse(c, err.Error())
			return
		}

		obj := bucket.Object(os.Getenv("BUCKET_FOLDER") + "/" + imageUrl)
		if err := obj.Delete(context.TODO()); err != nil {
			helpers.BadResponse(c, err.Error())
			return
		}
		obj2 := bucket.Object(os.Getenv("BUCKET_FOLDER") + "/" + image.Filename)

		f, err := image.Open()

		if err != nil {
			helpers.BadResponse(c, err.Error())
			return
		}
		defer f.Close()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		writer := obj2.NewWriter(ctx)
		defer writer.Close()
		id := uuid.New().String()

		writer.ObjectAttrs.Metadata = map[string]string{"firebaseStorageDownloadToken": id}
		if _, err := io.Copy(writer, f); err != nil {
			helpers.BadResponse(c, err.Error())
			return
		}

		finalImageUrl := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/%s/o/%s%s%s?alt=media&token=%s", os.Getenv("BUCKET_NAME"), os.Getenv("BUCKET_FOLDER"), "%2f", image.Filename, id)
		updateFields = append(updateFields, bson.E{Key: "image_url", Value: finalImageUrl})
	}

	update := bson.D{{Key: "$set", Value: updateFields}}
	filter := bson.D{{Key: "_id", Value: id}}

	res, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		helpers.BadResponse(c, err.Error())
		return
	}

	helpers.SuccessResponse(c, "item updated successfully !", res.ModifiedCount)

}

func AddBill(c *gin.Context) {
	collection := helpers.DB.Collection("bill")
	var input Bill
	if err := c.Bind(&input); err != nil {
		helpers.BadResponse(c, err.Error())
		return
	}

	var items []models.BillItemModel

	if input.Items != "" {
		if err := json.Unmarshal([]byte(input.Items), &items); err != nil {
			helpers.BadResponse(c, err.Error())
			return
		}
	}

	var parcel []models.BillItemModel

	if input.Parcels != "" {
		if err := json.Unmarshal([]byte(input.Parcels), &parcel); err != nil {
			helpers.BadResponse(c, err.Error())
			return
		}
	}

	bill := models.BillModel{CustomerNumber: input.CustomerNumber, CustomerName: input.CustomerName, Date: input.Date, BilledTime: input.BilledTime, PaidTime: input.PaidTime, PaymentStatus: input.PaymentStatus, PaymentType: input.PaymentType, Discount: input.Discount, TotalItemAmount: input.TotalItemAmount, TotalParcelAmount: input.TotalParcelAmount, NetTotal: input.NetTotal, Type: input.Type, Items: items, Parcels: parcel}

	res, err := collection.InsertOne(context.TODO(), &bill)
	if err != nil {
		helpers.BadResponse(c, err.Error())
		return
	}

	helpers.SuccessResponse(c, "bill generated successfully !", res.InsertedID)

}

func GetBill(c *gin.Context) {
	collection := helpers.DB.Collection("bill")
	type Format struct {
		From string `form:"from"`
		To   string `form:"to"`
	}

	var input Format

	if err := c.Bind(&input); err != nil {
		helpers.BadResponse(c, err.Error())
		return
	}

	if input.From != "" && input.To != "" {
		filter := bson.M{
			"date": bson.M{
				"$gte": input.From,
				"$lte": input.To,
			},
		}

		cursour, err := collection.Find(context.TODO(), filter)
		if err != nil {
			helpers.BadResponse(c, err.Error())
			return
		}
		var bills []models.BillModel
		for cursour.Next(context.TODO()) {
			var bill models.BillModel
			if err := cursour.Decode(&bill); err != nil {
				helpers.BadResponse(c, err.Error())
				return
			}
			bills = append(bills, bill)
		}
		defer cursour.Close(context.TODO())

		helpers.SuccessResponse(c, "bill from "+input.From+" to "+input.To, bills)
		return

	} else {
		var bills []models.BillModel

		cursor, err := collection.Find(context.Background(), bson.M{})
		if err != nil {
			helpers.BadResponse(c, err.Error())
			return
		}

		for cursor.Next(context.TODO()) {
			var bill models.BillModel
			if err := cursor.Decode(&bill); err != nil {
				helpers.BadResponse(c, err.Error())
				return
			}
			bills = append(bills, bill)
		}
		defer cursor.Close(context.TODO())

		helpers.SuccessResponse(c, "bills list", bills)
	}

}
