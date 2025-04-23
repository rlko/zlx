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
	ServerName   string `json:"servername"`
	HTTPInsecure bool   `json:"http_insecure"`
	Token        string `json:"token"`
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
		config.HTTPInsecure = false
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
		case "token":
			fmt.Println(config.Token)
		case "http_insecure":
			fmt.Println(config.HTTPInsecure)
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
		fmt.Println("Servername:", config.ServerName)
		fmt.Println("Token:", config.Token)
		fmt.Println("HTTPInsecure:", config.HTTPInsecure)
	},
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
		case "token":
			existingConfig.Token = value
		case "http_insecure":
			boolValue, err := strconv.ParseBool(value)
			if err != nil {
				fmt.Println("Error: http_insecure must be true or false")
				os.Exit(1)
			}
			existingConfig.HTTPInsecure = boolValue
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
