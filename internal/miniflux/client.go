package miniflux

import (
	api "miniflux.app/v2/client"
)

type CreateFeedOptions struct {
	FeedURL    string
	CategoryID int64
}

func CreateFeed(endpoint, apiKey string, opts CreateFeedOptions) (int64, error) {
	client := api.NewClient(endpoint, apiKey)

	req := &api.FeedCreationRequest{
		FeedURL:    opts.FeedURL,
		CategoryID: opts.CategoryID,
	}

	return client.CreateFeed(req)
}

type ListEntriesOptions struct {
	FeedID  int64
	Search  string
	Starred string
	Limit   int
	Status  string
	Offset  int
}

func Entries(endpoint, apiKey string, opts ListEntriesOptions) (*api.EntryResultSet, error) {
	client := api.NewClient(endpoint, apiKey)

	filter := &api.Filter{
		Search:    opts.Search,
		Limit:     opts.Limit,
		Offset:    opts.Offset,
		Order:     "published_at",
		Direction: "desc",
	}
	if opts.FeedID != 0 {
		filter.FeedID = opts.FeedID
	}
	if opts.Starred != "" {
		filter.Starred = opts.Starred
	}
	if opts.Status != "" {
		filter.Status = opts.Status
	}

	return client.Entries(filter)
}

func Feeds(endpoint, apiKey string) (api.Feeds, error) {
	client := api.NewClient(endpoint, apiKey)
	return client.Feeds()
}

func Validate(endpoint, apiKey string) error {
	client := api.NewClient(endpoint, apiKey)
	_, err := client.Me()
	return err
}
