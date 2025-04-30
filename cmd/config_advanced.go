package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/k8scat/nacos-cli/pkg/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	nid      string
	pageNo   int
	pageSize int
)

// init function for advanced config commands
func init() {
	// Get the config command from the root command
	var configCmd *cobra.Command
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "config" {
			configCmd = cmd
			break
		}
	}

	// If configCmd is not found, something is wrong, so create it
	if configCmd == nil {
		configCmd = &cobra.Command{
			Use:   "config",
			Short: "Manage Nacos configurations",
			Long:  `Configure operations for Nacos server including get, publish, delete configurations.`,
		}
		rootCmd.AddCommand(configCmd)
	}

	// Listen config command
	listenConfigCmd := &cobra.Command{
		Use:   "listen",
		Short: "Listen for configuration changes",
		Long:  `Listen for configuration changes in Nacos server.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Calculate MD5 of content if available
			var contentMD5 string
			if content != "" {
				contentMD5 = api.GetMD5(content)
			}

			client := api.NewClient(viper.GetString("server"), viper.GetString("username"), viper.GetString("password"))
			result, err := client.ListenConfig(dataId, group, contentMD5, viper.GetString("namespace"))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				os.Exit(1)
			}

			if result == "" {
				fmt.Println("No configuration changes detected")
			} else {
				fmt.Println("Configuration changes detected:")
				fmt.Println(result)
			}
		},
	}

	// History config command
	historyConfigCmd := &cobra.Command{
		Use:   "history",
		Short: "Get configuration history",
		Long:  `Get configuration history versions from Nacos server.`,
		Run: func(cmd *cobra.Command, args []string) {
			client := api.NewClient(viper.GetString("server"), viper.GetString("username"), viper.GetString("password"))
			result, err := client.GetConfigHistory(dataId, group, viper.GetString("namespace"), pageNo, pageSize)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				os.Exit(1)
			}

			if prettyFormat {
				var jsonData interface{}
				if err := json.Unmarshal([]byte(result), &jsonData); err != nil {
					fmt.Fprintf(os.Stderr, "Error parsing JSON response: %s\n", err)
					os.Exit(1)
				}

				// Pretty print JSON
				prettyJSON, err := json.MarshalIndent(jsonData, "", "  ")
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error formatting JSON: %s\n", err)
					os.Exit(1)
				}
				fmt.Println(string(prettyJSON))
			} else {
				fmt.Println(result)
			}
		},
	}

	// History detail config command
	historyDetailConfigCmd := &cobra.Command{
		Use:   "history-detail",
		Short: "Get configuration history detail",
		Long:  `Get configuration history version detail from Nacos server.`,
		Run: func(cmd *cobra.Command, args []string) {
			client := api.NewClient(viper.GetString("server"), viper.GetString("username"), viper.GetString("password"))
			result, err := client.GetConfigHistoryDetail(dataId, group, viper.GetString("namespace"), nid)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				os.Exit(1)
			}

			if prettyFormat {
				var jsonData interface{}
				if err := json.Unmarshal([]byte(result), &jsonData); err != nil {
					fmt.Fprintf(os.Stderr, "Error parsing JSON response: %s\n", err)
					os.Exit(1)
				}

				// Pretty print JSON
				prettyJSON, err := json.MarshalIndent(jsonData, "", "  ")
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error formatting JSON: %s\n", err)
					os.Exit(1)
				}
				fmt.Println(string(prettyJSON))
			} else {
				fmt.Println(result)
			}
		},
	}

	// Previous config command
	previousConfigCmd := &cobra.Command{
		Use:   "previous",
		Short: "Get previous configuration",
		Long:  `Get previous configuration version from Nacos server.`,
		Run: func(cmd *cobra.Command, args []string) {
			client := api.NewClient(viper.GetString("server"), viper.GetString("username"), viper.GetString("password"))
			result, err := client.GetPreviousConfigInfo(dataId, group, viper.GetString("namespace"))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				os.Exit(1)
			}

			if prettyFormat {
				var jsonData interface{}
				if err := json.Unmarshal([]byte(result), &jsonData); err != nil {
					fmt.Fprintf(os.Stderr, "Error parsing JSON response: %s\n", err)
					os.Exit(1)
				}

				// Pretty print JSON
				prettyJSON, err := json.MarshalIndent(jsonData, "", "  ")
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error formatting JSON: %s\n", err)
					os.Exit(1)
				}
				fmt.Println(string(prettyJSON))
			} else {
				fmt.Println(result)
			}
		},
	}

	// Add flags for listen config command
	listenConfigCmd.Flags().StringVar(&dataId, "data-id", "", "Configuration ID (required)")
	listenConfigCmd.Flags().StringVar(&group, "group", "", "Configuration group (required)")
	listenConfigCmd.Flags().StringVar(&content, "content", "", "Current content for MD5 calculation")
	listenConfigCmd.MarkFlagRequired("data-id")
	listenConfigCmd.MarkFlagRequired("group")

	// Add flags for history config command
	historyConfigCmd.Flags().StringVar(&dataId, "data-id", "", "Configuration ID (required)")
	historyConfigCmd.Flags().StringVar(&group, "group", "", "Configuration group (required)")
	historyConfigCmd.Flags().IntVar(&pageNo, "page", 1, "Page number")
	historyConfigCmd.Flags().IntVar(&pageSize, "size", 10, "Page size")
	historyConfigCmd.Flags().BoolVar(&prettyFormat, "pretty", false, "Pretty format JSON output")
	historyConfigCmd.MarkFlagRequired("data-id")
	historyConfigCmd.MarkFlagRequired("group")

	// Add flags for history detail config command
	historyDetailConfigCmd.Flags().StringVar(&dataId, "data-id", "", "Configuration ID (required)")
	historyDetailConfigCmd.Flags().StringVar(&group, "group", "", "Configuration group (required)")
	historyDetailConfigCmd.Flags().StringVar(&nid, "nid", "", "History ID (required)")
	historyDetailConfigCmd.Flags().BoolVar(&prettyFormat, "pretty", false, "Pretty format JSON output")
	historyDetailConfigCmd.MarkFlagRequired("data-id")
	historyDetailConfigCmd.MarkFlagRequired("group")
	historyDetailConfigCmd.MarkFlagRequired("nid")

	// Add flags for previous config command
	previousConfigCmd.Flags().StringVar(&dataId, "data-id", "", "Configuration ID (required)")
	previousConfigCmd.Flags().StringVar(&group, "group", "", "Configuration group (required)")
	previousConfigCmd.Flags().BoolVar(&prettyFormat, "pretty", false, "Pretty format JSON output")
	previousConfigCmd.MarkFlagRequired("data-id")
	previousConfigCmd.MarkFlagRequired("group")

	// Add commands to config command
	configCmd.AddCommand(listenConfigCmd)
	configCmd.AddCommand(historyConfigCmd)
	configCmd.AddCommand(historyDetailConfigCmd)
	configCmd.AddCommand(previousConfigCmd)
}
