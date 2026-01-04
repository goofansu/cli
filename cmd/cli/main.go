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
	Login    LoginCommand    `command:"login" description:"Authenticate with a service"`
	Logout   LogoutCommand   `command:"logout" description:"Remove credentials for a service"`
	Bookmark BookmarkCommand `command:"bookmark" description:"Manage bookmarks (linkding)"`
	Feed     FeedCommand     `command:"feed" description:"Manage feeds (miniflux)"`
	Entry    EntryCommand    `command:"entry" description:"Manage entries (miniflux)"`
	Page     PageCommand     `command:"page" description:"Manage pages (wallabag)"`
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
		Service string `positional-arg-name:"service" description:"Service name (miniflux, linkding, or wallabag)" required:"yes"`
	} `positional-args:"yes"`
	Endpoint     string `long:"endpoint" description:"Service endpoint URL" required:"yes"`
	APIKey       string `long:"api-key" description:"API key (for miniflux or linkding)"`
	ClientID     string `long:"client-id" description:"OAuth client ID (for wallabag)"`
	ClientSecret string `long:"client-secret" description:"OAuth client secret (for wallabag)"`
	Username     string `long:"username" description:"Username (for wallabag)"`
	Password     string `long:"password" description:"Password (for wallabag)"`
}

type LogoutCommand struct {
	BaseCommand
	Args struct {
		Service string `positional-arg-name:"service" description:"Service name (miniflux, linkding, or wallabag)" required:"yes"`
	} `positional-args:"yes"`
}

type FeedAddCommand struct {
	BaseCommand
	Args struct {
		URL string `positional-arg-name:"url" description:"URL of the feed to subscribe to" required:"yes"`
	} `positional-args:"yes"`
	CategoryID int64 `long:"category-id" description:"Miniflux category ID (defaults to 1)"`
}

type BookmarkAddCommand struct {
	BaseCommand
	Args struct {
		URL string `positional-arg-name:"url" description:"URL of the bookmark to add" required:"yes"`
	} `positional-args:"yes"`
	Notes string `long:"notes" description:"Optional notes for the bookmark"`
	Tags  string `long:"tags" description:"Optional tags separated by spaces"`
}

type EntryListCommand struct {
	BaseCommand
	JSONOutputOptions
	Limit   int    `long:"limit" description:"Maximum number of results" default:"10"`
	Offset  int    `long:"offset" description:"Number of results to skip" default:"0"`
	Search  string `long:"search" description:"Search query text"`
	Status  string `long:"status" value-name:"status" description:"Filter by entry status (read, unread, removed)" default:"unread"`
	Starred bool   `long:"starred" description:"Filter by starred entries"`
	FeedID  int64  `long:"feed-id" description:"Filter by feed ID"`
}

type BookmarkListCommand struct {
	BaseCommand
	JSONOutputOptions
	Limit  int    `long:"limit" description:"Maximum number of results" default:"10"`
	Offset int    `long:"offset" description:"Number of results to skip" default:"0"`
	Search string `long:"search" description:"Search query text"`
}

type BookmarkCommand struct {
	BaseCommand
	Add  BookmarkAddCommand  `command:"add" description:"Add a bookmark (linkding)"`
	List BookmarkListCommand `command:"list" description:"List bookmarks (linkding)"`
}

type FeedCommand struct {
	BaseCommand
	Add FeedAddCommand `command:"add" description:"Add a feed (miniflux)"`
}

type EntryCommand struct {
	BaseCommand
	List EntryListCommand `command:"list" description:"List entries (miniflux)"`
}

type PageAddCommand struct {
	BaseCommand
	Args struct {
		URL string `positional-arg-name:"url" description:"URL of the page to add" required:"yes"`
	} `positional-args:"yes"`
	Tags    string `long:"tags" description:"Tags separated by spaces"`
	Archive bool   `long:"archive" description:"Mark as archived"`
}

type PageListCommand struct {
	BaseCommand
	JSONOutputOptions
	Archive bool   `long:"archive" description:"Filter by archived status"`
	Starred bool   `long:"starred" description:"Filter by starred status"`
	Page    int    `long:"page" description:"Page number" default:"1"`
	PerPage int    `long:"per-page" description:"Items per page" default:"10"`
	Tags    string `long:"tags" description:"Tags separated by spaces"`
	Domain  string `long:"domain" description:"Filter by domain name"`
}

type PageCommand struct {
	BaseCommand
	Add  PageAddCommand  `command:"add" description:"Add a page (wallabag)"`
	List PageListCommand `command:"list" description:"List pages (wallabag)"`
}

func (c *LoginCommand) Execute(_ []string) error {
	switch c.Args.Service {
	case "miniflux", "linkding":
		if c.APIKey == "" {
			return fmt.Errorf("--api-key is required for %s", c.Args.Service)
		}
		return auth.Login(c.Args.Service, c.Endpoint, c.APIKey)
	case "wallabag":
		if c.ClientID == "" || c.ClientSecret == "" || c.Username == "" || c.Password == "" {
			return fmt.Errorf("--client-id, --client-secret, --username, --password required for wallabag")
		}
		return auth.LoginWallabag(c.Endpoint, c.ClientID, c.ClientSecret, c.Username, c.Password)
	default:
		return fmt.Errorf("unknown service: %s", c.Args.Service)
	}
}

func (c *LogoutCommand) Execute(_ []string) error {
	return auth.Logout(c.Args.Service)
}

func (c *FeedAddCommand) Execute(_ []string) error {
	opts := app.AddFeedOptions{
		URL:        c.Args.URL,
		CategoryID: c.CategoryID,
	}

	return c.App.AddFeed(opts)
}

func (c *BookmarkAddCommand) Execute(_ []string) error {
	opts := app.AddBookmarkOptions{
		URL:   c.Args.URL,
		Notes: c.Notes,
		Tags:  c.Tags,
	}

	return c.App.AddBookmark(opts)
}

func (c *EntryListCommand) Execute(_ []string) error {
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

func (c *BookmarkListCommand) Execute(_ []string) error {
	opts := app.ListBookmarksOptions{
		Query:  c.Search,
		Limit:  c.Limit,
		Offset: c.Offset,
		JSON:   c.JSON,
		JQ:     c.JQ,
	}
	return c.App.ListBookmarks(opts)
}

func (c *PageAddCommand) Execute(_ []string) error {
	opts := app.AddPageOptions{
		URL:     c.Args.URL,
		Tags:    c.Tags,
		Archive: c.Archive,
	}
	return c.App.AddPage(opts)
}

func (c *PageListCommand) Execute(_ []string) error {
	archive := -1
	if c.Archive {
		archive = 1
	}

	starred := -1
	if c.Starred {
		starred = 1
	}

	opts := app.ListPagesOptions{
		Archive: archive,
		Starred: starred,
		Page:    c.Page,
		PerPage: c.PerPage,
		Tags:    c.Tags,
		Domain:  c.Domain,
		JSON:    c.JSON,
		JQ:      c.JQ,
	}

	return c.App.ListPages(opts)
}

func (c *LoginCommand) Usage() string {
	return "<service> [OPTIONS]"
}

func (c *LogoutCommand) Usage() string {
	return "<service>"
}

func (c *FeedAddCommand) Usage() string {
	return "<url>"
}

func (c *BookmarkAddCommand) Usage() string {
	return "<url>"
}

func (c *EntryListCommand) Usage() string {
	return "[OPTIONS]"
}

func (c *BookmarkListCommand) Usage() string {
	return "[OPTIONS]"
}

func (c *PageAddCommand) Usage() string {
	return "<url>"
}

func (c *PageListCommand) Usage() string {
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
	opts.Bookmark.Add.App = application
	opts.Bookmark.List.App = application
	opts.Feed.Add.App = application
	opts.Entry.List.App = application
	opts.Page.Add.App = application
	opts.Page.List.App = application

	parser := flags.NewParser(&opts, flags.HelpFlag|flags.PassDoubleDash)
	parser.ShortDescription = "My command-line tool for agents"
	parser.LongDescription = "Manage bookmarks, RSS feeds, and pages from terminal.\n\nExamples:\ncli login linkding --endpoint https://linkding.example.com --api-key YOUR_API_KEY\ncli login miniflux --endpoint https://miniflux.example.com --api-key YOUR_API_KEY\ncli login wallabag --endpoint https://wallabag.example.com --client-id ID --client-secret SECRET --username USER --password PASS\ncli bookmark add https://example.com --tags \"cool useful\"\ncli bookmark list\ncli feed add https://blog.example.com/feed.xml\ncli entry list\ncli page add https://example.com/article --archive\ncli page list"

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
