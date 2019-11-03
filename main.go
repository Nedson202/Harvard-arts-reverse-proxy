package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nedson202/harvard-arts-reverse-proxy/reverse_proxy"
	"github.com/rs/cors"

	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"
)

func main() {
	err := gotenv.Load()
	if err != nil {
		log.Println(err)
	}

	router := mux.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
	})

	port := os.Getenv("PORT")
	redisHost := os.Getenv("REDIS_URL")
	baseURL := os.Getenv("HARVARD_API_BASE_URL")
	harvardAPIKey := os.Getenv("HARVARD_ARTS_API_KEY")
	elasticSearchURL := os.Getenv("BONSAI_URL")

	_, err = reverse_proxy.New(redisHost, router, baseURL, harvardAPIKey, elasticSearchURL)
	if err != nil {
		log.Println(err)
	}

	combineServerAddress := fmt.Sprintf("%s%s", ":", port)

	server := &http.Server{
		// launch server with CORS validations
		Handler:      c.Handler(router),
		Addr:         combineServerAddress,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	// Start Server
	func() {
		startMessage := fmt.Sprintf("%s%s", "Server started on http://localhost:", port)
		log.Println(startMessage)

		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	handleShutdown(server)
}

// Handle graceful shutdown
func handleShutdown(server *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-interruptChan

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	server.Shutdown(ctx)

	log.Println("Shutting down server")
	os.Exit(0)
}
