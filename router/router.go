package router

import (
	"dbmsbackend/controller"
	"dbmsbackend/middleware"
	"dbmsbackend/util"

	"github.com/gin-gonic/gin"
)

func setupHeader(c *gin.Context) {
	// Add header Access-Control-Allow-Origin
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Max-Age", "86400")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, UPDATE")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Max")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(200)
	} else {
		c.Next()
	}
}

func SetupRouter(router *gin.Engine, config *util.Config) {

	router.Use(setupHeader)

	//API route for version 1
	v1 := router.Group("/api/v1")

	userController := new(controller.User)
	userController.Initialize(config)

	productController := new(controller.Product)
	productController.Initialize(config)

	orderController := new(controller.Order)
	orderController.Initialize(config)

	v1.POST("users", userController.New)
	v1.PATCH("users/:id", middleware.Auth(config), userController.Update)
	v1.POST("users/signIn", userController.Login(config))
	v1.GET("users/me", middleware.Auth(config), userController.GetCurrentUser)

	v1.POST("products", middleware.Auth(config), productController.New)
	v1.PATCH("products/:id", middleware.Auth(config), productController.Update)
	v1.GET("products", productController.Query)
	v1.GET("products/:id", productController.GetByID)
	v1.DELETE("products/:id", middleware.Auth(config), productController.Delete)

	v1.POST("orders", middleware.Auth(config), orderController.New)
	v1.GET("orders", middleware.Auth(config), orderController.Query)
	v1.GET("orders/:id", middleware.Auth(config), orderController.GetByID)
	v1.DELETE("orders/:id", middleware.Auth(config), orderController.Delete)

	v1.POST("images", productController.NewImage(config), middleware.Auth(config))

}
