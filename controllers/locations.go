package controllers

import (
	"context"
	"net/http"
	"os"
	"time"

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
	locacao.Estado = "Em analise"

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

func GetLocations(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := database.DB.Database(os.Getenv("DB_NAME")).Collection("locations").Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}
	defer cursor.Close(ctx)

	var locations []models.Locacao
	if err := cursor.All(ctx, &locations); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing products"})
		return
	}

	c.JSON(http.StatusOK, locations)
}

func DeleteLocation(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var payload struct {
		Items []models.Item `json:"items"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao processar itens enviados"})
		return
	}

	// Devolver itens ao estoque
	for _, item := range payload.Items {
		filter := bson.M{"nome": item.Nome}
		update := bson.M{
			"$inc": bson.M{
				"quantidadeemlocacao": -item.Quantidade,
			},
		}
		_, err := database.DB.Database(os.Getenv("DB_NAME")).Collection("produtos").UpdateOne(context.Background(), filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao devolver produto: " + item.Nome})
			return
		}
	}

	// Deleta a locação
	result, err := database.DB.Database(os.Getenv("DB_NAME")).Collection("locations").DeleteOne(context.Background(), bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar locação"})
		return
	}
	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Locação não encontrada"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Locação excluída com sucesso"})
}


func UpdateLocation(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var updatedData models.Locacao
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao fazer binding do JSON"})
		return
	}

	// Devolve itens ao estoque se a locação for recusada ou concluída
	if updatedData.Estado == "Recusada" || updatedData.Estado == "Concluida" {
		for _, item := range updatedData.Items {
			filter := bson.M{"nome": item.Nome}
			update := bson.M{
				"$inc": bson.M{
					"quantidadeemlocacao": -item.Quantidade,
				},
			}
			_, err := database.DB.Database(os.Getenv("DB_NAME")).Collection("produtos").UpdateOne(context.Background(), filter, update)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao devolver produto: " + item.Nome})
				return
			}
		}
	}

	// Atualiza estado da locação
	update := bson.M{
		"$set": bson.M{
			"estado": updatedData.Estado,
		},
	}

	collection := database.DB.Database(os.Getenv("DB_NAME")).Collection("locations")
	result, err := collection.UpdateOne(context.Background(), bson.M{"_id": objID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar locação"})
		return
	}
	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Locação não encontrada"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Locação atualizada com sucesso"})
}


func LocationsByClient(c *gin.Context) {
	var request models.ClientLocation

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cliente field is required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"email": request.Email}

	cursor, err := database.DB.Database(os.Getenv("DB_NAME")).Collection("locations").Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch locations"})
		return
	}
	defer cursor.Close(ctx)

	var locations []models.Locacao
	if err := cursor.All(ctx, &locations); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing locations"})
		return
	}

	if len(locations) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"message": request})
		return
	}

	c.JSON(http.StatusOK, locations)
}

