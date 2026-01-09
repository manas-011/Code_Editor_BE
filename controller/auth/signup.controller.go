package auth

import (
	"context"
	"net/http"
	"time"

	"auth-app/internal/db"
	"auth-app/internal/model"
	"auth-app/internal/service"
	"auth-app/internal/validators"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Signup(c *gin.Context) {
	var req SignupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if !validators.IsValidEmail(req.Email) ||
		!validators.IsStrongPassword(req.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email or password"})
		return
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "password hashing failed"})
		return
	}

	otp := services.GenerateOTP()

	tempUser := models.TempUser{
		Email:     req.Email,
		Password:  string(hashedPwd),
		OTP:       otp,
		ExpiresAt: time.Now().Add(2 * time.Minute),
		CreatedAt: time.Now(),
	}

	res, err := db.TempUserCollection.InsertOne(context.TODO(), tempUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	services.SendOTPEmail(req.Email, otp)

	encID, err := services.Encrypt(res.InsertedID.(primitive.ObjectID).Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "encryption failed"})
		return
	}

	c.SetCookie(
		"signup_session",
		encID,
		120, // 2 minutes
		"/",
		"",
		true,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "otp sent, please verify",
	})
}
