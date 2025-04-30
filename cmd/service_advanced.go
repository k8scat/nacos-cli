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
	weight           float64
	enabled          bool
	clusters         string
	healthyOnly      string
	groupName        string
	protectThreshold string
	switchEntry      string
	switchValue      string
	healthy          bool
)

// init function for advanced service commands
func init() {
	// Get the service command from the root command
	var serviceCmd *cobra.Command
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "service" {
			serviceCmd = cmd
			break
		}
	}

	// If serviceCmd is not found, something is wrong, so create it
	if serviceCmd == nil {
		serviceCmd = &cobra.Command{
			Use:   "service",
			Short: "Manage Nacos services",
			Long:  `Service operations for Nacos server including register, deregister instances, etc.`,
		}
		rootCmd.AddCommand(serviceCmd)
	}

	// Modify instance command
	modifyInstanceCmd := &cobra.Command{
		Use:   "modify-instance",
		Short: "Modify a service instance",
		Long:  `Modify a service instance in Nacos server.`,
		Run: func(cmd *cobra.Command, args []string) {
			metadataMap := make(map[string]string)
			if metadata != "" {
				if err := json.Unmarshal([]byte(metadata), &metadataMap); err != nil {
					fmt.Fprintf(os.Stderr, "Error parsing metadata JSON: %s\n", err)
					os.Exit(1)
				}
			}

			client := api.NewClient(viper.GetString("server"), viper.GetString("username"), viper.GetString("password"))
			success, err := client.ModifyInstance(serviceName, ip, port, clusterName, viper.GetString("namespace"), weight, metadataMap, enabled)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				os.Exit(1)
			}

			if success {
				fmt.Println("Instance modified successfully")
			} else {
				fmt.Println("Failed to modify instance")
				os.Exit(1)
			}
		},
	}

	// List instances command
	listInstancesCmd := &cobra.Command{
		Use:   "list-instances",
		Short: "List service instances",
		Long:  `List service instances from Nacos server.`,
		Run: func(cmd *cobra.Command, args []string) {
			client := api.NewClient(viper.GetString("server"), viper.GetString("username"), viper.GetString("password"))
			result, err := client.ListInstances(serviceName, viper.GetString("namespace"), clusters, healthyOnly)
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

	// Get instance command
	getInstanceCmd := &cobra.Command{
		Use:   "get-instance",
		Short: "Get instance details",
		Long:  `Get service instance details from Nacos server.`,
		Run: func(cmd *cobra.Command, args []string) {
			client := api.NewClient(viper.GetString("server"), viper.GetString("username"), viper.GetString("password"))
			result, err := client.GetInstance(serviceName, ip, port, viper.GetString("namespace"), clusterName)
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

	// Beat command
	beatCmd := &cobra.Command{
		Use:   "beat",
		Short: "Send heartbeat for an instance",
		Long:  `Send heartbeat for a service instance to Nacos server.`,
		Run: func(cmd *cobra.Command, args []string) {
			metadataMap := make(map[string]string)
			if metadata != "" {
				if err := json.Unmarshal([]byte(metadata), &metadataMap); err != nil {
					fmt.Fprintf(os.Stderr, "Error parsing metadata JSON: %s\n", err)
					os.Exit(1)
				}
			}

			client := api.NewClient(viper.GetString("server"), viper.GetString("username"), viper.GetString("password"))
			success, err := client.SendInstanceHeartbeat(serviceName, ip, port, viper.GetString("namespace"), clusterName, weight, metadataMap)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				os.Exit(1)
			}

			if success {
				fmt.Println("Heartbeat sent successfully")
			} else {
				fmt.Println("Failed to send heartbeat")
				os.Exit(1)
			}
		},
	}

	// Create service command
	createServiceCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a service",
		Long:  `Create a new service in Nacos server.`,
		Run: func(cmd *cobra.Command, args []string) {
			metadataMap := make(map[string]string)
			if metadata != "" {
				if err := json.Unmarshal([]byte(metadata), &metadataMap); err != nil {
					fmt.Fprintf(os.Stderr, "Error parsing metadata JSON: %s\n", err)
					os.Exit(1)
				}
			}

			client := api.NewClient(viper.GetString("server"), viper.GetString("username"), viper.GetString("password"))
			success, err := client.CreateService(serviceName, viper.GetString("namespace"), groupName, protectThreshold, metadataMap)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				os.Exit(1)
			}

			if success {
				fmt.Println("Service created successfully")
			} else {
				fmt.Println("Failed to create service")
				os.Exit(1)
			}
		},
	}

	// Delete service command
	deleteServiceCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a service",
		Long:  `Delete a service from Nacos server.`,
		Run: func(cmd *cobra.Command, args []string) {
			client := api.NewClient(viper.GetString("server"), viper.GetString("username"), viper.GetString("password"))
			success, err := client.DeleteService(serviceName, viper.GetString("namespace"), groupName)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				os.Exit(1)
			}

			if success {
				fmt.Println("Service deleted successfully")
			} else {
				fmt.Println("Failed to delete service")
				os.Exit(1)
			}
		},
	}

	// Update service command
	updateServiceCmd := &cobra.Command{
		Use:   "update",
		Short: "Update a service",
		Long:  `Update a service in Nacos server.`,
		Run: func(cmd *cobra.Command, args []string) {
			metadataMap := make(map[string]string)
			if metadata != "" {
				if err := json.Unmarshal([]byte(metadata), &metadataMap); err != nil {
					fmt.Fprintf(os.Stderr, "Error parsing metadata JSON: %s\n", err)
					os.Exit(1)
				}
			}

			client := api.NewClient(viper.GetString("server"), viper.GetString("username"), viper.GetString("password"))
			success, err := client.UpdateService(serviceName, viper.GetString("namespace"), groupName, protectThreshold, metadataMap)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				os.Exit(1)
			}

			if success {
				fmt.Println("Service updated successfully")
			} else {
				fmt.Println("Failed to update service")
				os.Exit(1)
			}
		},
	}

	// List services command
	listServicesCmd := &cobra.Command{
		Use:   "list",
		Short: "List services",
		Long:  `List services from Nacos server.`,
		Run: func(cmd *cobra.Command, args []string) {
			client := api.NewClient(viper.GetString("server"), viper.GetString("username"), viper.GetString("password"))
			result, err := client.ListServices(viper.GetString("namespace"), groupName, pageNo, pageSize)
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

	// Get system switches command
	getSwitchesCmd := &cobra.Command{
		Use:   "get-switches",
		Short: "Get system switches",
		Long:  `Get system switches from Nacos server.`,
		Run: func(cmd *cobra.Command, args []string) {
			client := api.NewClient(viper.GetString("server"), viper.GetString("username"), viper.GetString("password"))
			result, err := client.GetSystemSwitches()
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

	// Update system switch command
	updateSwitchCmd := &cobra.Command{
		Use:   "update-switch",
		Short: "Update system switch",
		Long:  `Update system switch in Nacos server.`,
		Run: func(cmd *cobra.Command, args []string) {
			client := api.NewClient(viper.GetString("server"), viper.GetString("username"), viper.GetString("password"))
			success, err := client.UpdateSystemSwitch(switchEntry, switchValue)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				os.Exit(1)
			}

			if success {
				fmt.Println("System switch updated successfully")
			} else {
				fmt.Println("Failed to update system switch")
				os.Exit(1)
			}
		},
	}

	// Get metrics command
	getMetricsCmd := &cobra.Command{
		Use:   "metrics",
		Short: "Get system metrics",
		Long:  `Get system metrics from Nacos server.`,
		Run: func(cmd *cobra.Command, args []string) {
			client := api.NewClient(viper.GetString("server"), viper.GetString("username"), viper.GetString("password"))
			result, err := client.GetSystemMetrics()
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

	// Get server list command
	getServersCmd := &cobra.Command{
		Use:   "servers",
		Short: "Get server list",
		Long:  `Get server list from Nacos server.`,
		Run: func(cmd *cobra.Command, args []string) {
			client := api.NewClient(viper.GetString("server"), viper.GetString("username"), viper.GetString("password"))
			result, err := client.GetServerList()
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

	// Get leader command
	getLeaderCmd := &cobra.Command{
		Use:   "leader",
		Short: "Get cluster leader",
		Long:  `Get cluster leader from Nacos server.`,
		Run: func(cmd *cobra.Command, args []string) {
			client := api.NewClient(viper.GetString("server"), viper.GetString("username"), viper.GetString("password"))
			result, err := client.GetClusterLeader()
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

	// Update health command
	updateHealthCmd := &cobra.Command{
		Use:   "update-health",
		Short: "Update instance health",
		Long:  `Update instance health status in Nacos server.`,
		Run: func(cmd *cobra.Command, args []string) {
			client := api.NewClient(viper.GetString("server"), viper.GetString("username"), viper.GetString("password"))
			success, err := client.UpdateInstanceHealth(serviceName, ip, port, viper.GetString("namespace"), clusterName, healthy)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				os.Exit(1)
			}

			if success {
				fmt.Println("Instance health updated successfully")
			} else {
				fmt.Println("Failed to update instance health")
				os.Exit(1)
			}
		},
	}

	// Add flags for modify instance command
	modifyInstanceCmd.Flags().StringVar(&serviceName, "name", "", "Service name (required)")
	modifyInstanceCmd.Flags().StringVar(&ip, "ip", "", "IP address (required)")
	modifyInstanceCmd.Flags().StringVar(&port, "port", "", "Port (required)")
	modifyInstanceCmd.Flags().StringVar(&clusterName, "cluster", "", "Cluster name")
	modifyInstanceCmd.Flags().Float64Var(&weight, "weight", 1.0, "Instance weight")
	modifyInstanceCmd.Flags().StringVar(&metadata, "metadata", "", "Metadata in JSON format")
	modifyInstanceCmd.Flags().BoolVar(&enabled, "enabled", true, "Whether instance is enabled")
	modifyInstanceCmd.MarkFlagRequired("name")
	modifyInstanceCmd.MarkFlagRequired("ip")
	modifyInstanceCmd.MarkFlagRequired("port")

	// Add flags for list instances command
	listInstancesCmd.Flags().StringVar(&serviceName, "name", "", "Service name (required)")
	listInstancesCmd.Flags().StringVar(&clusters, "clusters", "", "Cluster names, comma separated")
	listInstancesCmd.Flags().StringVar(&healthyOnly, "healthy-only", "", "Only return healthy instances")
	listInstancesCmd.Flags().BoolVar(&prettyFormat, "pretty", false, "Pretty format JSON output")
	listInstancesCmd.MarkFlagRequired("name")

	// Add flags for get instance command
	getInstanceCmd.Flags().StringVar(&serviceName, "name", "", "Service name (required)")
	getInstanceCmd.Flags().StringVar(&ip, "ip", "", "IP address (required)")
	getInstanceCmd.Flags().StringVar(&port, "port", "", "Port (required)")
	getInstanceCmd.Flags().StringVar(&clusterName, "cluster", "", "Cluster name")
	getInstanceCmd.Flags().BoolVar(&prettyFormat, "pretty", false, "Pretty format JSON output")
	getInstanceCmd.MarkFlagRequired("name")
	getInstanceCmd.MarkFlagRequired("ip")
	getInstanceCmd.MarkFlagRequired("port")

	// Add flags for beat command
	beatCmd.Flags().StringVar(&serviceName, "name", "", "Service name (required)")
	beatCmd.Flags().StringVar(&ip, "ip", "", "IP address (required)")
	beatCmd.Flags().StringVar(&port, "port", "", "Port (required)")
	beatCmd.Flags().StringVar(&clusterName, "cluster", "", "Cluster name")
	beatCmd.Flags().Float64Var(&weight, "weight", 1.0, "Instance weight")
	beatCmd.Flags().StringVar(&metadata, "metadata", "", "Metadata in JSON format")
	beatCmd.MarkFlagRequired("name")
	beatCmd.MarkFlagRequired("ip")
	beatCmd.MarkFlagRequired("port")

	// Add flags for create service command
	createServiceCmd.Flags().StringVar(&serviceName, "name", "", "Service name (required)")
	createServiceCmd.Flags().StringVar(&groupName, "group", "", "Group name")
	createServiceCmd.Flags().StringVar(&protectThreshold, "protect", "0", "Protect threshold")
	createServiceCmd.Flags().StringVar(&metadata, "metadata", "", "Metadata in JSON format")
	createServiceCmd.MarkFlagRequired("name")

	// Add flags for delete service command
	deleteServiceCmd.Flags().StringVar(&serviceName, "name", "", "Service name (required)")
	deleteServiceCmd.Flags().StringVar(&groupName, "group", "", "Group name")
	deleteServiceCmd.MarkFlagRequired("name")

	// Add flags for update service command
	updateServiceCmd.Flags().StringVar(&serviceName, "name", "", "Service name (required)")
	updateServiceCmd.Flags().StringVar(&groupName, "group", "", "Group name")
	updateServiceCmd.Flags().StringVar(&protectThreshold, "protect", "", "Protect threshold")
	updateServiceCmd.Flags().StringVar(&metadata, "metadata", "", "Metadata in JSON format")
	updateServiceCmd.MarkFlagRequired("name")

	// Add flags for list services command
	listServicesCmd.Flags().StringVar(&groupName, "group", "", "Group name")
	listServicesCmd.Flags().IntVar(&pageNo, "page", 1, "Page number")
	listServicesCmd.Flags().IntVar(&pageSize, "size", 10, "Page size")
	listServicesCmd.Flags().BoolVar(&prettyFormat, "pretty", false, "Pretty format JSON output")

	// Add flags for update switch command
	updateSwitchCmd.Flags().StringVar(&switchEntry, "entry", "", "Switch entry (required)")
	updateSwitchCmd.Flags().StringVar(&switchValue, "value", "", "Switch value (required)")
	updateSwitchCmd.MarkFlagRequired("entry")
	updateSwitchCmd.MarkFlagRequired("value")

	// Add flags for get switches, metrics, servers, leader commands
	getSwitchesCmd.Flags().BoolVar(&prettyFormat, "pretty", false, "Pretty format JSON output")
	getMetricsCmd.Flags().BoolVar(&prettyFormat, "pretty", false, "Pretty format JSON output")
	getServersCmd.Flags().BoolVar(&prettyFormat, "pretty", false, "Pretty format JSON output")
	getLeaderCmd.Flags().BoolVar(&prettyFormat, "pretty", false, "Pretty format JSON output")

	// Add flags for update health command
	updateHealthCmd.Flags().StringVar(&serviceName, "name", "", "Service name (required)")
	updateHealthCmd.Flags().StringVar(&ip, "ip", "", "IP address (required)")
	updateHealthCmd.Flags().StringVar(&port, "port", "", "Port (required)")
	updateHealthCmd.Flags().StringVar(&clusterName, "cluster", "", "Cluster name")
	updateHealthCmd.Flags().BoolVar(&healthy, "healthy", true, "Health status")
	updateHealthCmd.MarkFlagRequired("name")
	updateHealthCmd.MarkFlagRequired("ip")
	updateHealthCmd.MarkFlagRequired("port")

	// Add commands to service command
	serviceCmd.AddCommand(modifyInstanceCmd)
	serviceCmd.AddCommand(listInstancesCmd)
	serviceCmd.AddCommand(getInstanceCmd)
	serviceCmd.AddCommand(beatCmd)
	serviceCmd.AddCommand(createServiceCmd)
	serviceCmd.AddCommand(deleteServiceCmd)
	serviceCmd.AddCommand(updateServiceCmd)
	serviceCmd.AddCommand(listServicesCmd)
	serviceCmd.AddCommand(getSwitchesCmd)
	serviceCmd.AddCommand(updateSwitchCmd)
	serviceCmd.AddCommand(getMetricsCmd)
	serviceCmd.AddCommand(getServersCmd)
	serviceCmd.AddCommand(getLeaderCmd)
	serviceCmd.AddCommand(updateHealthCmd)
}
