// Package client provides an HTTP client for the Manuals API.
package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	// APIVersion is the API version to use.
	APIVersion = "2025.12"
)

// Client is an HTTP client for the Manuals API.
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// New creates a new API client.
func New(baseURL, apiKey string) *Client {
	return &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SearchResult represents a search result.
type SearchResult struct {
	DeviceID string  `json:"device_id"`
	Name     string  `json:"name"`
	Domain   string  `json:"domain"`
	Type     string  `json:"type"`
	Path     string  `json:"path"`
	Score    float64 `json:"score"`
	Snippet  string  `json:"snippet"`
}

// SearchResponse is the response from the search endpoint.
type SearchResponse struct {
	Results []SearchResult `json:"results"`
	Total   int            `json:"total"`
	Query   string         `json:"query"`
}

// Device represents a device.
type Device struct {
	ID        string                 `json:"id"`
	Domain    string                 `json:"domain"`
	Type      string                 `json:"type"`
	Name      string                 `json:"name"`
	Path      string                 `json:"path"`
	Content   string                 `json:"content,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	IndexedAt string                 `json:"indexed_at"`
}

// DevicesResponse is the response from the devices list endpoint.
type DevicesResponse struct {
	Data   []Device `json:"data"`
	Total  int      `json:"total"`
	Limit  int      `json:"limit"`
	Offset int      `json:"offset"`
}

// Document represents a document.
type Document struct {
	ID        string `json:"id"`
	DeviceID  string `json:"device_id"`
	Path      string `json:"path"`
	Filename  string `json:"filename"`
	MimeType  string `json:"mime_type"`
	SizeBytes int64  `json:"size_bytes"`
	Checksum  string `json:"checksum"`
	IndexedAt string `json:"indexed_at"`
}

// DocumentsResponse is the response from the documents list endpoint.
type DocumentsResponse struct {
	Data   []Document `json:"data"`
	Total  int        `json:"total"`
	Limit  int        `json:"limit"`
	Offset int        `json:"offset"`
}

// ErrorResponse is an API error response.
type ErrorResponse struct {
	Error string `json:"error"`
}

// Search searches for devices.
func (c *Client) Search(query string, limit int) (*SearchResponse, error) {
	params := url.Values{}
	params.Set("q", query)
	if limit > 0 {
		params.Set("limit", fmt.Sprintf("%d", limit))
	}

	var resp SearchResponse
	if err := c.get("/search?"+params.Encode(), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListDevices lists devices with pagination.
func (c *Client) ListDevices(limit, offset int, domain, deviceType string) (*DevicesResponse, error) {
	params := url.Values{}
	if limit > 0 {
		params.Set("limit", fmt.Sprintf("%d", limit))
	}
	if offset > 0 {
		params.Set("offset", fmt.Sprintf("%d", offset))
	}
	if domain != "" {
		params.Set("domain", domain)
	}
	if deviceType != "" {
		params.Set("type", deviceType)
	}

	path := "/devices"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	var resp DevicesResponse
	if err := c.get(path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetDevice gets a device by ID.
func (c *Client) GetDevice(id string) (*Device, error) {
	var resp Device
	if err := c.get("/devices/"+id, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListDocuments lists documents with pagination.
func (c *Client) ListDocuments(limit, offset int, deviceID string) (*DocumentsResponse, error) {
	params := url.Values{}
	if limit > 0 {
		params.Set("limit", fmt.Sprintf("%d", limit))
	}
	if offset > 0 {
		params.Set("offset", fmt.Sprintf("%d", offset))
	}
	if deviceID != "" {
		params.Set("device_id", deviceID)
	}

	path := "/documents"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	var resp DocumentsResponse
	if err := c.get(path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetDocument gets a document by ID.
func (c *Client) GetDocument(id string) (*Document, error) {
	var resp Document
	if err := c.get("/documents/"+id, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DownloadDocument downloads a document and returns the content.
func (c *Client) DownloadDocument(id string) (io.ReadCloser, string, error) {
	req, err := http.NewRequest("GET", c.baseURL+"/api/"+APIVersion+"/documents/"+id+"/download", nil)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return nil, "", fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}

	// Get filename from Content-Disposition header
	filename := ""
	if cd := resp.Header.Get("Content-Disposition"); cd != "" {
		// Simple parsing - look for filename=
		if idx := len("attachment; filename="); len(cd) > idx {
			filename = cd[idx:]
			filename = filename[1 : len(filename)-1] // Remove quotes
		}
	}

	return resp.Body, filename, nil
}

// get performs a GET request and decodes the JSON response.
func (c *Client) get(path string, result interface{}) error {
	req, err := http.NewRequest("GET", c.baseURL+"/api/"+APIVersion+path, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
		}
		return fmt.Errorf("API error (%d): %s", resp.StatusCode, errResp.Error)
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}
