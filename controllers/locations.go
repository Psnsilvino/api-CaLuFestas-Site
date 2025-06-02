package controllers

import (
	"context"
	"net/http"
	"os"

	"github.com/Psnsilvino/CaluFestas-Site-api/database"
	"github.com/Psnsilvino/CaluFestas-Site-api/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateLocation(c *gin.Context) {
	var locacao models.Locacao
	
	if err := c.ShouldBindJSON(&locacao); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	locacao.ID = primitive.NewObjectID()

	for _, item := range locacao.Items {
    	filter := bson.M{"nome": item.Nome}
    	update := bson.M{
        	"$inc": bson.M{
        	    "quantidadeemlocacao": item.Quantidade, // subtrai do estoque
        	},
    	}

    	result := database.DB.Database(os.Getenv("DB_NAME")).Collection("produtos").FindOneAndUpdate(context.Background(), filter, update)
    	if result.Err() != nil {
        	c.JSON(http.StatusInternalServerError, gin.H{
            	"error": item,
        	})
        	return
    	}
	}

	_, err := database.DB.Database(os.Getenv("DB_NAME")).Collection("locations").InsertOne(context.Background(), locacao)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create location"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Location registered successfully"})
}