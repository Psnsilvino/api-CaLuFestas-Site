package routes

import (
	"github.com/Psnsilvino/CaluFestas-Site-api/controllers"
	"github.com/gin-gonic/gin"
)

func LocationRoutes(r *gin.RouterGroup) {
	location := r.Group("/locations")
	{
		location.POST("/", controllers.CreateLocation)
		location.GET("/", controllers.GetLocations)
		location.PUT("/:id", controllers.UpdateLocation)
		location.POST("/:id/delete", controllers.DeleteLocation)
		location.POST("/cliente", controllers.LocationsByClient)
	}
}