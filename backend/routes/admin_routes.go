package routes

import (
	"github.com/chekuthankl13/sparrow_dine/controllers"
	"github.com/chekuthankl13/sparrow_dine/middlewares"
	"github.com/gin-gonic/gin"
)

func AdminRoutes(r *gin.Engine) {
	r.POST("/login", controllers.AdminLogin)
	r.POST("/register", controllers.AdminRegister)

	staff := r.Group("/staff", middlewares.ValidateToken())
	staff.POST("", controllers.StaffCreate).GET("", controllers.GetStaffs).DELETE("/:id", controllers.DeleteStaff).PUT("/:id", controllers.EditStaff)

	kitchen := r.Group("/kitchen", middlewares.ValidateToken())
	kitchen.POST("", controllers.CreateKitchen).GET("", controllers.GetKitchen).DELETE("/:id", controllers.DeleteKitchen)

	table := r.Group("/table", middlewares.ValidateToken())
	table.POST("", controllers.CreateTable).GET("", controllers.GetTable).DELETE("/:id", controllers.DeleteTable)

	item := r.Group("/item", middlewares.ValidateToken())

	item.POST("", controllers.CreateItem).GET("", controllers.GetItems)
}
