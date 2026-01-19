package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/manas-011/code-editor-backend/controller/executor"
	"github.com/manas-011/code-editor-backend/middleware/limiter"
	"github.com/manas-011/code-editor-backend/model"
)

func ExecCode(c *gin.Context) {
	var req model.ExecuteRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// concurrency limit
	limiter.ExecSemaphore <- struct{}{}
	defer func() { <-limiter.ExecSemaphore }()

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	result, err := executor.Execute(ctx, req.Language, req.Code, req.Input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "SYSTEM_ERROR",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}
