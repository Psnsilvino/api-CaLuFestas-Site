package routes

import (
	"github.com/Psnsilvino/CaluFestas-Site-api/controllers"
	"github.com/gin-gonic/gin"
)

func ClientRoutes(r *gin.RouterGroup) {
	clients := r.Group("/clients")
	{
		clients.POST("/", controllers.Register)
		clients.POST("/login", controllers.Login)
		clients.POST("/ForgotPassword", controllers.ForgotPassword)
		clients.POST("/verifyCode", controllers.VerifyCode)
		clients.POST("/ResetPassword", controllers.UpdatePassword)
		
	}
}