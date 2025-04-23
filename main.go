package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	maxViewsFlag     int
	originalNameFlag bool
)

var rootCmd = &cobra.Command{
	Use:   "zlx",
	Short: "zlx is a CLI tool to upload files",
	Long: `zlx is a simple CLI tool to upload files to a server.
It reads the servername and token from a config file.`,
}

var uploadCmd = &cobra.Command{
	Use:     "upload <file_path>",
	Aliases: []string{"up"},
	Short:   "Upload a file",
	Long:    `Upload a file to the server specified in the config file.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]

		returnedURL, err := uploadFile(config.ServerName, config.Token, filePath, maxViewsFlag, originalNameFlag)
		if err != nil {
			fmt.Println("Error uploading file:", err)
			os.Exit(1)
		}

		fmt.Println(returnedURL)
	},
}

func Execute() {
	cobra.OnInitialize(initConfig)

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		os.Exit(1)
	}
	defaultConfig := filepath.Join(home, ".config", "zlx", "config.json")

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", defaultConfig, "config file")

	uploadCmd.Flags().IntVarP(&maxViewsFlag, "max-views", "m", 0, "Maximum views for the uploaded file")
	uploadCmd.Flags().BoolVarP(&originalNameFlag, "original-name", "o", false, "Use original name for the uploaded file")

	rootCmd.AddCommand(uploadCmd)
	rootCmd.AddCommand(configCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
