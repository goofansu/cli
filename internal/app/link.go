package app

import (
	"fmt"
	"strings"

	"github.com/goofansu/mlwcli/internal/format"
	"github.com/goofansu/mlwcli/internal/linkding"
)

type AddLinkOptions struct {
	URL   string
	Notes string
	Tags  string
}

type ListLinksOptions struct {
	Query  string
	Limit  int
	Offset int
	JSON   string
	JQ     string
}

func (a *App) AddLink(opts AddLinkOptions) error {
	tagNames := []string{}
	if opts.Tags != "" {
		tagNames = strings.Split(opts.Tags, " ")
	}

	_, err := linkding.CreateBookmark(a.Config.Linkding.Endpoint, a.Config.Linkding.APIKey, linkding.CreateBookmarkOptions{
		URL:      opts.URL,
		Notes:    opts.Notes,
		TagNames: tagNames,
	})
	if err != nil {
		return fmt.Errorf("failed to create link: %w", err)
	}

	fmt.Printf("✓ Link created successfully\n")
	return nil
}

func (a *App) ListLinks(opts ListLinksOptions) error {
	result, err := linkding.ListBookmarks(a.Config.Linkding.Endpoint, a.Config.Linkding.APIKey, linkding.ListBookmarksOptions{
		Query:  opts.Query,
		Limit:  opts.Limit,
		Offset: opts.Offset,
	})
	if err != nil {
		return fmt.Errorf("failed to list links: %w", err)
	}

	data := map[string]any{
		"total": result.Count,
		"items": result.Results,
	}

	return format.Output(data, opts.JSON, opts.JQ)
}
