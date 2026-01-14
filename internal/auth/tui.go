package auth

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/goofansu/mlwcli/internal/config"
)

// isServiceSignedIn checks if a service has credentials configured
func isServiceSignedIn(cfg *config.Config, service string) bool {
	if cfg == nil {
		return false
	}
	switch service {
	case config.ServiceLinkding:
		return cfg.Linkding.Endpoint != ""
	case config.ServiceMiniflux:
		return cfg.Miniflux.Endpoint != ""
	case config.ServiceWallabag:
		return cfg.Wallabag.Endpoint != ""
	}
	return false
}

// PromptServiceLoginTUI displays a radio button menu for service selection during login
// Shows a check mark next to already signed-in services
func PromptServiceLoginTUI() (string, error) {
	var service string

	// Load config to check which services are signed in
	cfg, err := config.Load()
	if err != nil && !os.IsNotExist(err) {
		return "", fmt.Errorf("failed to load config: %w", err)
	}

	// Build option labels with check marks for signed-in services
	linkdingLabel := "Linkding"
	if isServiceSignedIn(cfg, config.ServiceLinkding) {
		linkdingLabel = "Linkding ✓"
	}

	minifluxLabel := "Miniflux"
	if isServiceSignedIn(cfg, config.ServiceMiniflux) {
		minifluxLabel = "Miniflux ✓"
	}

	wallabagLabel := "Wallabag"
	if isServiceSignedIn(cfg, config.ServiceWallabag) {
		wallabagLabel = "Wallabag ✓"
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select a service to login").
				Description("Use arrow keys to navigate, Enter to select. ✓ = already signed in").
				Options(
					huh.NewOption(linkdingLabel, config.ServiceLinkding),
					huh.NewOption(minifluxLabel, config.ServiceMiniflux),
					huh.NewOption(wallabagLabel, config.ServiceWallabag),
				).
				Value(&service),
		),
	)

	err = form.Run()
	if err != nil {
		return "", fmt.Errorf("failed to select service: %w", err)
	}

	return service, nil
}

// PromptServiceLogoutTUI displays a radio button menu for service selection during logout
// Only shows services that are currently signed in
func PromptServiceLogoutTUI() (string, error) {
	var service string

	// Load config to check which services are signed in
	cfg, err := config.Load()
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("no services are currently signed in")
		}
		return "", fmt.Errorf("failed to load config: %w", err)
	}

	// Build list of signed-in services
	var options []huh.Option[string]
	if isServiceSignedIn(cfg, config.ServiceLinkding) {
		options = append(options, huh.NewOption("Linkding", config.ServiceLinkding))
	}
	if isServiceSignedIn(cfg, config.ServiceMiniflux) {
		options = append(options, huh.NewOption("Miniflux", config.ServiceMiniflux))
	}
	if isServiceSignedIn(cfg, config.ServiceWallabag) {
		options = append(options, huh.NewOption("Wallabag", config.ServiceWallabag))
	}

	if len(options) == 0 {
		return "", fmt.Errorf("no services are currently signed in")
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select a service to logout").
				Description("Use arrow keys to navigate, Enter to select").
				Options(options...).
				Value(&service),
		),
	)

	err = form.Run()
	if err != nil {
		return "", fmt.Errorf("failed to select service: %w", err)
	}

	return service, nil
}

// PromptLinkdingCredentialsTUI prompts for Linkding credentials using TUI
func PromptLinkdingCredentialsTUI() (endpoint, apiKey string, err error) {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Endpoint URL").
				Placeholder("https://linkding.example.com").
				Value(&endpoint).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("endpoint URL is required")
					}
					return nil
				}),
			huh.NewInput().
				Title("API Key").
				EchoMode(huh.EchoModePassword).
				Value(&apiKey).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("API key is required")
					}
					return nil
				}),
		),
	)

	err = form.Run()
	return
}

// PromptMinifluxCredentialsTUI prompts for Miniflux credentials using TUI
func PromptMinifluxCredentialsTUI() (endpoint, apiKey string, err error) {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Endpoint URL").
				Placeholder("https://miniflux.example.com").
				Value(&endpoint).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("endpoint URL is required")
					}
					return nil
				}),
			huh.NewInput().
				Title("API Key").
				EchoMode(huh.EchoModePassword).
				Value(&apiKey).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("API key is required")
					}
					return nil
				}),
		),
	)

	err = form.Run()
	return
}

// PromptWallabagCredentialsTUI prompts for Wallabag credentials using TUI
func PromptWallabagCredentialsTUI() (endpoint, clientID, clientSecret, username, password string, err error) {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Endpoint URL").
				Placeholder("https://wallabag.example.com").
				Value(&endpoint).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("endpoint URL is required")
					}
					return nil
				}),
			huh.NewInput().
				Title("Client ID").
				Value(&clientID).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("client ID is required")
					}
					return nil
				}),
			huh.NewInput().
				Title("Client Secret").
				EchoMode(huh.EchoModePassword).
				Value(&clientSecret).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("client secret is required")
					}
					return nil
				}),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("Username").
				Value(&username).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("username is required")
					}
					return nil
				}),
			huh.NewInput().
				Title("Password").
				EchoMode(huh.EchoModePassword).
				Value(&password).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("password is required")
					}
					return nil
				}),
		),
	)

	err = form.Run()
	return
}
