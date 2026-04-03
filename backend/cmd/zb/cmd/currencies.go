package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"zb/internal/api"
	"zb/internal/config"
	"zb/internal/output"
)

func newCurrenciesCmd(options *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "currencies",
		Short: "Manage currencies",
	}

	cmd.AddCommand(newCurrenciesListCmd(options))

	return cmd
}

func newCurrenciesListCmd(options *Options) *cobra.Command {
	var limit int
	var offset int
	var asJSON bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List currencies",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if limit < 0 {
				return fmt.Errorf("limit must be greater than or equal to 0")
			}

			if offset < 0 {
				return fmt.Errorf("offset must be greater than or equal to 0")
			}

			store, err := config.NewStore(options.ConfigPath)
			if err != nil {
				return err
			}

			cfg, err := store.Load()
			if err != nil {
				return err
			}

			client := api.NewClient(api.Config{
				BaseURL: cfg.APIURL,
				APIKey:  cfg.APIKey,
			})

			ctx, cancel := context.WithTimeout(cmd.Context(), 10*time.Second)
			defer cancel()

			currencies, err := client.ListCurrencies(ctx, api.ListCurrenciesParams{
				Limit:  limit,
				Offset: offset,
			})
			if err != nil {
				return err
			}

			if asJSON {
				encoder := json.NewEncoder(cmd.OutOrStdout())
				encoder.SetIndent("", "  ")
				return encoder.Encode(currencies)
			}

			return output.WriteCurrencies(cmd.OutOrStdout(), currencies)
		},
	}

	cmd.Flags().IntVar(&limit, "limit", 10, "Maximum number of currencies to return")
	cmd.Flags().IntVar(&offset, "offset", 0, "Number of currencies to skip")
	cmd.Flags().BoolVar(&asJSON, "json", false, "Output JSON")

	return cmd
}
