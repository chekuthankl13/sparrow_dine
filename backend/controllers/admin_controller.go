package controllers

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/chekuthankl13/sparrow_dine/helpers"
	"github.com/chekuthankl13/sparrow_dine/middlewares"
	"github.com/chekuthankl13/sparrow_dine/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

	var inputCred adminCred
	if err := c.Bind(&inputCred); err != nil {
		helpers.BadResponse(c, err.Error())
		return
	}

	count, _ := collection.CountDocuments(context.Background(), bson.M{"user_name": inputCred.Username})
	fmt.Println("count -", count)
	if count >= 1 {
		helpers.BadResponse(c, "username already exist !!")
		return
	}

	hashPsw, err := middlewares.HashPassword(inputCred.Password)
	if err != nil {
		helpers.BadResponse(c, err.Error())
		return
	}

	cred := models.AdminCredModel{UserName: inputCred.Username, Password: hashPsw}

	result, err := collection.InsertOne(context.Background(), &cred)
	if err != nil {
		helpers.BadResponse(c, err.Error())
		return
	}

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
