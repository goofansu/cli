package app

import (
	"fmt"

	"github.com/goofansu/cli/internal/format"
	"github.com/goofansu/cli/internal/linkding"
	"github.com/goofansu/cli/internal/miniflux"
)

type ListEntriesOptions struct {
	FeedID  int64
	Search  string
	Limit   int
	Status  string
	Starred string
	Offset  int
	JSON    string
	JQ      string
}

type ListBookmarksOptions struct {
	Query  string
	Limit  int
	Offset int
	JSON   string
	JQ     string
}

func (a *App) ListEntries(opts ListEntriesOptions) error {
	result, err := miniflux.Entries(a.Config.Miniflux.Endpoint, a.Config.Miniflux.APIKey, miniflux.ListEntriesOptions{
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

func (a *App) ListBookmarks(opts ListBookmarksOptions) error {
	result, err := linkding.ListBookmarks(a.Config.Linkding.Endpoint, a.Config.Linkding.APIKey, linkding.ListBookmarksOptions{
		Query:  opts.Query,
		Limit:  opts.Limit,
		Offset: opts.Offset,
	})
	if err != nil {
		return fmt.Errorf("failed to list bookmarks: %w", err)
	}

	output := map[string]any{
		"total": result.Count,
		"items": result.Results,
	}
	return format.Output(output, opts.JSON, opts.JQ)
}
