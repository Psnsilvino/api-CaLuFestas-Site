package controllers

import (
	"context"
	"net/http"
	"os"

	"github.com/Psnsilvino/CaluFestas-Site-api/database"
	"github.com/Psnsilvino/CaluFestas-Site-api/models"
	"github.com/gin-gonic/gin"
)

func CreatePasswordResetEntry(c *gin.Context, entry *models.PasswordResetEntry) error {

	_, err := database.DB.Database(os.Getenv("DB_NAME")).Collection("senhasEsquecidas").InsertOne(context.Background(), entry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register forgotten password"})
		return err
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Forgot password registered successfully"})

	return nil
}