package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/k8scat/nacos-cli/pkg/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	dataId     string
	group      string
	content    string
	configType string
	outputFile string
)

func init() {
	// Config command
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Manage Nacos configurations",
		Long:  `Configure operations for Nacos server including get, publish, delete configurations.`,
	}

	// Get config command
	getConfigCmd := &cobra.Command{
		Use:   "get",
		Short: "Get configuration from Nacos server",
		Long:  `Get configuration from Nacos server by dataId and group.`,
		Run: func(cmd *cobra.Command, args []string) {
			client := api.NewClient(viper.GetString("server"), viper.GetString("username"), viper.GetString("password"))
			content, err := client.GetConfig(dataId, group, viper.GetString("namespace"))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				os.Exit(1)
			}

			if outputFile != "" {
				// Create directory if not exists
				dir := filepath.Dir(outputFile)
				if dir != "." {
					if err := os.MkdirAll(dir, 0755); err != nil {
						fmt.Fprintf(os.Stderr, "Error creating directory: %s\n", err)
						os.Exit(1)
					}
				}

				// Write content to file
				if err := os.WriteFile(outputFile, []byte(content), 0644); err != nil {
					fmt.Fprintf(os.Stderr, "Error writing to file: %s\n", err)
					os.Exit(1)
				}
				fmt.Printf("Configuration saved to %s\n", outputFile)
			} else {
				fmt.Println(content)
			}
		},
	}

	// Publish config command
	publishConfigCmd := &cobra.Command{
		Use:   "publish",
		Short: "Publish configuration to Nacos server",
		Long:  `Publish configuration to Nacos server by dataId and group.`,
		Run: func(cmd *cobra.Command, args []string) {
			var contentData string
			if content != "" {
				contentData = content
			} else if outputFile != "" {
				// Read content from file
				data, err := os.ReadFile(outputFile)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err)
					os.Exit(1)
				}
				contentData = string(data)
			} else {
				// Read from stdin
				data, err := io.ReadAll(os.Stdin)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error reading from stdin: %s\n", err)
					os.Exit(1)
				}
				contentData = string(data)
			}

			client := api.NewClient(viper.GetString("server"), viper.GetString("username"), viper.GetString("password"))
			success, err := client.PublishConfig(dataId, group, contentData, configType, viper.GetString("namespace"))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				os.Exit(1)
			}

			if success {
				fmt.Println("Configuration published successfully")
			} else {
				fmt.Println("Failed to publish configuration")
				os.Exit(1)
			}
		},
	}

	// Delete config command
	deleteConfigCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete configuration from Nacos server",
		Long:  `Delete configuration from Nacos server by dataId and group.`,
		Run: func(cmd *cobra.Command, args []string) {
			client := api.NewClient(viper.GetString("server"), viper.GetString("username"), viper.GetString("password"))
			success, err := client.DeleteConfig(dataId, group, viper.GetString("namespace"))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				os.Exit(1)
			}

			if success {
				fmt.Println("Configuration deleted successfully")
			} else {
				fmt.Println("Failed to delete configuration")
				os.Exit(1)
			}
		},
	}

	// Add flags for get config command
	getConfigCmd.Flags().StringVar(&dataId, "data-id", "", "Configuration ID (required)")
	getConfigCmd.Flags().StringVar(&group, "group", "", "Configuration group (required)")
	getConfigCmd.Flags().StringVar(&outputFile, "output", "", "Output file path to save the configuration")
	getConfigCmd.MarkFlagRequired("data-id")
	getConfigCmd.MarkFlagRequired("group")

	// Add flags for publish config command
	publishConfigCmd.Flags().StringVar(&dataId, "data-id", "", "Configuration ID (required)")
	publishConfigCmd.Flags().StringVar(&group, "group", "", "Configuration group (required)")
	publishConfigCmd.Flags().StringVar(&content, "content", "", "Configuration content")
	publishConfigCmd.Flags().StringVar(&configType, "type", "", "Configuration type")
	publishConfigCmd.Flags().StringVar(&outputFile, "file", "", "File path to read configuration content")
	publishConfigCmd.MarkFlagRequired("data-id")
	publishConfigCmd.MarkFlagRequired("group")

	// Add flags for delete config command
	deleteConfigCmd.Flags().StringVar(&dataId, "data-id", "", "Configuration ID (required)")
	deleteConfigCmd.Flags().StringVar(&group, "group", "", "Configuration group (required)")
	deleteConfigCmd.MarkFlagRequired("data-id")
	deleteConfigCmd.MarkFlagRequired("group")

	// Add commands to config command
	configCmd.AddCommand(getConfigCmd)
	configCmd.AddCommand(publishConfigCmd)
	configCmd.AddCommand(deleteConfigCmd)

	// Add config command to root command
	rootCmd.AddCommand(configCmd)
}
