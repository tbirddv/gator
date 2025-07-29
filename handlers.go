package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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
