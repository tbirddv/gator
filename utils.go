package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/tbirddv/gator/internal/database"
)

func getLoggedInUser(s *state) (database.User, error) {
	if s.config.CurrentUserName == "" {
		return database.User{}, errors.New("no user is currently logged in")
	}
	user, err := s.queries.GetUserByName(context.Background(), s.config.CurrentUserName)
	if err != nil {
		return database.User{}, fmt.Errorf("failed to get current user: %w", err)
	}
	return user, nil
}
