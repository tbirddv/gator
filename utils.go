package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/tbirddv/gator/internal/database"
	"github.com/tbirddv/gator/internal/rssfeed"
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

func NewNullTime(t time.Time) sql.NullTime {
	return sql.NullTime{
		Time:  t,
		Valid: true,
	}
}

func scrapeFeeds(s *state) error {
	feed, err := s.queries.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get next feed to fetch: %w", err)
	}

	// Simulate feed scraping
	fmt.Printf("Scraping feed: %s\n", feed.Url)

	fetchedTime := database.MarkFeedFetchedParams{
		LastFetchedAt: NewNullTime(time.Now()),
		ID:            feed.ID,
	}

	if err := s.queries.MarkFeedFetched(context.Background(), fetchedTime); err != nil {
		return fmt.Errorf("failed to mark feed as fetched: %w", err)
	}

	fetchedFeed, err := rssfeed.FetchRSSFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("failed to fetch RSS feed: %w", err)
	}

	for _, item := range fetchedFeed.Channel.Items {
		fmt.Printf("Item Title: %s\n", item.Title)
	}

	return nil
}
