package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"zb/internal/api"
	"zb/internal/config"
)

func newLoginCmd(options *Options) *cobra.Command {
	var apiKey string
	var apiURL string
	var identityURL string

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Save an API key for the Zero Budget API",
		RunE: func(cmd *cobra.Command, _ []string) error {
			resolvedAPIKey, err := resolveAPIKey(apiKey, cmd.ErrOrStderr())
			if err != nil {
				return err
			}

			store, err := config.NewStore(options.ConfigPath)
			if err != nil {
				return err
			}

			client := api.NewClient(api.Config{
				BaseURL: apiURL,
				APIKey:  resolvedAPIKey,
			})

			ctx, cancel := context.WithTimeout(cmd.Context(), 10*time.Second)
			defer cancel()

			if _, err := client.ListCurrencies(ctx, api.ListCurrenciesParams{Limit: 1}); err != nil {
				return fmt.Errorf("failed to validate API key: %w", err)
			}

			cfg := config.Config{
				APIURL:      apiURL,
				IdentityURL: identityURL,
				APIKey:      resolvedAPIKey,
			}

			if err := store.Save(cfg); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Saved CLI config to %s\n", store.Path())

			return nil
		},
	}

	cmd.Flags().StringVar(&apiKey, "api-key", "", "Existing API key to save")
	cmd.Flags().StringVar(&apiURL, "api-url", envOrDefault("ZB_API_URL", "http://localhost:8080"), "Zero Budget API base URL")
	cmd.Flags().StringVar(&identityURL, "identity-url", envOrDefault("ZB_IDENTITY_URL", "http://localhost:8081"), "Identity service base URL")

	return cmd
}

func resolveAPIKey(apiKey string, output io.Writer) (string, error) {
	if strings.TrimSpace(apiKey) != "" {
		return strings.TrimSpace(apiKey), nil
	}

	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return "", fmt.Errorf("api key is required when stdin is not a terminal")
	}

	fmt.Fprint(output, "API key: ")

	value, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Fprintln(output)
	if err != nil {
		return "", fmt.Errorf("read API key: %w", err)
	}

	resolvedAPIKey := strings.TrimSpace(string(value))
	if resolvedAPIKey == "" {
		return "", fmt.Errorf("api key is required")
	}

	return resolvedAPIKey, nil
}

func envOrDefault(key string, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}

	return fallback
}
