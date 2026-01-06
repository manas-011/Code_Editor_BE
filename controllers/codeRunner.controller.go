package controller 

import (
	"net/http"
	"time"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/manas-011/code-editor-backend/models"
	"github.com/manas-011/code-editor-backend/controllers/executor"
	"github.com/manas-011/code-editor-backend/middlewares/limiter"
)



func ExecCode(c *gin.Context){
	var req model.ExecuteRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return 
	}

	// concurrency limit
	limiter.ExecSemaphore <- struct{}{}
	defer func(){ <-limiter.ExecSemaphore }()

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	stdout, stderr, err := executor.Execute(ctx, req.Language, req.Code, req.Input)

	result := models.ExecuteResult {
		Status: "success",
		Stdout: stdout,
		Stderr: stderr
	}

	if ctx.Err() == context.DeadlineExceeded {
		result.Status == "TLE"
	}else if err != nil {
		result.Status = "RE"
	}

	c.JSON(http.StatusOK, result)
}