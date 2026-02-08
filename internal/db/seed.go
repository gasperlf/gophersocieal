package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"ontopsolutions.net/gasperlf/social/internal/store"
)

var usernames = []string{
	"alice", "bob", "charlie", "david", "emma", "frank", "grace", "henry", "irene", "jack",
	"karen", "leo", "maria", "nathan", "olivia", "paul", "quinn", "rachel", "steve", "tina",
	"ursula", "victor", "wendy", "xavier", "yasmin", "zack", "adam", "bella", "carl", "diana",
	"ethan", "fiona", "george",
	"hannah34",
	"ian35",
	"julia36",
	"kevin37",
	"linda38",
	"mike39",
	"nina40",
	"oscar41",
	"penny42",
	"roger43",
	"sophia44",
	"tom45",
	"una46",
	"vincent47",
	"will48",
	"zoe49",
	"user50",
}

func Seed(store store.Storage, db *sql.DB) {
	ctx := context.Background()

	users := generateUsers(100)
	tx, _ := db.BeginTx(ctx, nil)

	for _, user := range users {
		if err := store.Users.Create(ctx, tx, user); err != nil {
			log.Panicln("failed to create users", err)
			return
		}
	}
	tx.Commit()

	posts := generatePosts(users, 500)

	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			_ = tx.Rollback()
			log.Panicln("failed to create posts", err)
			return
		}
	}

	comments := generateComments(users, posts, 1000)

	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Panicln("failed to create comments", err)
			return
		}
	}

	log.Println("database seeded successfully")
}

func generateUsers(count int) []*store.User {
	users := make([]*store.User, count)
	for i := 0; i < count; i++ {
		username := usernames[i%len(usernames)] + fmt.Sprintf("%d", i)
		users[i] = &store.User{
			Username: username,
			Email:    username + "@example.com",
			Role: store.Role{
				Name: "user", // default role for new users
			},
		}
	}
	return users
}

func generatePosts(users []*store.User, count int) []*store.Post {
	posts := make([]*store.Post, count)
	for i := 0; i < count; i++ {
		posts[i] = &store.Post{
			UserID:  users[i%len(users)].ID,
			Title:   fmt.Sprintf("Post Title %d", i+1),
			Content: fmt.Sprintf("This is post number %d", i+1),
			Tags:    []string{},
		}
	}

	return posts
}

func generateComments(users []*store.User, posts []*store.Post, count int) []*store.Comment {
	comments := make([]*store.Comment, count)
	for i := 0; i < count; i++ {
		comments[i] = &store.Comment{
			PostID:  posts[i%len(posts)].ID,
			UserID:  users[i%len(users)].ID,
			Content: fmt.Sprintf("This is comment number %d", i+1),
		}
	}
	return comments
}
