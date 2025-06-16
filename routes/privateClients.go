package routes

import (
	"github.com/Psnsilvino/CaluFestas-Site-api/controllers"
	"github.com/gin-gonic/gin"
)

func PrivateClientRoutes(r *gin.RouterGroup) {
	privateClients := r.Group("/privateClients")
	{
		privateClients.GET("/", controllers.GetClients)
		privateClients.GET("/me", controllers.Me)
	}
}