package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// ModifyInstance modifies a service instance
func (c *Client) ModifyInstance(serviceName, ip, port, clusterName, namespaceId string, weight float64, metadata map[string]string, enabled bool) (bool, error) {
	params := url.Values{}
	params.Set("serviceName", serviceName)
	params.Set("ip", ip)
	params.Set("port", port)
	if clusterName != "" {
		params.Set("clusterName", clusterName)
	}
	if namespaceId != "" {
		params.Set("namespaceId", namespaceId)
	}
	params.Set("weight", fmt.Sprintf("%f", weight))

	if len(metadata) > 0 {
		metadataBytes, err := json.Marshal(metadata)
		if err != nil {
			return false, err
		}
		params.Set("metadata", string(metadataBytes))
	}

	params.Set("enabled", fmt.Sprintf("%t", enabled))

	resp, err := c.doRequest(http.MethodPut, "/nacos/v1/ns/instance", params, []byte(params.Encode()))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("failed to modify instance: %s, status code: %d", string(body), resp.StatusCode)
	}

	return true, nil
}

// ListInstances lists all instances of a service
func (c *Client) ListInstances(serviceName, namespaceId, clusters, healthyOnly string) (string, error) {
	params := url.Values{}
	params.Set("serviceName", serviceName)

	if namespaceId != "" {
		params.Set("namespaceId", namespaceId)
	}
	if clusters != "" {
		params.Set("clusters", clusters)
	}
	if healthyOnly != "" {
		params.Set("healthyOnly", healthyOnly)
	}

	resp, err := c.doRequest(http.MethodGet, "/nacos/v1/ns/instance/list", params, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to list instances: %s, status code: %d", string(body), resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// GetInstance gets the details of a specific instance
func (c *Client) GetInstance(serviceName, ip, port, namespaceId, cluster string) (string, error) {
	params := url.Values{}
	params.Set("serviceName", serviceName)
	params.Set("ip", ip)
	params.Set("port", port)

	if namespaceId != "" {
		params.Set("namespaceId", namespaceId)
	}
	if cluster != "" {
		params.Set("cluster", cluster)
	}

	resp, err := c.doRequest(http.MethodGet, "/nacos/v1/ns/instance", params, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to get instance: %s, status code: %d", string(body), resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// SendInstanceHeartbeat sends a heartbeat for an instance
func (c *Client) SendInstanceHeartbeat(serviceName, ip, port, namespaceId, cluster string, weight float64, metadata map[string]string) (bool, error) {
	params := url.Values{}
	params.Set("serviceName", serviceName)
	params.Set("ip", ip)
	params.Set("port", port)
	params.Set("beat", "")

	if namespaceId != "" {
		params.Set("namespaceId", namespaceId)
	}
	if cluster != "" {
		params.Set("clusterName", cluster)
	}

	// Construct beat JSON
	beat := map[string]interface{}{
		"ip":          ip,
		"port":        port,
		"serviceName": serviceName,
	}

	if weight != 0 {
		beat["weight"] = weight
	}
	if cluster != "" {
		beat["cluster"] = cluster
	}
	if len(metadata) > 0 {
		beat["metadata"] = metadata
	}

	beatBytes, err := json.Marshal(beat)
	if err != nil {
		return false, err
	}
	params.Set("beat", string(beatBytes))

	resp, err := c.doRequest(http.MethodPut, "/nacos/v1/ns/instance/beat", params, []byte(params.Encode()))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("failed to send heartbeat: %s, status code: %d", string(body), resp.StatusCode)
	}

	return true, nil
}

// CreateService creates a new service
func (c *Client) CreateService(serviceName, namespaceId, groupName, protectThreshold string, metadata map[string]string) (bool, error) {
	params := url.Values{}
	params.Set("serviceName", serviceName)

	if namespaceId != "" {
		params.Set("namespaceId", namespaceId)
	}
	if groupName != "" {
		params.Set("groupName", groupName)
	}
	if protectThreshold != "" {
		params.Set("protectThreshold", protectThreshold)
	}

	if len(metadata) > 0 {
		metadataBytes, err := json.Marshal(metadata)
		if err != nil {
			return false, err
		}
		params.Set("metadata", string(metadataBytes))
	}

	resp, err := c.doRequest(http.MethodPost, "/nacos/v1/ns/service", params, []byte(params.Encode()))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("failed to create service: %s, status code: %d", string(body), resp.StatusCode)
	}

	return true, nil
}

// DeleteService deletes a service
func (c *Client) DeleteService(serviceName, namespaceId, groupName string) (bool, error) {
	params := url.Values{}
	params.Set("serviceName", serviceName)

	if namespaceId != "" {
		params.Set("namespaceId", namespaceId)
	}
	if groupName != "" {
		params.Set("groupName", groupName)
	}

	resp, err := c.doRequest(http.MethodDelete, "/nacos/v1/ns/service", params, nil)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("failed to delete service: %s, status code: %d", string(body), resp.StatusCode)
	}

	return true, nil
}

// UpdateService updates a service
func (c *Client) UpdateService(serviceName, namespaceId, groupName, protectThreshold string, metadata map[string]string) (bool, error) {
	params := url.Values{}
	params.Set("serviceName", serviceName)

	if namespaceId != "" {
		params.Set("namespaceId", namespaceId)
	}
	if groupName != "" {
		params.Set("groupName", groupName)
	}
	if protectThreshold != "" {
		params.Set("protectThreshold", protectThreshold)
	}

	if len(metadata) > 0 {
		metadataBytes, err := json.Marshal(metadata)
		if err != nil {
			return false, err
		}
		params.Set("metadata", string(metadataBytes))
	}

	resp, err := c.doRequest(http.MethodPut, "/nacos/v1/ns/service", params, []byte(params.Encode()))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("failed to update service: %s, status code: %d", string(body), resp.StatusCode)
	}

	return true, nil
}

// ListServices lists all services
func (c *Client) ListServices(namespaceId, groupName string, pageNo, pageSize int) (string, error) {
	params := url.Values{}

	if namespaceId != "" {
		params.Set("namespaceId", namespaceId)
	}
	if groupName != "" {
		params.Set("groupName", groupName)
	}

	params.Set("pageNo", fmt.Sprintf("%d", pageNo))
	params.Set("pageSize", fmt.Sprintf("%d", pageSize))

	resp, err := c.doRequest(http.MethodGet, "/nacos/v1/ns/service/list", params, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to list services: %s, status code: %d", string(body), resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// GetSystemSwitches gets system switches
func (c *Client) GetSystemSwitches() (string, error) {
	resp, err := c.doRequest(http.MethodGet, "/nacos/v1/ns/operator/switches", nil, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to get system switches: %s, status code: %d", string(body), resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// UpdateSystemSwitch updates a system switch
func (c *Client) UpdateSystemSwitch(entry, value string) (bool, error) {
	params := url.Values{}
	params.Set("entry", entry)
	params.Set("value", value)

	resp, err := c.doRequest(http.MethodPut, "/nacos/v1/ns/operator/switches", params, []byte(params.Encode()))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("failed to update system switch: %s, status code: %d", string(body), resp.StatusCode)
	}

	return true, nil
}

// GetSystemMetrics gets system metrics
func (c *Client) GetSystemMetrics() (string, error) {
	resp, err := c.doRequest(http.MethodGet, "/nacos/v1/ns/operator/metrics", nil, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to get system metrics: %s, status code: %d", string(body), resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// GetServerList gets the server list
func (c *Client) GetServerList() (string, error) {
	resp, err := c.doRequest(http.MethodGet, "/nacos/v1/ns/operator/servers", nil, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to get server list: %s, status code: %d", string(body), resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// GetClusterLeader gets the cluster leader
func (c *Client) GetClusterLeader() (string, error) {
	resp, err := c.doRequest(http.MethodGet, "/nacos/v1/ns/raft/leader", nil, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to get cluster leader: %s, status code: %d", string(body), resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// UpdateInstanceHealth updates the health status of an instance
func (c *Client) UpdateInstanceHealth(serviceName, ip, port, namespaceId, cluster string, healthy bool) (bool, error) {
	params := url.Values{}
	params.Set("serviceName", serviceName)
	params.Set("ip", ip)
	params.Set("port", port)

	if namespaceId != "" {
		params.Set("namespaceId", namespaceId)
	}
	if cluster != "" {
		params.Set("clusterName", cluster)
	}

	params.Set("healthy", fmt.Sprintf("%t", healthy))

	resp, err := c.doRequest(http.MethodPut, "/nacos/v1/ns/health/instance", params, []byte(params.Encode()))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("failed to update instance health: %s, status code: %d", string(body), resp.StatusCode)
	}

	return true, nil
}
