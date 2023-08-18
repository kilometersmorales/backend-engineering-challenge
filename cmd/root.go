package cmd

import (
	"os"
	"unbabel/internal/parser"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "unbabel",
	Short: `A simple command line application that parses a stream of events ` +
		`and produces an aggregated output`,
	Long: `A simple command line application that parses a stream of events ` +
		`and produces an aggregated output. In this case, we're interested in` +
		` calculating, for every minute, a moving average of the translation ` +
		`delivery time for the last X minutes.`,
	Run: func(cmd *cobra.Command, args []string) {
		parser.Parse(parserConfig)
	},
}

var parserConfig parser.Config

// Execute adds all child commands to the root command and sets flags
// appropriately. This is called by main.main(). It only needs to happen once to
// the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// Init cobra and viper
func init() {

	// Flags

	// Debug flag (persistent)
	rootCmd.PersistentFlags().BoolVarP(&parserConfig.Debug, "debug", "d", false,
		"Display debugging output in the console (default: false)")
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	// Input file flag
	rootCmd.Flags().StringVar(&parserConfig.InputFile, "input_file", "events.json",
		"Input file")
	viper.BindPFlag("input_file", rootCmd.Flags().Lookup("input_file"))

	// Source language flag
	rootCmd.Flags().StringVar(&parserConfig.SourceLang, "source_language", "all",
		"Specify a specific source language")
	viper.BindPFlag("source_language", rootCmd.Flags().Lookup("source_language"))

	// Target language flag
	rootCmd.Flags().StringVar(&parserConfig.TargetLang, "target_language", "all",
		"Specify a specific target language")
	viper.BindPFlag("target_language", rootCmd.Flags().Lookup("target_language"))

	// Client name flag
	rootCmd.Flags().StringVar(&parserConfig.ClientName, "client_name", "all",
		"Specify a specific client name")
	viper.BindPFlag("client_name", rootCmd.Flags().Lookup("client_name"))

	// Window size flag
	rootCmd.Flags().IntVar(&parserConfig.WindowSize, "window_size", 10,
		"Specify a window size")
	viper.BindPFlag("window_size", rootCmd.Flags().Lookup("window_size"))
}
