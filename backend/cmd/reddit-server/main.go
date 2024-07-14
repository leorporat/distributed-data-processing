// cmd/reddit_server/main.go

package main

import (
	"context"
	"encoding/json"
	"log"
	"net"

	"backend/internal/kafka"
	pb "backend/pkg/grpc"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedRedditServiceServer
	kafkaProducer *kafka.KafkaProducer
}

// SendRedditComment receives a Reddit comment via gRPC and sends it to Kafka
func (s *server) SendRedditComment(ctx context.Context, in *pb.RedditComment) (*pb.Response, error) {
	// Convert the comment to JSON
	commentJSON, err := json.Marshal(in)
	if err != nil {
		return &pb.Response{Success: false, Message: "Failed to marshal comment"}, err
	}

	// Send the JSON comment to Kafka
	err = s.kafkaProducer.Produce(string(commentJSON))
	if err != nil {
		return &pb.Response{Success: false, Message: "Failed to produce message to Kafka"}, err
	}

	return &pb.Response{Success: true, Message: "Comment sent to Kafka successfully"}, nil
}

func main() {
	// Create a new Kafka producer
	kafkaProducer, err := kafka.NewKafkaProducer("localhost:9092", "reddit_comments")
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer kafkaProducer.Close()

	// Set up gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterRedditServiceServer(s, &server{kafkaProducer: kafkaProducer})

	log.Println("gRPC server listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
