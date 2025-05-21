package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"backend/internal/reddit_implementation"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := reddit_implementation.NewRedditServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	for i := 0; i < 100; i++ {
		_, err := c.GetPostComments(ctx, &reddit_implementation.PostRequest{
			Subreddit:   "uiuc",
			SearchQuery: "Arboretum",
			Limit:       1,
		})
		if err != nil {
			fmt.Print("error")
		}
		// fmt.Printf(r.PostTitle)
		// fmt.Printf(r.Comments)
	}
	// if err != nil {
	// 	log.Fatalf("could not get comments: %v", err)
	// }

	// Print out results from request
	// fmt.Printf("Post Title: %s\n", r.PostTitle)
	// fmt.Printf("Number of comments: %d\n", len(r.Comments))
	// for i, comment := range r.Comments {
	// 	fmt.Printf("Comment %d: %s\n", i+1, comment.Body)
	// }
}
