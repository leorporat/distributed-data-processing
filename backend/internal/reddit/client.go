package reddit

import (
	"context"
	"fmt"
	"log"

	"github.com/vartanbeno/go-reddit/v2/reddit"
)

var ctx = context.Background()

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() (err error) {
	// Sets credentials for reddit API
	credentials := reddit.Credentials{ID: "client-id", Secret: "client-secret", Username: "username", Password: "password"}
	client, _ := reddit.NewClient(credentials)

	// Calls the Reddit API to search for posts. First string: search query, second string: subreddit
	posts, _, err := client.Subreddit.SearchPosts(ctx, "Best cafe for studying", "uiuc", &reddit.ListPostSearchOptions{
		ListPostOptions: reddit.ListPostOptions{
			ListOptions: reddit.ListOptions{
				Limit: 1,
			},
			Time: "all",
		},
		Sort: "relevance",
	})
	if err != nil {
		fmt.Println("There was an error")
		return
	}

	// Loops over queried posts and prints out the comments of each one
	for _, post := range posts {
		postandcomments, _, _ := client.Post.Get(ctx, post.ID)
		fmt.Print("TITLE: ")
		fmt.Print(post.Title)
		fmt.Println()
		for _, comment := range postandcomments.Comments {
			fmt.Println(comment.Body)
		}
	}

	return
}
