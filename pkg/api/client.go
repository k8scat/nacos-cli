package api

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client represents a Nacos API client
type Client struct {
	Server   string
	Username string
	Password string
	Timeout  time.Duration
	client   *http.Client
}

// NewClient creates a new Nacos API client
func NewClient(server, username, password string) *Client {
	return &Client{
		Server:   server,
		Username: username,
		Password: password,
		Timeout:  10 * time.Second,
		client:   &http.Client{Timeout: 10 * time.Second},
	}
}

// doRequest performs an HTTP request to the Nacos API
func (c *Client) doRequest(method, path string, params url.Values, body []byte) (*http.Response, error) {
	// Make a copy of params for query params and form params
	queryParams := url.Values{}

	// Add auth params to query string only
	if c.Username != "" && c.Password != "" {
		queryParams.Set("username", c.Username)
		queryParams.Set("password", c.Password)
	}

	// Build URL with only auth params in query string
	u, err := url.Parse(c.Server)
	if err != nil {
		return nil, err
	}
	u.Path = path
	u.RawQuery = queryParams.Encode()

	// Set request body
	var requestBody io.Reader
	if method == http.MethodPost || method == http.MethodPut {
		// For POST/PUT, use params as form body
		if body == nil {
			requestBody = strings.NewReader(params.Encode())
		} else {
			requestBody = bytes.NewReader(body)
		}
	} else {
		// For other methods, add params to query string and use empty body
		for k, v := range params {
			for _, val := range v {
				q := u.Query()
				q.Add(k, val)
				u.RawQuery = q.Encode()
			}
		}
		requestBody = bytes.NewReader(body)
	}

	// Create request
	req, err := http.NewRequest(method, u.String(), requestBody)
	if err != nil {
		return nil, err
	}

	// Set headers for form post
	if (method == http.MethodPost || method == http.MethodPut) && (body != nil || len(params) > 0) {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	// Execute request
	return c.client.Do(req)
}

// GetConfig retrieves a configuration from Nacos server
func (c *Client) GetConfig(dataId, group, tenant string) (string, error) {
	params := url.Values{}
	params.Set("dataId", dataId)
	params.Set("group", group)
	if tenant != "" {
		params.Set("tenant", tenant)
	}

	resp, err := c.doRequest(http.MethodGet, "/nacos/v1/cs/configs", params, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to get config: %s, status code: %d", string(body), resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// PublishConfig publishes a configuration to Nacos server
func (c *Client) PublishConfig(dataId, group, content, configType, tenant string) (bool, error) {
	params := url.Values{}
	params.Set("dataId", dataId)
	params.Set("group", group)
	params.Set("content", content)
	if configType != "" {
		params.Set("type", configType)
	}
	if tenant != "" {
		params.Set("tenant", tenant)
	}

	// Pass nil as body so doRequest will use params as the form body
	resp, err := c.doRequest(http.MethodPost, "/nacos/v1/cs/configs", params, nil)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("failed to publish config: %s, status code: %d", string(body), resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	return string(body) == "true", nil
}

// DeleteConfig deletes a configuration from Nacos server
func (c *Client) DeleteConfig(dataId, group, tenant string) (bool, error) {
	params := url.Values{}
	params.Set("dataId", dataId)
	params.Set("group", group)
	if tenant != "" {
		params.Set("tenant", tenant)
	}

	resp, err := c.doRequest(http.MethodDelete, "/nacos/v1/cs/configs", params, nil)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("failed to delete config: %s, status code: %d", string(body), resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	return string(body) == "true", nil
}

// GetMD5 calculates the MD5 hash of a string
func GetMD5(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

// RegisterInstance registers a service instance to Nacos server
func (c *Client) RegisterInstance(serviceName, ip, port, clusterName, namespaceId string, metadata map[string]string, ephemeral bool) (bool, error) {
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
	if len(metadata) > 0 {
		metadataBytes, err := json.Marshal(metadata)
		if err != nil {
			return false, err
		}
		params.Set("metadata", string(metadataBytes))
	}
	params.Set("ephemeral", fmt.Sprintf("%t", ephemeral))

	// Pass nil as body so doRequest will use params as the form body
	resp, err := c.doRequest(http.MethodPost, "/nacos/v1/ns/instance", params, nil)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("failed to register instance: %s, status code: %d", string(body), resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	return string(body) == "ok", nil
}

// DeregisterInstance deregisters a service instance from Nacos server
func (c *Client) DeregisterInstance(serviceName, ip, port, clusterName, namespaceId string) (bool, error) {
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

	resp, err := c.doRequest(http.MethodDelete, "/nacos/v1/ns/instance", params, nil)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("failed to deregister instance: %s, status code: %d", string(body), resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	return string(body) == "ok", nil
}

// GetService gets service information from Nacos server
func (c *Client) GetService(serviceName, namespaceId string) (string, error) {
	params := url.Values{}
	params.Set("serviceName", serviceName)
	if namespaceId != "" {
		params.Set("namespaceId", namespaceId)
	}

	resp, err := c.doRequest(http.MethodGet, "/nacos/v1/ns/service", params, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to get service: %s, status code: %d", string(body), resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// ListNamespaces lists all namespaces from Nacos server
func (c *Client) ListNamespaces() (string, error) {
	resp, err := c.doRequest(http.MethodGet, "/nacos/v1/console/namespaces", nil, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to list namespaces: %s, status code: %d", string(body), resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// CreateNamespace creates a new namespace in Nacos server
func (c *Client) CreateNamespace(namespaceId, namespaceName, namespaceDesc string) (bool, error) {
	params := url.Values{}
	params.Set("customNamespaceId", namespaceId)
	params.Set("namespaceName", namespaceName)
	if namespaceDesc != "" {
		params.Set("namespaceDesc", namespaceDesc)
	}

	resp, err := c.doRequest(http.MethodPost, "/nacos/v1/console/namespaces", params, nil)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("failed to create namespace: %s, status code: %d", string(body), resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	return string(body) == "true", nil
}

// ModifyNamespace modifies a namespace in Nacos server
func (c *Client) ModifyNamespace(namespaceId, namespaceName, namespaceDesc string) (bool, error) {
	params := url.Values{}
	params.Set("namespace", namespaceId)
	params.Set("namespaceShowName", namespaceName)
	params.Set("namespaceDesc", namespaceDesc)

	// Pass nil as body so doRequest will use params as the form body
	resp, err := c.doRequest(http.MethodPut, "/nacos/v1/console/namespaces", params, nil)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("failed to modify namespace: %s, status code: %d", string(body), resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	return string(body) == "true", nil
}

// DeleteNamespace deletes a namespace from Nacos server
func (c *Client) DeleteNamespace(namespaceId string) (bool, error) {
	params := url.Values{}
	params.Set("namespaceId", namespaceId)

	// Pass nil as body so doRequest will use params as the form body for DELETE
	resp, err := c.doRequest(http.MethodDelete, "/nacos/v1/console/namespaces", params, nil)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("failed to delete namespace: %s, status code: %d", string(body), resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	return string(body) == "true", nil
}
