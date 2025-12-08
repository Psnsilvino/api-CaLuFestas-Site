package routes

import (
	"github.com/Psnsilvino/CaluFestas-Site-api/controllers"
	"github.com/gin-gonic/gin"
)

func PrivateProductRoutes(r *gin.RouterGroup) {
	privateProducts := r.Group("/privateProducts")
	{
		privateProducts.POST("/register", controllers.CreateProduct)
		privateProducts.DELETE("/:nome", controllers.DeleteProduct)
		privateProducts.PUT("/:nome", controllers.UpdateProduct)
	}
}