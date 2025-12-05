package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"service-3-gateway/internal/handler"

	articlepb "github.com/thatlq1812/service-2-article/proto"

	userpb "github.com/thatlq1812/service-1-user/proto"
)

func main() {
	// Get service addresses from environment
	userServiceAddr := getEnv("USER_SERVICE_ADDR", "localhost:50051")
	articleServiceAddr := getEnv("ARTICLE_SERVICE_ADDR", "localhost:50052")
	gatewayPort := getEnv("GATEWAY_PORT", "8080")

	log.Printf("Starting API Gateway...")
	log.Printf("User Service: %s", userServiceAddr)
	log.Printf("Article Service: %s", articleServiceAddr)

	// Connect to User Service (gRPC)
	log.Printf("Connecting to User Service...")
	userConn, err := grpc.Dial(
		userServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithTimeout(5*time.Second),
	)
	if err != nil {
		log.Fatalf("Failed to connect to User Service: %v", err)
	}
	defer userConn.Close()
	userClient := userpb.NewUserServiceClient(userConn)
	log.Printf("✓ Connected to User Service")

	// Connect to Article Service (gRPC)
	log.Printf("Connecting to Article Service...")
	articleConn, err := grpc.Dial(
		articleServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithTimeout(5*time.Second),
	)
	if err != nil {
		log.Fatalf("Failed to connect to Article Service: %v", err)
	}
	defer articleConn.Close()
	articleClient := articlepb.NewArticleServiceClient(articleConn)
	log.Printf("✓ Connected to Article Service")

	// Initialize handlers
	userHandler := handler.NewUserHandler(userClient)
	articleHandler := handler.NewArticleHandler(articleClient)

	// Create HTTP Router (receive REST request)
	router := mux.NewRouter()

	// Define REST endpoints
	router.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	router.HandleFunc("/articles", articleHandler.CreateArticle).Methods("POST")

	// Add logging middleware
	router.Use(loggingMiddleware)

	// Add CORS middleware for development
	router.Use(corsMiddleware)

	// API v1 routes
	api := router.PathPrefix("/api/v1").Subrouter()

	// User routes
	api.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	api.HandleFunc("/users", userHandler.ListUsers).Methods("GET")
	api.HandleFunc("/users/{id}", userHandler.GetUser).Methods("GET")
	api.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PUT")
	api.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")

	// Auth routes
	api.HandleFunc("/auth/login", userHandler.Login).Methods("POST")
	api.HandleFunc("/auth/refresh", userHandler.RefreshToken).Methods("POST")
	api.HandleFunc("/auth/validate", userHandler.ValidateToken).Methods("POST")
	api.HandleFunc("/auth/logout", userHandler.Logout).Methods("POST")

	// Article routes
	api.HandleFunc("/articles", articleHandler.CreateArticle).Methods("POST")
	api.HandleFunc("/articles", articleHandler.ListArticles).Methods("GET")
	api.HandleFunc("/articles/{id}", articleHandler.GetArticle).Methods("GET")
	api.HandleFunc("/articles/{id}", articleHandler.UpdateArticle).Methods("PUT")
	api.HandleFunc("/articles/{id}", articleHandler.DeleteArticle).Methods("DELETE")

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	}).Methods("GET")

	// Start server
	addr := ":" + gatewayPort
	log.Printf("API Gateway listening on %s", addr)
	log.Printf("Health check: http://localhost%s/health", addr)
	log.Printf("API Base URL: http://localhost%s/api/v1", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}

// loggingMiddleware logs all incoming requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("[%s] %s %s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
		log.Printf("[%s] %s completed in %v", r.Method, r.RequestURI, time.Since(start))
	})
}

// corsMiddleware adds CORS headers for development
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
