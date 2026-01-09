package ai 

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/manas-011/code-editor-backend/service"
)

type AskRequest struct {
	Question string `json:"question"`
}

func AskHandler(c *gin.Context) {
	var req AskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	answer, err := service.AskAI(req.Question)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"answer": answer,
	})
}
