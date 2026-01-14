package auth

import (
	"bufio"
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

var stdinReader = bufio.NewReader(os.Stdin)

func PromptInput(prompt string) (string, error) {
	fmt.Print(prompt)
	input, err := stdinReader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read input: %w", err)
	}
	return strings.TrimSpace(input), nil
}

func PromptService() (string, error) {
	fmt.Println("Select a service:")
	fmt.Println("1. Linkding")
	fmt.Println("2. Miniflux")
	fmt.Println("3. Wallabag")
	fmt.Print("Enter your choice (1-3): ")

	choice, err := stdinReader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read choice: %w", err)
	}

	choice = strings.TrimSpace(choice)
	switch choice {
	case "1":
		return config.ServiceLinkding, nil
	case "2":
		return config.ServiceMiniflux, nil
	case "3":
		return config.ServiceWallabag, nil
	default:
		return "", fmt.Errorf("invalid choice: %s", choice)
	}
}

func GetSecretOrPrompt(secret, prompt string) (string, error) {
	if secret != "" {
		return secret, nil
	}
	return PromptSecret(prompt)
}

// normalizeEndpoint removes trailing slashes and trims whitespace from endpoint URLs
func normalizeEndpoint(endpoint string) string {
	endpoint = strings.TrimSpace(endpoint)
	return strings.TrimRight(endpoint, "/")
}

func LoginMiniflux(endpoint, apiKey string) error {
	endpoint = normalizeEndpoint(endpoint)
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
	endpoint = normalizeEndpoint(endpoint)
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
	endpoint = normalizeEndpoint(endpoint)
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

func Login() error {
	service, err := PromptServiceLoginTUI()
	if err != nil {
		return err
	}

	switch service {
	case config.ServiceLinkding:
		return loginLinkdingInteractive()
	case config.ServiceMiniflux:
		return loginMinifluxInteractive()
	case config.ServiceWallabag:
		return loginWallabagInteractive()
	default:
		return fmt.Errorf("unknown service: %s", service)
	}
}

func loginLinkdingInteractive() error {
	endpoint, apiKey, err := PromptLinkdingCredentialsTUI()
	if err != nil {
		return err
	}

	return LoginLinkding(endpoint, apiKey)
}

func loginMinifluxInteractive() error {
	endpoint, apiKey, err := PromptMinifluxCredentialsTUI()
	if err != nil {
		return err
	}

	return LoginMiniflux(endpoint, apiKey)
}

func loginWallabagInteractive() error {
	endpoint, clientID, clientSecret, username, password, err := PromptWallabagCredentialsTUI()
	if err != nil {
		return err
	}

	return LoginWallabag(endpoint, clientID, clientSecret, username, password)
}

func Logout() error {
	service, err := PromptServiceLogoutTUI()
	if err != nil {
		return err
	}

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
