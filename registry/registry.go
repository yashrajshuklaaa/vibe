package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/vibefile-dev/vibe/config"
)

type SkillEntry struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	DownloadURL string `json:"download_url"`
}

type Client struct {
	cfg  config.RegistryConfig
	http *http.Client
}

func New(cfg *config.RegistryConfig) *Client {
	if cfg == nil || cfg.URL == "" {
		return nil
	}
	t := cfg.Timeout
	if t == 0 {
		t = 10
	}
	return &Client{cfg: *cfg, http: &http.Client{Timeout: time.Duration(t) * time.Second}}
}

func (c *Client) LookupSkill(ctx context.Context, name string) (*SkillEntry, error) {
	if c == nil {
		return nil, nil
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("%s/skills/%s", c.cfg.URL, name), nil)
	if err != nil {
		return nil, fmt.Errorf("registry: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("registry: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("registry: unexpected status %d for skill %q", resp.StatusCode, name)
	}

	var e SkillEntry
	if err := json.NewDecoder(resp.Body).Decode(&e); err != nil {
		return nil, fmt.Errorf("registry: decode: %w", err)
	}
	return &e, nil
}
