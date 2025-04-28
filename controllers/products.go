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
)

func GetProducts(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := database.DB.Database(os.Getenv("DB_NAME")).Collection("produtos").Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}
	defer cursor.Close(ctx)

	var products []models.Product
	if err := cursor.All(ctx, &products); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing products"})
		return
	}

	var productResponses []models.ProductResponse
	for _, product := range products {
    productResponses = append(productResponses, models.ProductResponse{
        ID:                  product.ID.Hex(),
        Nome:                product.Nome,
        Categoria:           product.Categoria,
        Quantidade:          product.Quantidade,
        QuantidadeEmLocacao: product.QuantidadeEmLocacao,
        Preco:               product.Preco,
        CreatedAt:           product.CreatedAt,
        UpdatedAt:           product.UpdatedAt,
    })
}

	c.JSON(http.StatusOK, productResponses)
}

func CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := database.DB.Database(os.Getenv("DB_NAME")).Collection("produtos").InsertOne(context.Background(), product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Product registered successfully"})
}