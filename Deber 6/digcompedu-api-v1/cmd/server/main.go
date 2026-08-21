// cmd/server/main.go

package main

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"digcompedu-api/internal/api"
	"digcompedu-api/internal/config"
	"digcompedu-api/internal/evaluator"
	"digcompedu-api/internal/llm"
	"digcompedu-api/internal/storage"
)

//go:embed all:web
var webFS embed.FS

func main() {
	if err := run(); err != nil {
		log.Fatalf("error fatal: %v", err)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	repo, err := storage.NewSQLiteRepository(cfg.DBPath)
	if err != nil {
		return err
	}

	//llmClient := llm.NewAnthropicClient(cfg.AnthropicKey, cfg.AnthropicModel)
	llmClient := llm.NewGroqClient(cfg.GroqKey, cfg.GroqModel)
	service := evaluator.NewService(llmClient, repo)
	handler := api.NewHandler(service, repo)

	// fs.Sub "pela" el prefijo "web" para que index.html quede en la raíz del FS
	staticFS, err := fs.Sub(webFS, "web")
	if err != nil {
		return err
	}

	router := api.NewRouter(handler, staticFS)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		log.Printf("servidor escuchando en :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("error al iniciar servidor: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("apagando servidor...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return srv.Shutdown(ctx)
}