package githubstatus

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const defaultBaseURL = "https://www.githubstatus.com"

type Client struct {
	baseURL    string
	timeout    time.Duration
	httpClient *http.Client
}

func NewClient(baseURL string, timeout time.Duration) *Client {
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	baseURL = strings.TrimRight(baseURL, "/")
	if timeout <= 0 {
		timeout = 10 * time.Second
	}

	return &Client{
		baseURL: baseURL,
		timeout: timeout,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) Timeout() time.Duration {
	return c.timeout
}

func (c *Client) GetStatus(ctx context.Context) (*StatusResponse, error) {
	var resp StatusResponse
	if err := c.getJSON(ctx, "/api/v2/status.json", &resp); err != nil {
		return nil, fmt.Errorf("fetch status: %w", err)
	}
	return &resp, nil
}

func (c *Client) GetComponents(ctx context.Context) (*ComponentsResponse, error) {
	var resp ComponentsResponse
	if err := c.getJSON(ctx, "/api/v2/components.json", &resp); err != nil {
		return nil, fmt.Errorf("fetch components: %w", err)
	}
	return &resp, nil
}

func (c *Client) GetUnresolvedIncidents(ctx context.Context) (*IncidentsResponse, error) {
	var resp IncidentsResponse
	if err := c.getJSON(ctx, "/api/v2/incidents/unresolved.json", &resp); err != nil {
		return nil, fmt.Errorf("fetch unresolved incidents: %w", err)
	}
	return &resp, nil
}

func (c *Client) GetActiveMaintenances(ctx context.Context) (*ScheduledMaintenancesResponse, error) {
	var resp ScheduledMaintenancesResponse
	if err := c.getJSON(ctx, "/api/v2/scheduled-maintenances/active.json", &resp); err != nil {
		return nil, fmt.Errorf("fetch active maintenances: %w", err)
	}
	return &resp, nil
}

func (c *Client) getJSON(ctx context.Context, path string, dest any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "github-outage-exporter")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request %s: %w", path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("request %s: unexpected status %d: %s", path, resp.StatusCode, strings.TrimSpace(string(body)))
	}

	if err := json.NewDecoder(resp.Body).Decode(dest); err != nil {
		return fmt.Errorf("decode response from %s: %w", path, err)
	}

	return nil
}
