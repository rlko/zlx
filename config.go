package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

type Config struct {
	ServerName string        `json:"servername"`
	PathName   string        `json:"pathname"`
	Token      string        `json:"token"`
	Upload     UploadOptions `json:"upload"`
}

var config Config
var cfgFile string

func initConfig() {
	// Determine config file path
	configPath := cfgFile

	jsonFile, err := os.Open(configPath)
	if err != nil {
		fmt.Printf("Config file not found at %s\n", configPath)
		// If the file doesn't exist, set default values
		config.PathName = "/api/upload"
		return
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		fmt.Printf("Error unmarshaling config: %s\n", err)
		return
	}
}

func validateConfig() []error {
	var errors []error
	if config.ServerName == "" {
		errors = append(errors, fmt.Errorf("\"servername\" cannot be empty in config file"))
	}
	if config.Token == "" {
		errors = append(errors, fmt.Errorf("\"token\" cannot be empty in config file"))
	}
	return errors
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage zlx configuration",
	Long:  `The config command allows you to list, get, and set configuration variables.`,
}

func init() {
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configValidateCmd)
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get the value of a configuration variable",
	Long:  `Get the value of a configuration variable.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		switch key {
		case "servername":
			fmt.Println(config.ServerName)
		case "pathname":
			fmt.Println(config.PathName)
		case "token":
			fmt.Println(config.Token)
		case "upload.max_views":
			fmt.Println(config.Upload.MaxViews)
		case "upload.original_name":
			fmt.Println(config.Upload.OriginalName)
		case "upload.clipboard":
			fmt.Println(config.Upload.Clipboard)
		default:
			fmt.Printf("Error: key '%s' not found\n", key)
			os.Exit(1)
		}
	},
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configuration variables",
	Long:  `List all configuration variables and their values.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("servername:", config.ServerName)
		fmt.Println("pathname:", config.PathName)
		fmt.Println("token:", config.Token)
		fmt.Println("upload.max_views:", config.Upload.MaxViews)
		fmt.Println("upload.original_name:", config.Upload.OriginalName)
		fmt.Println("upload.clipboard:", config.Upload.Clipboard)
	},
}

func getBoolValue(value string) bool {
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		fmt.Println("Error: invalid value for boolean, must be 'true' or 'false'")
		os.Exit(1)
	}
	return boolValue
}

func getIntValue(value string) int {
	intValue, err := strconv.Atoi(value)
	if err != nil {
		fmt.Println("Error: invalid value for integer, must be a valid number")
		os.Exit(1)
	}
	return intValue
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set the value of a configuration variable",
	Long:  `Set the value of a configuration variable.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := args[1]

		configDir := filepath.Join(os.Getenv("HOME"), ".config", "zlx")
		configPath := filepath.Join(configDir, "config.json")

		if _, err := os.Stat(configDir); os.IsNotExist(err) {
			os.MkdirAll(configDir, 0700)
		}

		existingConfig := Config{}
		jsonFile, err := os.Open(configPath)
		if err == nil {
			byteValue, _ := io.ReadAll(jsonFile)
			json.Unmarshal(byteValue, &existingConfig)
			jsonFile.Close()
		}

		switch key = strings.ToLower(key); key {
		case "servername":
			existingConfig.ServerName = value
		case "pathname":
			existingConfig.PathName = value
		case "token":
			existingConfig.Token = value
		case "upload.clipboard":
			existingConfig.Upload.Clipboard = getBoolValue(value)
		case "upload.max_views":
			existingConfig.Upload.MaxViews = getIntValue(value)
		case "upload.original_name":
			existingConfig.Upload.OriginalName = getBoolValue(value)
		default:
			fmt.Printf("Error: key '%s' not found\n", key)
			os.Exit(1)
		}

		jsonValue, _ := json.MarshalIndent(existingConfig, "", "    ")

		err = os.WriteFile(configPath, jsonValue, 0600)
		if err != nil {
			fmt.Println("Error writing config file:", err)
			os.Exit(1)
		}

		fmt.Printf("Set %s to %s\n", key, value)
	},
}

var configValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate the configuration file",
	Long:  `Validate the configuration file.`,
	Run: func(cmd *cobra.Command, args []string) {
		errors := validateConfig()
		if len(errors) > 0 {
			fmt.Println("Invalid configuration:")
			for _, err := range errors {
				fmt.Println("-", err)
			}
			os.Exit(1)
		}
		fmt.Println("Config file is valid.")
	},
}
