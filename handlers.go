package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/tbirddv/gator/internal/database"
)

func HandleLogin(s *state) error {
	if len(s.args) < 1 {
		return errors.New("username is required")
	}
	username := s.args[0]
	if _, err := s.queries.GetUserByName(context.Background(), username); errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("user %s does not exist", username)
	}
	return s.config.SetUser(username)
}

func HandleHelp(commands map[string]Command, s *state) error {
	if len(s.args) < 1 {
		for _, command := range commands {
			println(command.Name + ": " + command.Description)
		}
		return nil
	}
	commandName := s.args[0]
	if command, exists := commands[commandName]; exists {
		println(command.Name + ": " + command.Description)
		return nil
	}
	return errors.New("unknown command: " + commandName)
}

func HandleRegister(s *state) error {
	if len(s.args) < 1 {
		return errors.New("username is required")
	}
	username := s.args[0]
	if username == "" {
		return errors.New("username cannot be empty")
	}

	_, err := s.queries.GetUserByName(context.Background(), username)
	if err == nil {
		return errors.New("user already exists")
	}

	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	}

	user, err := s.queries.CreateUser(context.Background(), userParams)
	if err != nil {
		return err
	}

	err = s.config.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Println("User registered:", user.Name)
	fmt.Println(user.ID, user.CreatedAt, user.UpdatedAt, user.Name)

	return nil
}

func HandleResetUsers(s *state) error {
	err := s.queries.ResetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to reset users: %w", err)
	}
	err = s.config.ClearUser()
	if err != nil {
		return fmt.Errorf("failed to clear current user: %w", err)
	}
	fmt.Println("All users have been reset.")
	return nil
}

func HandleGetUsers(s *state) error {
	users, err := s.queries.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get users: %w", err)
	}
	if len(users) == 0 {
		fmt.Println("No users found.")
		return nil
	}
	for _, user := range users {
		if user.Name == s.config.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}

func HandleAgg(s *state) error {
	var timeBetweenRequests time.Duration
	if len(s.args) < 1 {
		return errors.New("time between requests is required")
	} else {
		var err error
		timeBetweenRequests, err = time.ParseDuration(s.args[0])
		if err != nil {
			return fmt.Errorf("invalid duration: %w", err)
		}
	}

	ticker := time.NewTicker(timeBetweenRequests)
	defer ticker.Stop()
	for ; ; <-ticker.C {
		err := scrapeFeeds(s)
		if err != nil {
			return fmt.Errorf("error scraping feeds: %v", err)
		}
	}
}

func HandleCreateFeed(s *state) error {
	if len(s.args) < 2 {
		return errors.New("feed name and URL are required")
	}
	name := s.args[0]
	url := s.args[1]

	if name == "" || url == "" {
		return errors.New("feed name and URL cannot be empty")
	}

	user, err := getLoggedInUser(s)
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}

	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	}

	newFeed, err := s.queries.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return fmt.Errorf("failed to create feed: %w", err)
	}
	s.args[0] = s.args[1] // change args for expected input of HandleFollow

	fmt.Printf("Feed created successfully: %s (%s)\n", newFeed.Name, newFeed.Url)
	fmt.Printf("Feed ID: %s\n", newFeed.ID)
	fmt.Printf("Created At: %s\n", newFeed.CreatedAt)
	fmt.Printf("Updated At: %s\n", newFeed.UpdatedAt)
	if err := HandleFollow(s); err != nil {
		return fmt.Errorf("failed to follow feed after creation: %w", err)
	}
	return nil
}

func HandleGetFeeds(s *state) error {
	feeds, err := s.queries.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get feeds: %w", err)
	}
	if len(feeds) == 0 {
		fmt.Println("No feeds found.")
		return nil
	}
	for _, feed := range feeds {
		fmt.Printf("Feed Name: %s\n", feed.Name)
		fmt.Printf("Feed URL: %s\n", feed.Url)
		fmt.Printf("User: %s\n", feed.UserName)
		fmt.Println("-----------------------------")
	}
	return nil
}

func HandleFollow(s *state) error {
	if len(s.args) < 1 {
		return errors.New("feed URL is required")
	}
	url := s.args[0]

	// Check if the feed already exists
	feed, err := s.queries.GetFeedByURL(context.Background(), url)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("failed to check feed existence: %w", err)
	}
	if errors.Is(err, sql.ErrNoRows) {
		// Feed does not exist, but creating a new feed has a different usage pattern, inform the user
		return fmt.Errorf("feed with URL %s does not exist. Use 'addfeed' command to add a new feed.\n Usage: addfeed <name> <url>", url)
	}
	user, err := getLoggedInUser(s)
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}
	followParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	follow, err := s.queries.CreateFeedFollow(context.Background(), followParams)
	if err != nil {
		return fmt.Errorf("failed to create feed follow: %w", err)
	}
	fmt.Printf("User %s is now following feed %s\n", follow.UserName, follow.FeedName)
	return nil
}

func HandleGetFollows(s *state) error {

	user, err := getLoggedInUser(s)
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}

	follows, err := s.queries.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("failed to get followed feeds for user %s: %w", user.Name, err)
	}

	if len(follows) == 0 {
		fmt.Printf("User %s is not following any feeds.\n", user.Name)
		return nil
	}

	fmt.Printf("Feeds followed by %s:\n", user.Name)
	for _, follow := range follows {
		fmt.Printf("- %s\n", follow.FeedName)
	}
	return nil
}

func HandleUnfollow(s *state) error {
	if len(s.args) < 1 {
		return errors.New("feed URL is required")
	}
	url := s.args[0]

	// Check if the feed exists
	feed, err := s.queries.GetFeedByURL(context.Background(), url)
	if err != nil {
		return fmt.Errorf("failed to get feed by URL: %w", err)
	}

	user, err := getLoggedInUser(s)
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}

	unfollowParams := database.UnfollowFeedParams{
		FeedID: feed.ID,
		UserID: user.ID,
	}
	// Attempt to delete the feed follow

	err = s.queries.UnfollowFeed(context.Background(), unfollowParams)
	if err != nil {
		return fmt.Errorf("failed to unfollow feed: %w", err)
	}

	fmt.Printf("User %s has unfollowed feed %s\n", user.Name, feed.Name)
	return nil
}

func HandleBrowse(s *state) error {
	var limit int32 = 2
	if len(s.args) >= 1 {
		input, err := strconv.Atoi(s.args[0])
		if err != nil {
			return fmt.Errorf("invalid limit value: %w", err)
		}
		limit = int32(input)
	}
	user, err := getLoggedInUser(s)
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}
	browseParams := database.GetPostsForUserParams{
		ID:    user.ID,
		Limit: limit,
	}
	posts, err := s.queries.GetPostsForUser(context.Background(), browseParams)
	if err != nil {
		return fmt.Errorf("failed to get posts for user %s: %w", user.Name, err)
	}
	if len(posts) == 0 {
		fmt.Printf("User %s is has no posts from followed feeds.\n", user.Name)
		return nil
	}
	for _, post := range posts {
		fmt.Printf("Post Title: %s\n", post.Title)
		fmt.Printf("Post URL: %s\n", post.Url)
		fmt.Printf("Published At: %s\n", post.PublishedAt)
		fmt.Println("-----------------------------")
	}
	return nil
}
