package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"

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

func parseFlexibleTimestamp(timestampStr string) (time.Time, error) {
	formats := []string{
		time.RFC3339,                // "2006-01-02T15:04:05Z07:00"
		time.RFC1123Z,               // "Mon, 02 Jan 2006 15:04:05 -0700"
		time.RFC1123,                // "Mon, 02 Jan 2006 15:04:05 MST"
		"2006-01-02 15:04:05-07:00", // PostgreSQL format
		"2006-01-02 15:04:05",       // Without timezone
		"2006-01-02T15:04:05",       // ISO without timezone
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timestampStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse timestamp: %s", timestampStr)
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
		pubDate, err := parseFlexibleTimestamp(item.PubDate)
		if err != nil {
			fmt.Printf("Skipping item with invalid pubDate: %s\n", item.PubDate)
			continue
		}
		postParams := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: item.Description != ""},
			PublishedAt: pubDate,
			FeedID:      feed.ID,
		}
		err = s.queries.CreatePost(context.Background(), postParams)
		if err != nil {
			if err, ok := err.(*pq.Error); ok && err.Code == "23505" {
				continue // Skip duplicate posts
			}
			fmt.Printf("Error creating post: %v\n", err)
		}
	}

	return nil
}
