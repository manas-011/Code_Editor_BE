package ai 

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/manas-011/code-editor-backend/service"
)


type AskAIRequest struct {
	Prompt   string `json:"prompt" binding:"required"`
	Language string `json:"language"`
}

func AskAI(c *gin.Context) {

	var req AskAIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload",
		})
		return
	}

	prompt := buildPrompt(req.Language, req.Prompt)

	response, err := service.CallHuggingFace(prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"output": response,
	})
}

// prompt engineering (important for quality)
func buildPrompt(language, userPrompt string) string {

	if language == "" {
		return userPrompt
	}

	var sb strings.Builder
	sb.WriteString("You are an expert ")
	sb.WriteString(language)
	sb.WriteString(" developer.\n")
	sb.WriteString("Write clean, production-ready code.\n\n")
	sb.WriteString(userPrompt)

	return sb.String()
}
