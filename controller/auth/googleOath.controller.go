package auth

import (
	"context"
	"net/http"
	"fmt"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/idtoken"
	"github.com/manas-011/code-editor-backend/config"
)

type GoogleAuthRequest struct {
	Code string `json:"code"`
}

func GoogleAuthHandler(c *gin.Context) {
	var req GoogleAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Exchange code for tokens
	token, err := config.GoogleOAuthConfig.Exchange(
		context.Background(),
		req.Code,
	)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "code exchange failed"})
		return
	}

	// Extract ID token
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing id token"})
		return
	}

	// Verify ID token
	payload, err := idtoken.Validate(
		context.Background(),
		rawIDToken,
		config.GoogleOAuthConfig.ClientID,
	)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid id token"})
		return
	}

	// User info
	email, _ := payload.Claims["email"].(string)
	name, _ := payload.Claims["name"].(string)
	picture, _ := payload.Claims["picture"].(string)

	if email == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "email not available"})
		return
	}

	fmt.Println(email)
	fmt.Println(name)
	fmt.Println(picture) 
	

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
