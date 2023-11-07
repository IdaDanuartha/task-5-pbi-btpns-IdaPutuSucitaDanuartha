package main

import (
	"btpn-syariah-final-project/controllers"
	"btpn-syariah-final-project/database"
	"btpn-syariah-final-project/middlewares"

	"github.com/gin-gonic/gin"
)

func init() {
	database.LoadEnvVariables()
	database.ConnectToDb()
	database.SyncDatabase()
}

func main() {
	r := gin.Default()

	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.GET("/validate", middlewares.RequireAuth, controllers.Validate)

	r.Run()
}