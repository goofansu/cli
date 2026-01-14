package auth

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/goofansu/mlwcli/internal/config"
	"github.com/goofansu/mlwcli/internal/linkding"
	"github.com/goofansu/mlwcli/internal/miniflux"
	"github.com/goofansu/mlwcli/internal/wallabag"
	"golang.org/x/term"
)

func PromptSecret(prompt string) (string, error) {
	fmt.Print(prompt)
	byteSecret, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return "", fmt.Errorf("failed to read secret: %w", err)
	}
	return string(byteSecret), nil
}

func GetSecretOrPrompt(secret, prompt string) (string, error) {
	if secret != "" {
		return secret, nil
	}
	return PromptSecret(prompt)
}

func LoginMiniflux(endpoint, apiKey string) error {
	endpoint = strings.TrimSpace(endpoint)
	apiKey = strings.TrimSpace(apiKey)

	if err := miniflux.Validate(endpoint, apiKey); err != nil {
		return fmt.Errorf("failed to verify miniflux connection: %w", err)
	}

	if err := saveMinifluxConfig(endpoint, apiKey); err != nil {
		return err
	}
	fmt.Println("✓ Configuration saved successfully")

	return nil
}

func saveMinifluxConfig(endpoint, apiKey string) error {
	cfg := config.ServiceConfig{
		Endpoint: endpoint,
		APIKey:   apiKey,
	}

	appCfg, err := config.Load()
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to load config: %w", err)
	}
	if appCfg == nil {
		appCfg = &config.Config{}
	}

	appCfg.Miniflux = cfg

	if err := config.Save(appCfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

func LoginLinkding(endpoint, apiKey string) error {
	endpoint = strings.TrimSpace(endpoint)
	apiKey = strings.TrimSpace(apiKey)

	if err := linkding.Validate(endpoint, apiKey); err != nil {
		return fmt.Errorf("failed to verify linkding connection: %w", err)
	}

	if err := saveLinkdingConfig(endpoint, apiKey); err != nil {
		return err
	}
	fmt.Println("✓ Configuration saved successfully")

	return nil
}

func saveLinkdingConfig(endpoint, apiKey string) error {
	cfg := config.ServiceConfig{
		Endpoint: endpoint,
		APIKey:   apiKey,
	}

	appCfg, err := config.Load()
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to load config: %w", err)
	}
	if appCfg == nil {
		appCfg = &config.Config{}
	}

	appCfg.Linkding = cfg

	if err := config.Save(appCfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

func LoginWallabag(endpoint, clientID, clientSecret, username, password string) error {
	endpoint = strings.TrimSpace(endpoint)
	clientID = strings.TrimSpace(clientID)
	clientSecret = strings.TrimSpace(clientSecret)
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)

	wallabag.LoadConfig(endpoint, clientID, clientSecret, username, password)

	if err := wallabag.Validate(); err != nil {
		return fmt.Errorf("failed to verify wallabag connection: %w", err)
	}

	cfg := config.WallabagConfig{
		Endpoint:     endpoint,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Username:     username,
		Password:     password,
	}

	if err := saveWallabagConfig(cfg); err != nil {
		return err
	}
	fmt.Println("✓ Configuration saved successfully")

	return nil
}

func saveWallabagConfig(cfg config.WallabagConfig) error {
	appCfg, err := config.Load()
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to load config: %w", err)
	}
	if appCfg == nil {
		appCfg = &config.Config{}
	}

	appCfg.Wallabag = cfg

	if err := config.Save(appCfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

func Logout(service string) error {
	service = strings.ToLower(strings.TrimSpace(service))

	switch service {
	case config.ServiceMiniflux, config.ServiceLinkding, config.ServiceWallabag:
		if err := config.RemoveService(service); err != nil {
			return fmt.Errorf("failed to remove %s config: %w", service, err)
		}
		fmt.Printf("✓ Logged out from %s successfully\n", service)
	default:
		return fmt.Errorf("invalid service: %s (must be '%s', '%s', or '%s')", service, config.ServiceMiniflux, config.ServiceLinkding, config.ServiceWallabag)
	}

	return nil
}
