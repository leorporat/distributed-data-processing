package reddit_implementation

import (
	"context"
	"fmt"

	"github.com/vartanbeno/go-reddit/v2/reddit"
)

type Server struct {
	UnimplementedRedditServiceServer
	redditClient *reddit.Client
}

func NewServer(credentials reddit.Credentials) (*Server, error) {
	redditClient, err := reddit.NewClient(credentials)
	if err != nil {
		return nil, fmt.Errorf("error creating Reddit client: %v", err)
	}
	return &Server{redditClient: redditClient}, nil
}

func (s *Server) GetPostComments(ctx context.Context, req *PostRequest) (*CommentsResponse, error) {

	// Calls the Reddit API to search for posts. First string: search query, second string: subreddit
	posts, _, err := s.redditClient.Subreddit.SearchPosts(ctx, req.SearchQuery, req.Subreddit, &reddit.ListPostSearchOptions{
		ListPostOptions: reddit.ListPostOptions{
			ListOptions: reddit.ListOptions{
				Limit: int(req.Limit),
			},
			Time: "all",
		},
		Sort: "relevance",
	})

	// Check for errors
	if err != nil {
		return nil, fmt.Errorf("error searching posts %v", err)
	}

	if len(posts) == 0 {
		return nil, fmt.Errorf("no posts found")
	}

	post := posts[0]
	postandcomments, _, err := s.redditClient.Post.Get(ctx, post.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting comments %v", err)
	}

	// Populate comments array and create response
	updatedComments := make([]*Comment, len(postandcomments.Comments))
	for i, comment := range postandcomments.Comments {
		updatedComments[i] = &Comment{Body: comment.Body}
	}
	response := &CommentsResponse{
		PostTitle: post.Title,
		Comments:  updatedComments,
	}

	return response, nil
}
