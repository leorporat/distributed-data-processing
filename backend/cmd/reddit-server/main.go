package main

// THIS FILE RUNS THE GRPC SERVER WITH GRPC-WEB SUPPORT AND JSON FALLBACK

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"backend/internal/reddit_implementation"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/vartanbeno/go-reddit/v2/reddit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

const (
	// Port the server will listen on
	port = "8080"
)

// PostRequest JSON equivalent for easier frontend-backend communication
type PostRequestJSON struct {
	Subreddit    string `json:"subreddit"`
	SearchQuery  string `json:"searchQuery"`
	Limit        int32  `json:"limit"`
}

// CommentsResponse JSON equivalent
type CommentsResponseJSON struct {
	PostTitle string        `json:"postTitle"`
	Comments  []CommentJSON `json:"comments"`
}

// Comment JSON equivalent
type CommentJSON struct {
	Body string `json:"body"`
}

func main() {
	// Set up logging with timestamps
	logger := log.New(os.Stdout, "[GRPC-WEB] ", log.LstdFlags|log.Lmicroseconds)
	grpclog.SetLoggerV2(grpclog.NewLoggerV2(logger.Writer(), logger.Writer(), logger.Writer()))

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// Check if Reddit credentials are set
	if os.Getenv("redditID") == "" || os.Getenv("redditSecret") == "" {
		return fmt.Errorf("Reddit API credentials not set. Please set environment variables: redditID, redditSecret, redditUser, redditPassword")
	}

	// Set up Reddit credentials
	credentials := reddit.Credentials{
		ID:       os.Getenv("redditID"),
		Secret:   os.Getenv("redditSecret"),
		Username: os.Getenv("redditUser"),
		Password: os.Getenv("redditPassword"),
	}

	// Create the Reddit service server
	server, err := reddit_implementation.NewServer(credentials)
	if err != nil {
		return fmt.Errorf("failed to create Reddit server: %v", err)
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Register our Reddit service
	reddit_implementation.RegisterRedditServiceServer(grpcServer, server)

	// Create the gRPC-Web wrapper for the gRPC server
	wrappedGrpc := grpcweb.WrapServer(grpcServer, 
		grpcweb.WithOriginFunc(func(origin string) bool {
			// Allow all origins for development
			log.Printf("Origin check: %s", origin)
			return true
		}),
		grpcweb.WithWebsockets(true),
		grpcweb.WithWebsocketOriginFunc(func(req *http.Request) bool {
			return true
		}),
		grpcweb.WithAllowNonRootResource(true),
		grpcweb.WithCorsForRegisteredEndpointsOnly(false),
	)

	// Create an HTTP handler that intercepts gRPC-Web requests
	handler := http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		startTime := time.Now()
		
		// Log request details
		log.Printf("Request: %s %s", req.Method, req.URL.Path)
		log.Printf("Headers: %v", req.Header)

		// Add CORS headers for browser support
		resp.Header().Set("Access-Control-Allow-Origin", "*")
		resp.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		resp.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-User-Agent, X-Grpc-Web, x-grpc-web")
		resp.Header().Set("Access-Control-Max-Age", "86400")

		// Handle preflight OPTIONS requests
		if req.Method == "OPTIONS" {
			log.Printf("Handling OPTIONS request")
			resp.WriteHeader(http.StatusOK)
			return
		}

		// Check if this is a JSON request 
		isJSON := strings.Contains(req.Header.Get("Content-Type"), "json") || 
			       strings.Contains(req.Header.Get("Accept"), "json")

		// Handle Reddit API JSON fallback
		if isJSON && strings.Contains(req.URL.Path, "/reddit_implementation.RedditService/GetPostComments") {
			log.Printf("Handling as JSON request")
			handleJSONRequest(resp, req, server)
			log.Printf("Request completed in %v", time.Since(startTime))
			return
		}

		// Check if this is a gRPC-Web request
		if wrappedGrpc.IsGrpcWebRequest(req) || wrappedGrpc.IsAcceptableGrpcCorsRequest(req) {
			log.Printf("Handling as gRPC-Web request")
			wrappedGrpc.ServeHTTP(resp, req)
			log.Printf("Request completed in %v", time.Since(startTime))
			return
		}

		// If not a gRPC-Web request but has the right path, show error
		if strings.Contains(req.URL.Path, "/reddit_implementation.RedditService/") {
			log.Printf("Received gRPC request, but not properly formatted as gRPC-Web")
			log.Printf("Content-Type: %s", req.Header.Get("Content-Type"))
			
			resp.WriteHeader(http.StatusBadRequest)
			resp.Write([]byte("This endpoint requires a gRPC-Web client or JSON format"))
			return
		}

		// Default response for other requests
		resp.WriteHeader(http.StatusOK)
		resp.Write([]byte("Reddit Summarizer gRPC-Web Server. Use a gRPC-Web client or JSON format to communicate."))
		log.Printf("Request completed in %v", time.Since(startTime))
	})

	// Start the HTTP server
	httpServer := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	log.Printf("Starting gRPC-Web server with JSON fallback on port %s", port)
	if err := httpServer.ListenAndServe(); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}

// handleJSONRequest processes the Reddit API request in JSON format
func handleJSONRequest(resp http.ResponseWriter, req *http.Request, server *reddit_implementation.Server) {
	// Set response headers
	resp.Header().Set("Content-Type", "application/json")
	
	// Read request body
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		resp.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(resp).Encode(map[string]string{"error": "Could not read request body"})
		return
	}
	
	// Parse JSON request
	var jsonReq PostRequestJSON
	if err := json.Unmarshal(body, &jsonReq); err != nil {
		log.Printf("Error parsing JSON request: %v", err)
		resp.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(resp).Encode(map[string]string{"error": "Invalid JSON format"})
		return
	}
	
	log.Printf("JSON Request data: %+v", jsonReq)
	
	// Create protobuf request
	protoReq := &reddit_implementation.PostRequest{
		Subreddit:   jsonReq.Subreddit,
		SearchQuery: jsonReq.SearchQuery,
		Limit:       jsonReq.Limit,
	}
	
	// Process request using the gRPC service
	protoResp, err := server.GetPostComments(req.Context(), protoReq)
	if err != nil {
		log.Printf("Error processing request: %v", err)
		resp.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(resp).Encode(map[string]string{"error": err.Error()})
		return
	}
	
	// Convert protobuf response to JSON
	jsonResp := CommentsResponseJSON{
		PostTitle: protoResp.PostTitle,
		Comments:  make([]CommentJSON, len(protoResp.Comments)),
	}
	
	for i, c := range protoResp.Comments {
		jsonResp.Comments[i] = CommentJSON{Body: c.Body}
	}
	
	log.Printf("Response prepared with %d comments for post: %s", len(jsonResp.Comments), jsonResp.PostTitle)
	
	// Write JSON response
	resp.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(resp).Encode(jsonResp); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}
