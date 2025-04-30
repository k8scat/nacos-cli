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
	namespaceId   string
	namespaceName string
	namespaceDesc string
)

func init() {
	// Namespace command
	namespaceCmd := &cobra.Command{
		Use:   "namespace",
		Short: "Manage Nacos namespaces",
		Long:  `Namespace operations for Nacos server including list, create, modify, delete namespaces.`,
	}

	// List namespace command
	listNamespaceCmd := &cobra.Command{
		Use:   "list",
		Short: "List namespaces from Nacos server",
		Long:  `List all namespaces from Nacos server.`,
		Run: func(cmd *cobra.Command, args []string) {
			client := api.NewClient(viper.GetString("server"), viper.GetString("username"), viper.GetString("password"))
			namespaces, err := client.ListNamespaces()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				os.Exit(1)
			}

			if prettyFormat {
				var jsonData interface{}
				if err := json.Unmarshal([]byte(namespaces), &jsonData); err != nil {
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
				fmt.Println(namespaces)
			}
		},
	}

	// Create namespace command
	createNamespaceCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a namespace in Nacos server",
		Long:  `Create a new namespace in Nacos server.`,
		Run: func(cmd *cobra.Command, args []string) {
			client := api.NewClient(viper.GetString("server"), viper.GetString("username"), viper.GetString("password"))
			success, err := client.CreateNamespace(namespaceId, namespaceName, namespaceDesc)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				os.Exit(1)
			}

			if success {
				fmt.Println("Namespace created successfully")
			} else {
				fmt.Println("Failed to create namespace")
				os.Exit(1)
			}
		},
	}

	// Modify namespace command
	modifyNamespaceCmd := &cobra.Command{
		Use:   "modify",
		Short: "Modify a namespace in Nacos server",
		Long:  `Modify an existing namespace in Nacos server.`,
		Run: func(cmd *cobra.Command, args []string) {
			client := api.NewClient(viper.GetString("server"), viper.GetString("username"), viper.GetString("password"))
			success, err := client.ModifyNamespace(namespaceId, namespaceName, namespaceDesc)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				os.Exit(1)
			}

			if success {
				fmt.Println("Namespace modified successfully")
			} else {
				fmt.Println("Failed to modify namespace")
				os.Exit(1)
			}
		},
	}

	// Delete namespace command
	deleteNamespaceCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a namespace from Nacos server",
		Long:  `Delete an existing namespace from Nacos server.`,
		Run: func(cmd *cobra.Command, args []string) {
			client := api.NewClient(viper.GetString("server"), viper.GetString("username"), viper.GetString("password"))
			success, err := client.DeleteNamespace(namespaceId)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				os.Exit(1)
			}

			if success {
				fmt.Println("Namespace deleted successfully")
			} else {
				fmt.Println("Failed to delete namespace")
				os.Exit(1)
			}
		},
	}

	// Add flags for list namespace command
	listNamespaceCmd.Flags().BoolVar(&prettyFormat, "pretty", false, "Pretty format JSON output")

	// Add flags for create namespace command
	createNamespaceCmd.Flags().StringVar(&namespaceId, "id", "", "Namespace ID (required)")
	createNamespaceCmd.Flags().StringVar(&namespaceName, "name", "", "Namespace name (required)")
	createNamespaceCmd.Flags().StringVar(&namespaceDesc, "desc", "", "Namespace description")
	createNamespaceCmd.MarkFlagRequired("id")
	createNamespaceCmd.MarkFlagRequired("name")

	// Add flags for modify namespace command
	modifyNamespaceCmd.Flags().StringVar(&namespaceId, "id", "", "Namespace ID (required)")
	modifyNamespaceCmd.Flags().StringVar(&namespaceName, "name", "", "Namespace name (required)")
	modifyNamespaceCmd.Flags().StringVar(&namespaceDesc, "desc", "", "Namespace description (required)")
	modifyNamespaceCmd.MarkFlagRequired("id")
	modifyNamespaceCmd.MarkFlagRequired("name")
	modifyNamespaceCmd.MarkFlagRequired("desc")

	// Add flags for delete namespace command
	deleteNamespaceCmd.Flags().StringVar(&namespaceId, "id", "", "Namespace ID (required)")
	deleteNamespaceCmd.MarkFlagRequired("id")

	// Add commands to namespace command
	namespaceCmd.AddCommand(listNamespaceCmd)
	namespaceCmd.AddCommand(createNamespaceCmd)
	namespaceCmd.AddCommand(modifyNamespaceCmd)
	namespaceCmd.AddCommand(deleteNamespaceCmd)

	// Add namespace command to root command
	rootCmd.AddCommand(namespaceCmd)
}
