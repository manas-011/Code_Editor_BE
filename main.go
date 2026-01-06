package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/manas-011/code-editor-backend/configs"
	"github.com/manas-011/code-editor-backend/routes"
)

func main(){
	// Load configuration
	cfg := config.Load()

	// Setup router (Gin)
	router := route.SetupRouter()

	// Create HTTP server
	server := &http.Server{
		Addr: ":" + cfg.Port,
		Handler: router,
		ReadTimeout: time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
		IdleTimeout: 60 * time.Second,
	}

	// Start server in a goroutine
	go func(){
		log.Printf("Server started on port %s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server error: %v", err) 
		}
	}()

	// Graceful shutdown handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")

}