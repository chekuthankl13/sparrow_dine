package main

import (
	"os"

	"github.com/chekuthankl13/sparrow_dine/helpers"
	"github.com/chekuthankl13/sparrow_dine/routes"
	"github.com/gin-gonic/gin"
)

func init() {
	helpers.LoadEnv()
	helpers.ConnectMongoDb()
}

func main() {

	port := os.Getenv("PORT")

	if port == "" {
		port = "3035"
	}

	r := gin.Default()

	routes.AdminRoutes(r)

	r.Run(":" + port)
}
