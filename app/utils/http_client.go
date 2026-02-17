package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type HTTPClient struct {
	client  *http.Client
	baseURL string
}

func NewHTTPClient(baseURL string, timeout time.Duration) *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: timeout,
		},
		baseURL: baseURL,
	}
}

type MembershipCountResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		MembershipCount int    `json:"membership_count"`
		LastUpdatedAt   string `json:"last_updated_at"`
	} `json:"data"`
}

type MembershipCountData struct {
	Count         int
	LastUpdatedAt string
}

func (h *HTTPClient) GetMembershipCount() (*MembershipCountData, error) {
	url := fmt.Sprintf("%s/admin/customers/membership-count", h.baseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call customer service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("customer service returned status %d: %s", resp.StatusCode, string(body))
	}

	var result MembershipCountResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &MembershipCountData{
		Count:         result.Data.MembershipCount,
		LastUpdatedAt: result.Data.LastUpdatedAt,
	}, nil
}
