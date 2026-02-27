package feed

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/mmcdole/gofeed"
)

type Info struct {
	Title   string
	SiteURL string
}

func Fetch(ctx context.Context, feedURL string) (*Info, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("invalid feed URL: %w", err)
	}

	req.Header.Set("User-Agent", "GoSift/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to reach feed URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("feed URL returned status %d", resp.StatusCode)
	}

	parser := gofeed.NewParser()
	parsed, err := parser.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("URL does not contain a valid RSS/Atom feed: %w", err)
	}

	return &Info{
		Title:   parsed.Title,
		SiteURL: parsed.Link,
	}, nil
}
