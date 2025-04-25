package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

type UploadOptions struct {
	MaxViews     int  `json:"max_views"`
	OriginalName bool `json:"original_name"`
	Clipboard    bool `json:"clipboard"`
}

var rootCmd = &cobra.Command{
	Use:   "zlx",
	Short: "zlx is a CLI tool to upload files",
	Long: `zlx is a simple CLI tool to upload files to a Zipline server.
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

		returnedURL, err := uploadFile(config, filePath)
		if err != nil {
			fmt.Println("Error uploading file:", err)
			os.Exit(1)
		}

		fmt.Println(returnedURL)

		if config.Upload.Clipboard {
			err := clipboard.WriteAll(returnedURL)
			if err != nil {
				fmt.Println("Clipboard write failed:", err)
			} else {
				fmt.Println("URL copied to clipboard!")
			}
		}
	},
}

func Execute() {
	cobra.OnInitialize(initConfig)
	var flags UploadOptions

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		os.Exit(1)
	}
	defaultConfig := filepath.Join(home, ".config", "zlx", "config.json")

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", defaultConfig, "config file")

	uploadCmd.Flags().IntVarP(&flags.MaxViews, "max-views", "m", 0, "Maximum views for the uploaded file")
	uploadCmd.Flags().BoolVarP(&flags.OriginalName, "original-name", "o", false, "Use original name for the uploaded file")
	uploadCmd.Flags().BoolVarP(&flags.Clipboard, "clipboard", "c", false, "Copy the returned URL to the clipboard")
	uploadCmd.PreRun = func(cmd *cobra.Command, args []string) {
		if cmd.Flags().Changed("max-views") {
			config.Upload.MaxViews = flags.MaxViews
		}
		if cmd.Flags().Changed("original-name") {
			config.Upload.OriginalName = flags.OriginalName
		}
		if cmd.Flags().Changed("clipboard") {
			config.Upload.Clipboard = flags.Clipboard
		}
	}

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
