package router

import (
	"btpn-syariah-final-project/controllers"
	"btpn-syariah-final-project/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Photo authentication
	r.POST("/users/register", controllers.Signup)
	r.POST("/users/login", controllers.Login)	
	r.PUT("/users/:userId", middlewares.RequireAuth, controllers.UpdateUser)	
	r.DELETE("/users/:userId", middlewares.RequireAuth, controllers.DeleteUser)	

	// Photo routes
	r.GET("/photos", middlewares.RequireAuth, controllers.GetPhoto)	
	r.POST("/photos", middlewares.RequireAuth, controllers.StorePhoto)
	r.PUT("/:photoId", middlewares.RequireAuth, controllers.UpdatePhoto)	
	r.DELETE("/:photoId", middlewares.RequireAuth, controllers.DeletePhoto)

	return r
}