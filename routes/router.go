package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	middle"github.com/Psnsilvino/CaluFestas-Site-api/middleware"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"}, // CaluFestas-Site (5173) and calu-chat (3000)
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))

	api := router.Group("/api") // Agrupa todas as rotas dentro de /api
	ClientRoutes(api)   // Adiciona rotas de usuários
    ProductRoutes(api)

	// Agora criamos um grupo protegido pelo AuthMiddleware
    protected := api.Group("/")
    protected.Use(middle.AuthMiddleware()) // tudo que estiver aqui exigirá o JWT

    // Rotas protegidas
    PrivateClientRoutes(protected)
    PrivateProductRoutes(protected)
    LocationRoutes(protected)
	

	return router
}