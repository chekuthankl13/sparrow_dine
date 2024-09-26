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
}
