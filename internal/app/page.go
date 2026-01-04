package app

import (
	"fmt"

	"github.com/goofansu/cli/internal/wallabag"
)

type AddPageOptions struct {
	URL     string
	Tags    string
	Archive bool
}

func (a *App) AddPage(opts AddPageOptions) error {
	if a.Config.Wallabag.Endpoint == "" {
		return fmt.Errorf("wallabag not configured. Run 'cli login wallabag' first")
	}

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
