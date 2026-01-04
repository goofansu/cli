package app

import (
	"fmt"

	"github.com/goofansu/cli/internal/format"
	"github.com/goofansu/cli/internal/wallabag"
)

type AddPageOptions struct {
	URL     string
	Tags    string
	Archive bool
}

type ListPagesOptions struct {
	Archive int
	Starred int
	Page    int
	PerPage int
	Tags    string
	Domain  string
	JSON    string
	JQ      string
}

func (a *App) AddPage(opts AddPageOptions) error {
	wallabag.LoadConfig(
		a.Config.Wallabag.Endpoint,
		a.Config.Wallabag.ClientID,
		a.Config.Wallabag.ClientSecret,
		a.Config.Wallabag.Username,
		a.Config.Wallabag.Password,
	)

	if err := wallabag.CreatePage(opts.URL, opts.Tags, opts.Archive); err != nil {
		return err
	}

	fmt.Printf("✓ Page created successfully\n")
	return nil
}

func (a *App) ListPages(opts ListPagesOptions) error {
	wallabag.LoadConfig(
		a.Config.Wallabag.Endpoint,
		a.Config.Wallabag.ClientID,
		a.Config.Wallabag.ClientSecret,
		a.Config.Wallabag.Username,
		a.Config.Wallabag.Password,
	)

	listOpts := wallabag.ListPagesOptions{
		Archive: opts.Archive,
		Starred: opts.Starred,
		Page:    opts.Page,
		PerPage: opts.PerPage,
		Tags:    opts.Tags,
		Domain:  opts.Domain,
	}

	result, err := wallabag.ListPages(listOpts)
	if err != nil {
		return err
	}

	data := map[string]any{
		"total": result.Total,
		"items": result.Items,
	}

	return format.Output(data, opts.JSON, opts.JQ)
}
