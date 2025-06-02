package routes

import (
	"github.com/Psnsilvino/CaluFestas-Site-api/controllers"
	"github.com/gin-gonic/gin"
)

func LocationRoutes(r *gin.RouterGroup) {
	location := r.Group("/locations")
	{
		location.POST("/", controllers.CreateLocation)
	}
}