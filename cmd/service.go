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
	serviceName  string
	ip           string
	port         string
	clusterName  string
	metadata     string
	ephemeral    bool
	prettyFormat bool
)

func init() {
	// Service command
	serviceCmd := &cobra.Command{
		Use:   "service",
		Short: "Manage Nacos services",
		Long:  `Service operations for Nacos server including register, deregister instances, etc.`,
	}

	// Get service command
	getServiceCmd := &cobra.Command{
		Use:   "get",
		Short: "Get service information from Nacos server",
		Long:  `Get service information from Nacos server by serviceName.`,
		Run: func(cmd *cobra.Command, args []string) {
			client := api.NewClient(viper.GetString("server"), viper.GetString("username"), viper.GetString("password"))
			serviceInfo, err := client.GetService(serviceName, viper.GetString("namespace"))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				os.Exit(1)
			}

			if prettyFormat {
				var jsonData interface{}
				if err := json.Unmarshal([]byte(serviceInfo), &jsonData); err != nil {
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
				fmt.Println(serviceInfo)
			}
		},
	}

	// Register instance command
	registerCmd := &cobra.Command{
		Use:   "register",
		Short: "Register an instance to Nacos server",
		Long:  `Register a service instance to Nacos server.`,
		Run: func(cmd *cobra.Command, args []string) {
			metadataMap := make(map[string]string)
			if metadata != "" {
				if err := json.Unmarshal([]byte(metadata), &metadataMap); err != nil {
					fmt.Fprintf(os.Stderr, "Error parsing metadata JSON: %s\n", err)
					os.Exit(1)
				}
			}

			client := api.NewClient(viper.GetString("server"), viper.GetString("username"), viper.GetString("password"))
			success, err := client.RegisterInstance(serviceName, ip, port, clusterName, viper.GetString("namespace"), metadataMap, ephemeral)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				os.Exit(1)
			}

			if success {
				fmt.Println("Instance registered successfully")
			} else {
				fmt.Println("Failed to register instance")
				os.Exit(1)
			}
		},
	}

	// Deregister instance command
	deregisterCmd := &cobra.Command{
		Use:   "deregister",
		Short: "Deregister an instance from Nacos server",
		Long:  `Deregister a service instance from Nacos server.`,
		Run: func(cmd *cobra.Command, args []string) {
			client := api.NewClient(viper.GetString("server"), viper.GetString("username"), viper.GetString("password"))
			success, err := client.DeregisterInstance(serviceName, ip, port, clusterName, viper.GetString("namespace"))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				os.Exit(1)
			}

			if success {
				fmt.Println("Instance deregistered successfully")
			} else {
				fmt.Println("Failed to deregister instance")
				os.Exit(1)
			}
		},
	}

	// Add flags for get service command
	getServiceCmd.Flags().StringVar(&serviceName, "name", "", "Service name (required)")
	getServiceCmd.Flags().BoolVar(&prettyFormat, "pretty", false, "Pretty format JSON output")
	getServiceCmd.MarkFlagRequired("name")

	// Add flags for register instance command
	registerCmd.Flags().StringVar(&serviceName, "name", "", "Service name (required)")
	registerCmd.Flags().StringVar(&ip, "ip", "", "IP address (required)")
	registerCmd.Flags().StringVar(&port, "port", "", "Port (required)")
	registerCmd.Flags().StringVar(&clusterName, "cluster", "", "Cluster name")
	registerCmd.Flags().StringVar(&metadata, "metadata", "", "Metadata in JSON format")
	registerCmd.Flags().BoolVar(&ephemeral, "ephemeral", true, "Whether instance is ephemeral")
	registerCmd.MarkFlagRequired("name")
	registerCmd.MarkFlagRequired("ip")
	registerCmd.MarkFlagRequired("port")

	// Add flags for deregister instance command
	deregisterCmd.Flags().StringVar(&serviceName, "name", "", "Service name (required)")
	deregisterCmd.Flags().StringVar(&ip, "ip", "", "IP address (required)")
	deregisterCmd.Flags().StringVar(&port, "port", "", "Port (required)")
	deregisterCmd.Flags().StringVar(&clusterName, "cluster", "", "Cluster name")
	deregisterCmd.MarkFlagRequired("name")
	deregisterCmd.MarkFlagRequired("ip")
	deregisterCmd.MarkFlagRequired("port")

	// Add commands to service command
	serviceCmd.AddCommand(getServiceCmd)
	serviceCmd.AddCommand(registerCmd)
	serviceCmd.AddCommand(deregisterCmd)

	// Add service command to root command
	rootCmd.AddCommand(serviceCmd)
}
