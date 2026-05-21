package external

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"

	"k8s-mon/internal/models"

	"golang.org/x/time/rate"
)

type ExternalAPI struct {
	client     *http.Client
	baseURL    string
	limiter    *rate.Limiter
	retryCount int
	source     string
}

type jsonPlaceholderPost struct {
	UserID int    `json:"userId"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func NewExternalAPI(baseURL string, timeout time.Duration, retryCount int, ratePerSecond float64) (*ExternalAPI, error) {
	if strings.TrimSpace(baseURL) == "" {
		return nil, errors.New("external API URL is required")
	}
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	if retryCount < 0 {
		retryCount = 2
	}
	if ratePerSecond <= 0 {
		ratePerSecond = 5
	}
	return &ExternalAPI{
		client:     &http.Client{Timeout: timeout},
		baseURL:    strings.TrimRight(baseURL, "/"),
		limiter:    rate.NewLimiter(rate.Limit(ratePerSecond), int(math.Max(1, math.Ceil(ratePerSecond)))),
		retryCount: retryCount,
		source:     "jsonplaceholder",
	}, nil
}

func (api *ExternalAPI) FetchPosts(ctx context.Context, userID, limit int, search string) ([]models.ExternalPost, error) {
	endpoint, err := api.buildURL(userID, limit)
	if err != nil {
		return nil, err
	}
	var result []jsonPlaceholderPost
	if err := api.doGet(ctx, endpoint, &result); err != nil {
		return nil, err
	}
	posts := make([]models.ExternalPost, 0, len(result))
	for _, item := range result {
		post := models.ExternalPost{
			ID:          item.ID,
			UserID:      item.UserID,
			Title:       item.Title,
			Summary:     summarize(item.Body),
			Content:     item.Body,
			Source:      api.source,
			PublishedAt: time.Now().UTC(),
		}
		if search == "" || strings.Contains(strings.ToLower(post.Title), strings.ToLower(search)) || strings.Contains(strings.ToLower(post.Content), strings.ToLower(search)) {
			posts = append(posts, post)
		}
	}
	return posts, nil
}

func (api *ExternalAPI) buildURL(userID, limit int) (string, error) {
	parsed, err := url.Parse(api.baseURL)
	if err != nil {
		return "", err
	}
	parsed.Path = strings.TrimRight(parsed.Path, "/") + "/posts"
	query := url.Values{}
	if userID > 0 {
		query.Set("userId", fmt.Sprintf("%d", userID))
	}
	if limit > 0 {
		query.Set("_limit", fmt.Sprintf("%d", limit))
	}
	parsed.RawQuery = query.Encode()
	return parsed.String(), nil
}

func (api *ExternalAPI) doGet(ctx context.Context, endpoint string, target interface{}) error {
	var lastErr error
	for attempt := 0; attempt <= api.retryCount; attempt++ {
		if err := api.limiter.Wait(ctx); err != nil {
			return err
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
		if err != nil {
			return err
		}
		resp, err := api.client.Do(req)
		if err != nil {
			lastErr = err
			if attempt < api.retryCount {
				time.Sleep(time.Duration(attempt+1) * 250 * time.Millisecond)
				continue
			}
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
			return fmt.Errorf("external API returned status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
		}
		if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
			return err
		}
		return nil
	}
	return lastErr
}

func summarize(text string) string {
	const maxLen = 120
	text = strings.TrimSpace(text)
	if len(text) <= maxLen {
		return text
	}
	return strings.TrimSpace(text[:maxLen]) + "..."
}
