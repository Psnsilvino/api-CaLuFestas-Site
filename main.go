package main

import (
	"log"
	"os"
	"strconv"

	"github.com/Psnsilvino/CaluFestas-Site-api/controllers"
	"github.com/Psnsilvino/CaluFestas-Site-api/database"
	"github.com/Psnsilvino/CaluFestas-Site-api/routes"
	"github.com/Psnsilvino/CaluFestas-Site-api/utils/email"
	"github.com/joho/godotenv"
)

func toInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("Invalid SMTP_PORT: %v", err)
	}
	return i
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Connect to MongoDB
	database.ConnectDB()
	
	// Inicializar o sender global
	sender := email.NewSender(os.Getenv("SMTP_HOST"), toInt(os.Getenv("SMTP_PORT")), os.Getenv("SMTP_EMAIL"), os.Getenv("SMTP_PSW"))
	controllers.InitEmailSender(sender)

	// Setup routes
	r := routes.SetupRouter()

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(r.Run(": " + port))
}