package route

import (
	"github.com/gin-gonic/gin"
	"github.com/manas-011/code-editor-backend/controller"
	"github.com/manas-011/code-editor-backend/controller/ai"
	"github.com/manas-011/code-editor-backend/controller/auth"
)

func SetupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	// middleware
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	r.GET("/health", controller.HealthCheck)
	r.POST("/app/exec", controller.ExecCode)
	r.POST("/app/askAI", ai.AskHandler)
	r.POST("/auth/signup", auth.Signup)
	r.POST("/auth/verifyOtp", auth.VerifyOTP)
	r.POST("/auth/login", auth.Login)

	return r 
}