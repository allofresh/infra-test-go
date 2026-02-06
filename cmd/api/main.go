package main

import (
	"log"
	"net/http"
	"os"

	"github.com/allofresh/infra-test-go/internal/handler"
	"github.com/allofresh/infra-test-go/internal/service"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	svc := service.NewProductService()
	h := handler.NewProductHandler(svc)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	log.Printf("Starting server on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
