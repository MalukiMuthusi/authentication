package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	log "github.com/malukimuthusi/authentication/pkg/internal/logger"
)

func main() {
	r := mux.NewRouter()

	// CORS
	methods := []string{http.MethodPost, http.MethodGet, http.MethodOptions}
	origins := []string{"*"}

	co := cors.New(cors.Options{
		AllowedOrigins:   origins,
		AllowCredentials: false,
		AllowedMethods:   methods,
	})
	handler := co.Handler(r)

	r.Use(mux.CORSMethodMiddleware(r))

	port := viper.GetString("PORT")
	if port == "" {
		log.Info("port not set, setting it to 8080")
		port = "8080"
	}

	srv := &http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%s", port),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      handler,
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {

		log.Info(fmt.Sprintf("Starting server %s\n", srv.Addr))

		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server", zap.Error(err))
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(viper.GetInt64("wait")))
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Info("shutting down")
	os.Exit(0)
}
