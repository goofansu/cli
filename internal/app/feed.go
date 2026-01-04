package app

import (
	"fmt"

	"github.com/goofansu/cli/internal/format"
	"github.com/goofansu/cli/internal/miniflux"
)

type AddFeedOptions struct {
	URL        string
	CategoryID int64
}

type EntriesOptions struct {
	FeedID  int64
	Search  string
	Limit   int
	Status  string
	Starred string
	Offset  int
	JSON    string
	JQ      string
}

func (a *App) AddFeed(opts AddFeedOptions) error {
	categoryID := opts.CategoryID
	if categoryID == 0 {
		categoryID = 1
	}

	feedID, err := miniflux.CreateFeed(a.Config.Miniflux.Endpoint, a.Config.Miniflux.APIKey, miniflux.CreateFeedOptions{
		FeedURL:    opts.URL,
		CategoryID: categoryID,
	})
	if err != nil {
		return fmt.Errorf("failed to create feed: %w", err)
	}

	fmt.Printf("✓ Feed created successfully (ID: %d)\n", feedID)
	return nil
}

func (a *App) ListEntries(opts EntriesOptions) error {
	result, err := miniflux.Entries(a.Config.Miniflux.Endpoint, a.Config.Miniflux.APIKey, miniflux.EntriesOptions{
		FeedID:  opts.FeedID,
		Search:  opts.Search,
		Limit:   opts.Limit,
		Offset:  opts.Offset,
		Status:  opts.Status,
		Starred: opts.Starred,
	})
	if err != nil {
		return fmt.Errorf("failed to list entries: %w", err)
	}

	output := map[string]any{
		"total": result.Total,
		"items": result.Entries,
	}
	return format.Output(output, opts.JSON, opts.JQ)
}
