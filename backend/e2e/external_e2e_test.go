package e2e

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestExternalPostsEndpoint(t *testing.T) {
	client := newCookieClient()
	email := makeUniqueEmail("extuser")
	_, _ = postJSON(client, "/api/v1/auth/register", map[string]string{"name": "Ext", "email": email, "position": "Dev", "team": "Ext", "password": "Pass"})
	res, _ := postJSON(client, "/api/v1/auth/login", map[string]string{"email": email, "password": "Pass"})
	var lr struct {
		AccessToken string `json:"AccessToken"`
	}
	_ = json.NewDecoder(res.Body).Decode(&lr)

	// Call external posts
	req, _ := http.NewRequest(http.MethodGet, baseURL+"/api/v1/external/posts", nil)
	res, err := client.Do(req)
	if err != nil {
		t.Fatalf("external posts request failed: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("external posts bad status: %d", res.StatusCode)
	}
	// expect JSON array
	var posts []interface{}
	_ = json.NewDecoder(res.Body).Decode(&posts)
}
