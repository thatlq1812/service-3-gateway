package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/thatlq1812/service-3-gateway/internal/circuit"
	"github.com/thatlq1812/service-3-gateway/internal/handler"
	"github.com/thatlq1812/service-3-gateway/internal/middleware"

	articlepb "github.com/thatlq1812/service-2-article/proto"

	userpb "github.com/thatlq1812/service-1-user/proto"
)

// connectWithRetry attempts to establish gRPC connection with exponential backoff
func connectWithRetry(address string, serviceName string, maxRetries int) (*grpc.ClientConn, error) {
	backoff := 1 * time.Second
	maxBackoff := 30 * time.Second

	for i := 0; i < maxRetries; i++ {
		log.Printf("[%s] Connection attempt %d/%d to %s", serviceName, i+1, maxRetries, address)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		conn, err := grpc.DialContext(
			ctx,
			address,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(), // Block until connected or timeout
		)
		cancel()

		if err == nil {
			log.Printf("[%s] Successfully connected to %s", serviceName, address)
			return conn, nil
		}

		log.Printf("[%s] Connection attempt %d failed: %v", serviceName, i+1, err)

		if i < maxRetries-1 {
			log.Printf("[%s] Retrying in %v...", serviceName, backoff)
			time.Sleep(backoff)

			// Exponential backoff
			backoff *= 2
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}
	}

	return nil, fmt.Errorf("failed to connect after %d retries", maxRetries)
}

func main() {
	// Get service addresses from environment
	userServiceAddr := getEnv("USER_SERVICE_ADDR", "localhost:50051")
	articleServiceAddr := getEnv("ARTICLE_SERVICE_ADDR", "localhost:50052")
	gatewayPort := getEnv("GATEWAY_PORT", "8080")

	log.Printf("Starting API Gateway...")
	log.Printf("User Service: %s", userServiceAddr)
	log.Printf("Article Service: %s", articleServiceAddr)

	// Connect to User Service (gRPC) with retry logic
	log.Printf("Connecting to User Service...")
	userConn, err := connectWithRetry(userServiceAddr, "User Service", 5)
	if err != nil {
		log.Fatalf("Failed to connect to User Service after retries: %v", err)
	}
	defer userConn.Close()
	userClient := userpb.NewUserServiceClient(userConn)
	log.Printf("✓ Connected to User Service")

	// Connect to Article Service (gRPC) with retry logic
	log.Printf("Connecting to Article Service...")
	articleConn, err := connectWithRetry(articleServiceAddr, "Article Service", 5)
	if err != nil {
		log.Fatalf("Failed to connect to Article Service after retries: %v", err)
	}
	defer articleConn.Close()
	articleClient := articlepb.NewArticleServiceClient(articleConn)
	log.Printf("✓ Connected to Article Service")

	// Initialize circuit breakers for each service
	// maxFailures: 5 consecutive failures trigger circuit open
	// resetTimeout: 30s before attempting half-open state
	userCircuit := circuit.NewBreaker(5, 30*time.Second)
	articleCircuit := circuit.NewBreaker(5, 30*time.Second)

	log.Printf("Circuit Breakers initialized")
	log.Printf("- User Service: max_failures=5, reset_timeout=30s")
	log.Printf("- Article Service: max_failures=5, reset_timeout=30s")

	// Initialize handlers with circuit breakers
	userHandler := handler.NewUserHandlerWithCircuit(userClient, userCircuit)
	articleHandler := handler.NewArticleHandlerWithCircuit(articleClient, articleCircuit)

	// Create HTTP Router (receive REST request)
	router := mux.NewRouter()

	// Define REST endpoints
	router.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	router.HandleFunc("/articles", articleHandler.CreateArticle).Methods("POST")

	// Add global timeout middleware (5 seconds per request)
	router.Use(middleware.TimeoutMiddleware(5 * time.Second))

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

	// Health check with backend service status
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Check User Service connection
		userState := userConn.GetState().String()
		userHealthy := userState == "READY"

		// Check Article Service connection
		articleState := articleConn.GetState().String()
		articleHealthy := articleState == "READY"

		// Overall health
		overallHealthy := userHealthy && articleHealthy
		statusCode := http.StatusOK
		status := "healthy"

		if !overallHealthy {
			statusCode = http.StatusServiceUnavailable
			status = "degraded"
		}

		response := fmt.Sprintf(`{
	"status": "%s",
	"services": {
		"user_service": {
			"status": "%s",
			"healthy": %v
		},
		"article_service": {
			"status": "%s",
			"healthy": %v
		}
	}
}`, status, userState, userHealthy, articleState, articleHealthy)

		w.WriteHeader(statusCode)
		w.Write([]byte(response))
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
