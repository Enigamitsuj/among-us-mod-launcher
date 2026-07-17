package githubapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Enigamitsuj/among-us-mod-launcher/internal/game"
	"github.com/Enigamitsuj/among-us-mod-launcher/internal/mods"
)

const apiBase = "https://api.github.com"

type releaseDTO struct {
	TagName     string    `json:"tag_name"`
	Name        string    `json:"name"`
	Body        string    `json:"body"`
	PublishedAt time.Time `json:"published_at"`
	Prerelease  bool      `json:"prerelease"`
	Draft       bool      `json:"draft"`
	HTMLURL     string    `json:"html_url"`
	Assets      []assetDTO `json:"assets"`
}

type assetDTO struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int64  `json:"size"`
	ContentType        string `json:"content_type"`
}

// Client talks to the GitHub Releases API.
type Client struct {
	HTTP *http.Client
}

func NewClient() *Client {
	return &Client{
		HTTP: &http.Client{Timeout: 20 * time.Second},
	}
}

// ListReleases fetches non-draft releases, picking the ZIP that matches platform.
func (c *Client) ListReleases(owner, repo string, platform game.Platform, limit int) ([]mods.Release, error) {
	if owner == "" || repo == "" {
		return nil, fmt.Errorf("GitHub repository is not configured for this mod")
	}
	if limit <= 0 {
		limit = 15
	}

	url := fmt.Sprintf("%s/repos/%s/%s/releases?per_page=%d", apiBase, owner, repo, limit)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "AmongUs-Mod-Launcher")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GitHub unavailable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("GitHub unavailable (%d): %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var raw []releaseDTO
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("failed to read GitHub releases")
	}

	hint := game.AssetPattern(platform)
	out := make([]mods.Release, 0, len(raw))
	for _, r := range raw {
		if r.Draft {
			continue
		}
		assetName, assetURL, assetSize, matched := pickZipAsset(r.Assets, hint)
		if assetURL == "" {
			continue
		}
		name := r.Name
		if name == "" {
			name = r.TagName
		}
		out = append(out, mods.Release{
			TagName:      r.TagName,
			Name:           name,
			Body:           r.Body,
			PublishedAt:    r.PublishedAt.Format(time.RFC3339),
			Prerelease:     r.Prerelease,
			DownloadURL:    assetURL,
			DownloadName:   assetName,
			DownloadSize:   assetSize,
			HTMLURL:        r.HTMLURL,
			AssetMatched:   matched,
			AssetHint:      hint,
		})
	}
	return out, nil
}

func pickZipAsset(assets []assetDTO, hint string) (name, url string, size int64, matched bool) {
	hint = strings.ToLower(strings.TrimSpace(hint))

	if hint != "" {
		for _, a := range assets {
			lower := strings.ToLower(a.Name)
			if strings.HasSuffix(lower, ".zip") && strings.Contains(lower, hint) {
				return a.Name, a.BrowserDownloadURL, a.Size, true
			}
		}
	}

	// Prefer any windows zip over mac/linux when hint missing/unmatched.
	for _, a := range assets {
		lower := strings.ToLower(a.Name)
		if !strings.HasSuffix(lower, ".zip") {
			continue
		}
		if strings.Contains(lower, "macos") || strings.Contains(lower, "linux") {
			continue
		}
		return a.Name, a.BrowserDownloadURL, a.Size, false
	}

	for _, a := range assets {
		lower := strings.ToLower(a.Name)
		if strings.HasSuffix(lower, ".zip") {
			return a.Name, a.BrowserDownloadURL, a.Size, false
		}
	}
	return "", "", 0, false
}
