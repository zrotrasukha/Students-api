package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zrotrasukha/Students-api/http/handlers/student"
	"github.com/zrotrasukha/Students-api/internal/config"
)

func main() {
	// load config
	cfg := config.MustLoad()
	// logger
	// database setup
	// setup router
	router := http.NewServeMux()

	router.HandleFunc("GET /api/students", student.New())
	// setup server

	server := http.Server{
		Addr:    cfg.HttpServer.Addr,
		Handler: router,
	}

	fmt.Println("Server started successfully")

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Failed to start server", err)
		}
	}()

	// graceful shutdown
	<-done

	slog.Info("server is shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown:", slog.String("error", err.Error()))
	}

}
