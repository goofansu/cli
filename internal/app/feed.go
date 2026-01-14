package app

import (
	"fmt"

	"github.com/goofansu/mlwcli/internal/format"
	"github.com/goofansu/mlwcli/internal/miniflux"
)

type AddFeedOptions struct {
	URL        string
	CategoryID int64
}

type ListFeedsOptions struct {
	JSON string
	JQ   string
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

func (a *App) ListFeeds(opts ListFeedsOptions) error {
	feeds, err := miniflux.Feeds(a.Config.Miniflux.Endpoint, a.Config.Miniflux.APIKey)
	if err != nil {
		return fmt.Errorf("failed to list feeds: %w", err)
	}

	data := map[string]any{
		"total": len(feeds),
		"items": feeds,
	}

	return format.Output(data, opts.JSON, opts.JQ)
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

func (a *App) SaveEntry(entryID int64) error {
	err := miniflux.SaveEntry(a.Config.Miniflux.Endpoint, a.Config.Miniflux.APIKey, entryID)
	if err != nil {
		return fmt.Errorf("failed to save entry: %w", err)
	}
	fmt.Printf("Entry %d saved successfully\n", entryID)
	return nil
}
