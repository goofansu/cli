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
	Link   LinkCommand   `command:"link" description:"Manage links (linkding)"`
	Feed   FeedCommand   `command:"feed" description:"Manage feeds (miniflux)"`
	Page   PageCommand   `command:"page" description:"Manage pages (wallabag)"`
}

type BaseCommand struct {
	App *app.App
}

type JSONOutputOptions struct {
	JSON string `long:"json" value-name:"fields" description:"Output JSON with the specified fields (comma-separated)"`
	JQ   string `long:"jq" value-name:"expression" description:"Filter JSON output using a jq expression (requires --json)"`
}

type LoginMinifluxCommand struct {
	BaseCommand
	Endpoint string `long:"endpoint" description:"Miniflux endpoint URL" required:"yes"`
	APIKey   string `long:"api-key" description:"API key" required:"yes"`
}

type LoginLinkdingCommand struct {
	BaseCommand
	Endpoint string `long:"endpoint" description:"Linkding endpoint URL" required:"yes"`
	APIKey   string `long:"api-key" description:"API key" required:"yes"`
}

type LoginWallabagCommand struct {
	BaseCommand
	Endpoint     string `long:"endpoint" description:"Wallabag endpoint URL" required:"yes"`
	ClientID     string `long:"client-id" description:"OAuth client ID" required:"yes"`
	ClientSecret string `long:"client-secret" description:"OAuth client secret" required:"yes"`
	Username     string `long:"username" description:"Username" required:"yes"`
	Password     string `long:"password" description:"Password" required:"yes"`
}

type LoginCommand struct {
	Miniflux LoginMinifluxCommand `command:"miniflux" description:"Authenticate with Miniflux"`
	Linkding LoginLinkdingCommand `command:"linkding" description:"Authenticate with Linkding"`
	Wallabag LoginWallabagCommand `command:"wallabag" description:"Authenticate with Wallabag"`
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

type LinkAddCommand struct {
	BaseCommand
	Args struct {
		URL string `positional-arg-name:"url" description:"URL of link to add" required:"yes"`
	} `positional-args:"yes"`
	Notes string `long:"notes" description:"Optional notes for link"`
	Tags  string `long:"tags" description:"Optional tags separated by spaces"`
}

type FeedListCommand struct {
	BaseCommand
	JSONOutputOptions
	Limit   int    `long:"limit" description:"Maximum number of results" default:"10"`
	Offset  int    `long:"offset" description:"Number of results to skip" default:"0"`
	Search  string `long:"search" description:"Search query text"`
	Status  string `long:"status" value-name:"status" description:"Filter by feed entry status (read, unread, removed)" default:"unread"`
	Starred bool   `long:"starred" description:"Filter by starred feed entries"`
	FeedID  int64  `long:"feed-id" description:"Filter by feed ID"`
}

type LinkListCommand struct {
	BaseCommand
	JSONOutputOptions
	Limit  int    `long:"limit" description:"Maximum number of results" default:"10"`
	Offset int    `long:"offset" description:"Number of results to skip" default:"0"`
	Search string `long:"search" description:"Search query text"`
}

type LinkCommand struct {
	BaseCommand
	Add  LinkAddCommand  `command:"add" description:"Add a link (linkding)"`
	List LinkListCommand `command:"list" description:"List links (linkding)"`
}

type FeedCommand struct {
	BaseCommand
	Add  FeedAddCommand  `command:"add" description:"Add a feed (miniflux)"`
	List FeedListCommand `command:"list" description:"List feed entries (miniflux)"`
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

func (c *LoginMinifluxCommand) Execute(_ []string) error {
	return auth.LoginMiniflux(c.Endpoint, c.APIKey)
}

func (c *LoginLinkdingCommand) Execute(_ []string) error {
	return auth.LoginLinkding(c.Endpoint, c.APIKey)
}

func (c *LoginWallabagCommand) Execute(_ []string) error {
	return auth.LoginWallabag(c.Endpoint, c.ClientID, c.ClientSecret, c.Username, c.Password)
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

func (c *LinkAddCommand) Execute(_ []string) error {
	opts := app.AddLinkOptions{
		URL:   c.Args.URL,
		Notes: c.Notes,
		Tags:  c.Tags,
	}

	return c.App.AddLink(opts)
}

func (c *FeedListCommand) Execute(_ []string) error {
	starred := ""
	if c.Starred {
		starred = "1"
	}

	opts := app.EntriesOptions{
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

func (c *LinkListCommand) Execute(_ []string) error {
	opts := app.ListLinksOptions{
		Query:  c.Search,
		Limit:  c.Limit,
		Offset: c.Offset,
		JSON:   c.JSON,
		JQ:     c.JQ,
	}
	return c.App.ListLinks(opts)
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

func (c *LogoutCommand) Usage() string {
	return "<service>"
}

func (c *FeedAddCommand) Usage() string {
	return "<url>"
}

func (c *LinkAddCommand) Usage() string {
	return "<url>"
}

func (c *FeedListCommand) Usage() string {
	return "[OPTIONS]"
}

func (c *LinkListCommand) Usage() string {
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
	opts.Login.Miniflux.App = application
	opts.Login.Linkding.App = application
	opts.Login.Wallabag.App = application
	opts.Logout.App = application
	opts.Link.Add.App = application
	opts.Link.List.App = application
	opts.Feed.Add.App = application
	opts.Feed.List.App = application
	opts.Page.Add.App = application
	opts.Page.List.App = application

	parser := flags.NewParser(&opts, flags.HelpFlag|flags.PassDoubleDash)
	parser.ShortDescription = "My command-line tool for agents"
	parser.LongDescription = "Manage links, RSS feeds, and pages from terminal.\n\nExamples:\ncli link add https://example.com --tags \"cool useful\"\ncli link list\ncli feed add https://blog.example.com/feed.xml\ncli feed list\ncli page add https://example.com/article --archive\ncli page list"

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
