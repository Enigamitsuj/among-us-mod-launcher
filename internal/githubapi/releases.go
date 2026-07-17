package githubapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Enigamitsuj/among-us-mod-launcher/internal/mods"
)

const apiBase = "https://api.github.com"

type releaseDTO struct {
	TagName    string    `json:"tag_name"`
	Name       string    `json:"name"`
	Body       string    `json:"body"`
	PublishedAt time.Time `json:"published_at"`
	Prerelease bool      `json:"prerelease"`
	Draft      bool      `json:"draft"`
	HTMLURL    string    `json:"html_url"`
	Assets     []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
		Size               int64  `json:"size"`
		ContentType        string `json:"content_type"`
	} `json:"assets"`
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

// ListReleases fetches non-draft releases for a mod's GitHub repo.
func (c *Client) ListReleases(owner, repo string, limit int) ([]mods.Release, error) {
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

	out := make([]mods.Release, 0, len(raw))
	for _, r := range raw {
		if r.Draft {
			continue
		}
		assetName, assetURL, assetSize := pickZipAsset(r.Assets)
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
		})
	}
	return out, nil
}

// Categorize splits releases into latest / stable / beta helpers for the UI.
func Categorize(releases []mods.Release) (latest, stable, beta *mods.Release) {
	for i := range releases {
		r := &releases[i]
		if latest == nil {
			latest = r
		}
		if r.Prerelease {
			if beta == nil {
				beta = r
			}
			continue
		}
		if stable == nil {
			stable = r
		}
	}
	return latest, stable, beta
}

func pickZipAsset(assets []struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int64  `json:"size"`
	ContentType        string `json:"content_type"`
}) (name, url string, size int64) {
	for _, a := range assets {
		lower := strings.ToLower(a.Name)
		if strings.HasSuffix(lower, ".zip") {
			return a.Name, a.BrowserDownloadURL, a.Size
		}
	}
	for _, a := range assets {
		lower := strings.ToLower(a.Name)
		if strings.Contains(lower, "zip") || strings.Contains(a.ContentType, "zip") {
			return a.Name, a.BrowserDownloadURL, a.Size
		}
	}
	return "", "", 0
}
