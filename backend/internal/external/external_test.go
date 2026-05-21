package external

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestExternalAPI_FetchPosts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]any{{
			"userId": 1,
			"id":     1,
			"title":  "Hello World",
			"body":   "Sample content from external API",
		}})
	}))
	defer server.Close()

	api, err := NewExternalAPI(server.URL, 5*time.Second, 1, 10)
	if err != nil {
		t.Fatal(err)
	}

	posts, err := api.FetchPosts(context.Background(), 1, 1, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(posts) != 1 {
		t.Fatalf("expected 1 post, got %d", len(posts))
	}
	if posts[0].Source != "jsonplaceholder" {
		t.Fatalf("expected source jsonplaceholder, got %q", posts[0].Source)
	}
}
