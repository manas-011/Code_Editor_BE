package route

import (
	"github.com/gin-gonic/gin"
	"github.com/manas-011/code-editor-backend/controller"
)

func SetupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	// middleware
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	r.GET("/health", controller.HealthCheck)
	r.POST("/app/exec", controller.ExecCode)

	return r 
}