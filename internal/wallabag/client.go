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
