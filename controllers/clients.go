package controllers

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/Psnsilvino/CaluFestas-Site-api/database"
	"github.com/Psnsilvino/CaluFestas-Site-api/models"
	"github.com/Psnsilvino/CaluFestas-Site-api/utils/email"
	"github.com/Psnsilvino/CaluFestas-Site-api/utils/otputil"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var emailSender email.EmailSender

func InitEmailSender(sender email.EmailSender) {
	emailSender = sender
}

func Register(c *gin.Context) {
	var client models.Client
	if err := c.ShouldBindJSON(&client); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(client.Senha), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	client.Senha = string(hashedPassword)
	client.ID = primitive.NewObjectID()

	_, err = database.DB.Database(os.Getenv("DB_NAME")).Collection("clients").InsertOne(context.Background(), client)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func GetClients(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := database.DB.Database(os.Getenv("DB_NAME")).Collection("clients").Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	defer cursor.Close(ctx)

	var clients []models.Client
	if err := cursor.All(ctx, &clients); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing users"})
		return
	}

	c.JSON(http.StatusOK, clients)
}


func Login(c *gin.Context) {
	var loginData models.Client
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var client models.Client
	err := database.DB.Database(os.Getenv("DB_NAME")).Collection("clients").FindOne(ctx, bson.M{"email": loginData.Email}).Decode(&client)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(client.Senha), []byte(loginData.Senha))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"nome": client.Nome,
		"email": client.Email,
		"senha": client.Senha,
	})
}

func ForgotPassword(c *gin.Context)  {

	var email string
	var emailModel models.Email
	if err := c.ShouldBindJSON(&emailModel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	email = emailModel.Email

	// verificar se o cliente existe
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var client models.Client
	err := database.DB.Database(os.Getenv("DB_NAME")).Collection("clients").FindOne(ctx, bson.M{"email": email}).Decode(&client)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Generate reset token (OTP)
	resetToken, err := otputil.GenerateOTP(6)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating otp"})
		return 
	}

	expiresAt := time.Now().UTC().Add(15 * time.Minute)
	passwordResetEntry := &models.PasswordResetEntry{
		Email:     client.Email,
		OTPCode:   resetToken,
		ExpiresAt: expiresAt,
	}

	err = CreatePasswordResetEntry(c, passwordResetEntry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating entry"})
		return 
	}

	err = emailSender.SendPasswordResetToken(client.Email, resetToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return 
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})

}