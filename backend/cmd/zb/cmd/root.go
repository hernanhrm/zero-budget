package cmd

import "github.com/spf13/cobra"

type Options struct {
	ConfigPath string
}

func NewRootCmd() *cobra.Command {
	options := &Options{}

	rootCmd := &cobra.Command{
		Use:           "zb",
		Short:         "Zero Budget CLI",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	rootCmd.PersistentFlags().StringVar(&options.ConfigPath, "config", "", "Path to CLI config file")

	rootCmd.AddCommand(newLoginCmd(options))
	rootCmd.AddCommand(newCurrenciesCmd(options))

	return rootCmd
}
