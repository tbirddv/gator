package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/tbirddv/gator/internal/database"
	"github.com/tbirddv/gator/internal/rssfeed"
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
	var url string
	if len(s.args) < 1 {
		url = "https://www.wagslane.dev/index.xml" // Default RSS feed URL
	} else {
		url = s.args[0]
	}

	ctx := context.Background()
	feed, err := rssfeed.FetchRSSFeed(ctx, url)
	if err != nil {
		return fmt.Errorf("failed to fetch RSS feed: %w", err)
	}

	feed.UnescapeTitleandDescription()

	fmt.Printf("Title: %s\n", feed.Channel.Title)
	fmt.Printf("Link: %s\n", feed.Channel.Link)
	fmt.Printf("Description: %s\n", feed.Channel.Description)

	for _, item := range feed.Channel.Items {
		fmt.Printf("\nItem Title: %s\n", item.Title)
		fmt.Printf("Item Link: %s\n", item.Link)
		fmt.Printf("Item Description: %s\n", item.Description)
		fmt.Printf("Published Date: %s\n", item.PubDate)
	}

	return nil
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

	user, err := s.queries.GetUserByName(context.Background(), s.config.CurrentUserName)
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

	fmt.Printf("Feed created successfully: %s (%s)\n", newFeed.Name, newFeed.Url)
	fmt.Printf("Feed ID: %s\n", newFeed.ID)
	fmt.Printf("Created At: %s\n", newFeed.CreatedAt)
	fmt.Printf("Updated At: %s\n", newFeed.UpdatedAt)
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
