package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)


func Logout(c *gin.Context) {

	// Clear auth token cookie
	c.SetCookie(
		"auth_token",
		"",
		-1, // expire immediately
		"/",
		"",
		true, // Secure
		true, // HttpOnly
	)

	// Clear signup session cookie (optional safety)
	c.SetCookie(
		"signup_session",
		"",
		-1,
		"/",
		"",
		true,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "logout successful",
	})
}
