package routes

import (
	"github.com/Psnsilvino/CaluFestas-Site-api/controllers"
	"github.com/gin-gonic/gin"
)

func ProductRoutes(r *gin.RouterGroup) {
	products := r.Group("/products")
	{
		products.GET("/", controllers.GetProducts)
		products.POST("/register", controllers.CreateProduct)
	}
}