package controller

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func HealthCheck(c *gin.Context){
	// c = request + response + helper utilities
	c.JSON(http.StatusOK, gin.H{
		"status":"ok",
	})
}