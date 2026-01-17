package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/aira-id/gribe/internal/config"
	"github.com/aira-id/gribe/internal/delivery/websocket"
	"github.com/aira-id/gribe/internal/usecase"
)

func main() {
	// Load configuration from environment
	cfg := config.Load()

	// Log configuration (without sensitive data)
	log.Printf("Starting Gribe STT Server")
	log.Printf("Port: %s", cfg.Server.Port)
	log.Printf("Max audio buffer size: %d bytes", cfg.Audio.MaxBufferSize)
	log.Printf("Max connections per IP: %d", cfg.Rate.MaxConnectionsPerIP)

	if len(cfg.Server.AllowedOrigins) == 0 {
		log.Println("Allowed origins: * (all)")
	} else {
		log.Printf("Allowed origins: %v", cfg.Server.AllowedOrigins)
	}

	if len(cfg.Auth.APIKeys) == 0 {
		log.Println("Authentication: disabled (no API keys configured)")
	} else {
		log.Printf("Authentication: enabled (%d API key(s) configured)", len(cfg.Auth.APIKeys))
	}

	// Initialize Usecase with configuration
	sessionUsecase := usecase.NewSessionUsecaseWithConfig(cfg)

	// Initialize Delivery Handler
	wsHandler := websocket.NewHandler(sessionUsecase, cfg)

	// Set up routes
	http.Handle("/v1/realtime", wsHandler)

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Graceful shutdown handling
	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Server shutting down...")
		wsHandler.Close()
		done <- true
	}()

	// Start server
	addr := ":" + cfg.Server.Port
	log.Printf("Server listening on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("Server error:", err)
	}

	<-done
	log.Println("Server stopped")
}
