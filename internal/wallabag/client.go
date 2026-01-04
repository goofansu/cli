package wallabag

import (
	"fmt"
	"strings"

	"github.com/Strubbl/wallabago/v8"
)

func LoadConfig(wallabagURL, clientID, clientSecret, userName, userPassword string) {
	cfg := wallabago.NewWallabagConfig(wallabagURL, clientID, clientSecret, userName, userPassword)
	wallabago.SetConfig(cfg)
}

func Validate() error {
	if _, err := wallabago.GetAuthTokenHeader(); err != nil {
		return fmt.Errorf("failed to authenticate with wallabag: %w", err)
	}

	return nil
}

func CreatePage(url, tags string, archive bool) error {
	var archiveInt int
	if archive {
		archiveInt = 1
	}

	commaTags := ""
	if tags != "" {
		commaTags = strings.Join(strings.Fields(tags), ",")
	}

	if err := wallabago.PostEntry(url, "", commaTags, 0, archiveInt); err != nil {
		return fmt.Errorf("failed to create page: %w", err)
	}

	return nil
}

type ListPagesOptions struct {
	Archive int
	Starred int
	Page    int
	PerPage int
	Tags    string
	Domain  string
}

type ListPagesResult struct {
	Total int
	Items []wallabago.Item
}

func ListPages(opts ListPagesOptions) (*ListPagesResult, error) {
	tags := ""
	if opts.Tags != "" {
		tags = strings.Join(strings.Fields(opts.Tags), ",")
	}

	entries, err := wallabago.GetEntries(
		wallabago.APICall,
		opts.Archive,
		opts.Starred,
		"created",
		"desc",
		opts.Page,
		opts.PerPage,
		tags,
		0,
		0,
		"",
		opts.Domain,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list pages: %w", err)
	}

	return &ListPagesResult{
		Total: entries.Total,
		Items: entries.Embedded.Items,
	}, nil
}
