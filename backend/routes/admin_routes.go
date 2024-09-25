package routes

import (
	"github.com/chekuthankl13/sparrow_dine/controllers"
	"github.com/gin-gonic/gin"
)

func AdminRoutes(r *gin.Engine) {
	r.GET("/", controllers.AdminLogin)
	r.POST("/login", controllers.AdminLogin)
	r.POST("/register", controllers.AdminRegister)
}
