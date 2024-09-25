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
)

type adminCred struct {
	Username string `form:"user_name" binding:"required"`
	Password string `form:"password" binding:"required"`
}

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
	// id, _ := primitive.ObjectIDFromHex(result.ID.String())
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
