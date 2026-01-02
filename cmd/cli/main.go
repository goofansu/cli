package main

import (
	"fmt"
	"os"

	flags "github.com/jessevdk/go-flags"

	"github.com/goofansu/cli/internal/app"
	"github.com/goofansu/cli/internal/auth"
	"github.com/goofansu/cli/internal/config"
)

type Options struct {
	Login  LoginCommand  `command:"login" description:"Authenticate with a service"`
	Logout LogoutCommand `command:"logout" description:"Remove credentials for a service"`
	Add    AddCommand    `command:"add" description:"Add a resource (feed or bookmark)"`
	List   ListCommand   `command:"list" description:"List resources (entries or bookmarks)"`
}

type AddCommand struct {
	BaseCommand
	Feed     AddFeedCommand     `command:"feed" description:"Add a feed (miniflux)"`
	Bookmark AddBookmarkCommand `command:"bookmark" description:"Add a bookmark (linkding)"`
}

type ListCommand struct {
	BaseCommand
	Entries   ListEntriesCommand   `command:"entries" description:"List entries (miniflux)"`
	Bookmarks ListBookmarksCommand `command:"bookmarks" description:"List bookmarks (linkding)"`
}

type BaseCommand struct {
	App *app.App
}

type JSONOutputOptions struct {
	JSON string `long:"json" value-name:"fields" description:"Output JSON with the specified fields (comma-separated)"`
	JQ   string `long:"jq" value-name:"expression" description:"Filter JSON output using a jq expression (requires --json)"`
}

type LoginCommand struct {
	BaseCommand
	Args struct {
		Service string `positional-arg-name:"service" description:"Service name (miniflux or linkding)" required:"yes"`
	} `positional-args:"yes"`
	Endpoint string `long:"endpoint" description:"Service endpoint URL" required:"yes"`
	APIKey   string `long:"api-key" description:"API key" required:"yes"`
}

type LogoutCommand struct {
	BaseCommand
	Args struct {
		Service string `positional-arg-name:"service" description:"Service name (miniflux or linkding)" required:"yes"`
	} `positional-args:"yes"`
}

type AddFeedCommand struct {
	BaseCommand
	Args struct {
		URL string `positional-arg-name:"url" description:"URL of the feed to subscribe to" required:"yes"`
	} `positional-args:"yes"`
	CategoryID int64 `long:"category-id" description:"Miniflux category ID (defaults to 1)"`
}

type AddBookmarkCommand struct {
	BaseCommand
	Args struct {
		URL string `positional-arg-name:"url" description:"URL of the bookmark to add" required:"yes"`
	} `positional-args:"yes"`
	Notes string `long:"notes" description:"Optional notes for the bookmark"`
	Tags  string `long:"tags" description:"Optional tags separated by spaces"`
}

type ListEntriesCommand struct {
	BaseCommand
	JSONOutputOptions
	Limit   int    `long:"limit" description:"Maximum number of results" default:"10"`
	Offset  int    `long:"offset" description:"Number of results to skip" default:"0"`
	Search  string `long:"search" description:"Search query text"`
	Status  string `long:"status" value-name:"status" description:"Filter by entry status (read, unread, removed)" default:"unread"`
	Starred bool   `long:"starred" description:"Filter by starred entries"`
	FeedID  int64  `long:"feed-id" description:"Filter by feed ID"`
}

type ListBookmarksCommand struct {
	BaseCommand
	JSONOutputOptions
	Limit  int    `long:"limit" description:"Maximum number of results" default:"10"`
	Offset int    `long:"offset" description:"Number of results to skip" default:"0"`
	Search string `long:"search" description:"Search query text"`
}

func (c *LoginCommand) Execute(_ []string) error {
	return auth.Login(c.Args.Service, c.Endpoint, c.APIKey)
}

func (c *LogoutCommand) Execute(_ []string) error {
	return auth.Logout(c.Args.Service)
}

func (c *AddFeedCommand) Execute(_ []string) error {
	opts := app.AddFeedOptions{
		URL:        c.Args.URL,
		CategoryID: c.CategoryID,
	}

	return c.App.AddFeed(opts)
}

func (c *AddBookmarkCommand) Execute(_ []string) error {
	opts := app.AddBookmarkOptions{
		URL:   c.Args.URL,
		Notes: c.Notes,
		Tags:  c.Tags,
	}

	return c.App.AddBookmark(opts)
}

func (c *ListEntriesCommand) Execute(_ []string) error {
	starred := ""
	if c.Starred {
		starred = "1"
	}

	opts := app.ListEntriesOptions{
		FeedID:  c.FeedID,
		Search:  c.Search,
		Limit:   c.Limit,
		Offset:  c.Offset,
		Status:  c.Status,
		Starred: starred,
		JSON:    c.JSON,
		JQ:      c.JQ,
	}

	return c.App.ListEntries(opts)
}

func (c *ListBookmarksCommand) Execute(_ []string) error {
	opts := app.ListBookmarksOptions{
		Query:  c.Search,
		Limit:  c.Limit,
		Offset: c.Offset,
		JSON:   c.JSON,
		JQ:     c.JQ,
	}
	return c.App.ListBookmarks(opts)
}

func (c *LoginCommand) Usage() string {
	return "<service> [OPTIONS]"
}

func (c *LogoutCommand) Usage() string {
	return "<service>"
}

func (c *AddFeedCommand) Usage() string {
	return "<url>"
}

func (c *AddBookmarkCommand) Usage() string {
	return "<url>"
}

func (c *ListEntriesCommand) Usage() string {
	return "[OPTIONS]"
}

func (c *ListBookmarksCommand) Usage() string {
	return "[OPTIONS]"
}

func main() {
	cfg, err := config.Load()
	if err != nil && !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "warning: failed to load config: %v\n", err)
	}
	if cfg == nil {
		cfg = &config.Config{}
	}

	application := app.New(cfg)

	opts := Options{}
	opts.Login.App = application
	opts.Logout.App = application
	opts.Add.Feed.App = application
	opts.Add.Bookmark.App = application
	opts.List.Entries.App = application
	opts.List.Bookmarks.App = application

	parser := flags.NewParser(&opts, flags.HelpFlag|flags.PassDoubleDash)
	parser.ShortDescription = "My command-line tool for agents"
	parser.LongDescription = "Manage bookmarks and RSS feeds from terminal.\n\nExamples:\ncli login linkding --endpoint https://linkding.example.com --api-key YOUR_API_KEY\ncli login miniflux --endpoint https://miniflux.example.com --api-key YOUR_API_KEY\ncli add bookmark https://example.com --tags \"cool useful\"\ncli list bookmarks\ncli add feed https://blog.example.com/feed.xml\ncli list entries"

	if len(os.Args) == 1 {
		parser.WriteHelp(os.Stdout)
		return
	}

	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok {
			if flagsErr.Type == flags.ErrHelp {
				fmt.Fprint(os.Stdout, flagsErr.Message)
				return
			}
			fmt.Fprintf(os.Stderr, "error: %s\n\n", flagsErr.Message)
			parser.WriteHelp(os.Stderr)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}
