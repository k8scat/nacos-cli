package api

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// ListenConfig listens for configuration changes
func (c *Client) ListenConfig(dataId, group, contentMD5, tenant string) (string, error) {
	// Prepare the listening data
	var listenData string
	if tenant != "" {
		listenData = fmt.Sprintf("%s%s%s%s%s%s%s%s%s",
			dataId, string(rune(2)), group, string(rune(2)), contentMD5, string(rune(2)), tenant, string(rune(1)))
	} else {
		listenData = fmt.Sprintf("%s%s%s%s%s%s",
			dataId, string(rune(2)), group, string(rune(2)), contentMD5, string(rune(1)))
	}

	// URL encode the listening data
	params := url.Values{}
	params.Set("Listening-Configs", listenData)

	// Create request
	u, err := url.Parse(c.Server)
	if err != nil {
		return "", err
	}
	u.Path = "/nacos/v1/cs/configs/listener"

	req, err := http.NewRequest(http.MethodPost, u.String(), strings.NewReader(params.Encode()))
	if err != nil {
		return "", err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Long-Pulling-Timeout", "30000")

	// Execute request
	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to listen config: %s, status code: %d", string(body), resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// GetConfigHistory retrieves configuration history versions
func (c *Client) GetConfigHistory(dataId, group, tenant string, pageNo, pageSize int) (string, error) {
	params := url.Values{}
	params.Set("dataId", dataId)
	params.Set("group", group)
	params.Set("pageNo", fmt.Sprintf("%d", pageNo))
	params.Set("pageSize", fmt.Sprintf("%d", pageSize))
	if tenant != "" {
		params.Set("tenant", tenant)
	}

	resp, err := c.doRequest(http.MethodGet, "/nacos/v1/cs/history", params, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to get config history: %s, status code: %d", string(body), resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// GetConfigHistoryDetail retrieves a specific history version of a configuration
func (c *Client) GetConfigHistoryDetail(dataId, group, tenant, nid string) (string, error) {
	params := url.Values{}
	params.Set("dataId", dataId)
	params.Set("group", group)
	params.Set("nid", nid)
	if tenant != "" {
		params.Set("tenant", tenant)
	}

	resp, err := c.doRequest(http.MethodGet, "/nacos/v1/cs/history/detail", params, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to get config history detail: %s, status code: %d", string(body), resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// GetPreviousConfigInfo retrieves the previous version of a configuration
func (c *Client) GetPreviousConfigInfo(dataId, group, tenant string) (string, error) {
	params := url.Values{}
	params.Set("dataId", dataId)
	params.Set("group", group)
	if tenant != "" {
		params.Set("tenant", tenant)
	}

	resp, err := c.doRequest(http.MethodGet, "/nacos/v1/cs/history/previous", params, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to get previous config: %s, status code: %d", string(body), resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
