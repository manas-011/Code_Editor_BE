package auth

import (
	"context"
	"net/http"
	"time"

	"auth-app/internal/db"
	"auth-app/internal/models"
	"auth-app/internal/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VerifyOTPRequest struct {
	OTP string `json:"otp"`
}

func VerifyOTP(c *gin.Context) {
	var req VerifyOTPRequest

	if err := c.ShouldBindJSON(&req); err != nil || req.OTP == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid otp"})
		return
	}

	cookie, err := c.Cookie("signup_session")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "session expired"})
		return
	}

	idStr, err := services.Decrypt(cookie)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
		return
	}

	objID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session id"})
		return
	}

	var tempUser models.TempUser
	err = db.TempUserCollection.FindOne(
		context.TODO(),
		bson.M{"_id": objID},
	).Decode(&tempUser)

	if err != nil || time.Now().After(tempUser.ExpiresAt) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "otp expired"})
		return
	}

	if tempUser.OTP != req.OTP {
		c.JSON(http.StatusBadRequest, gin.H{"error": "otp verification failed"})
		return
	}

	user := models.User{
		Email:     tempUser.Email,
		Password:  tempUser.Password,
		CreatedAt: time.Now(),
	}

	res, err := db.UserCollection.InsertOne(context.TODO(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user creation failed"})
		return
	}

	db.TempUserCollection.DeleteOne(context.TODO(), bson.M{"_id": objID})

	token, err := services.GenerateJWT(res.InsertedID.(primitive.ObjectID).Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
		return
	}

	c.SetCookie(
		"auth_token",
		token,
		7*24*3600,
		"/",
		"",
		true,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "signup successful",
	})
}
