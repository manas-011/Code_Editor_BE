package auth

import (
	"context"
	"net/http"

	"github.com/manas-011/code-editor-backend/config"
	"github.com/manas-011/code-editor-backend/util"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}


func Login(c *gin.Context) {
	var req LoginRequest

	// Parse request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	// Find user by email
	var user struct {
		ID       primitive.ObjectID `bson:"_id"`
		Email    string             `bson:"email"`
		Password string             `bson:"password"`
	}

	err := config.DB.Collection("verified_users").FindOne(
		context.TODO(),
		bson.M{"email": req.Email},
	).Decode(&user)

	if err != nil {
		// Same error message to prevent user enumeration
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid email or password",
		})
		return
	}

	// Compare password hash
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid email or password",
		})
		return
	}

	// Generate JWT (7 days)
	token, err := util.GenerateJWT(user.ID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "token generation failed",
		})
		return
	}

	// Set JWT cookie
	c.SetCookie(
		"auth_token",
		token,
		7*24*3600, // 7 days
		"/",
		"",
		true, // Secure
		true, // HttpOnly
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
	})
}
