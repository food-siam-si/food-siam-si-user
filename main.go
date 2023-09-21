package main

import (
	"food-siam-si/food-siam-si-user/controllers"
	"food-siam-si/food-siam-si-user/models"

	"github.com/gin-gonic/gin"
)

func main() {

	models.ConnectDataBase()

	r := gin.Default()

	userGroup := r.Group("/users")

	userGroup.POST("/register", controllers.Register)
	userGroup.POST("/login", controllers.Login)
	userGroup.GET("/verify", controllers.CurrentUser)

	r.Run(":8080")

}
