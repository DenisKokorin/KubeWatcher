package e2e

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

type loginResp struct {
	AccessToken  string `json:"AccessToken"`
	RefreshToken string `json:"RefreshToken"`
}

func TestRegisterLoginRefreshLogout(t *testing.T) {
	// ensure clean database state
	if err := ResetDatabase(); err != nil {
		t.Fatalf("reset db failed: %v", err)
	}
	defer ResetDatabase()

	client := newCookieClient()

	email := makeUniqueEmail("e2euser")
	// Register
	regBody := map[string]string{
		"name":     "E2E User",
		"email":    email,
		"position": "Dev",
		"team":     "E2E",
		"password": "P@ssw0rd",
	}
	res, err := postJSON(client, "/api/v1/auth/register", regBody)
	if err != nil {
		t.Fatalf("register request failed: %v", err)
	}
	if res.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(res.Body)
		t.Fatalf("register unexpected status: %d body: %s", res.StatusCode, string(body))
	}

	// Login
	loginBody := map[string]string{"email": email, "password": "P@ssw0rd"}
	res, err = postJSON(client, "/api/v1/auth/login", loginBody)
	if err != nil {
		t.Fatalf("login request failed: %v", err)
	}
	if res.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(res.Body)
		t.Fatalf("login unexpected status: %d body: %s", res.StatusCode, string(body))
	}
	var lr loginResp
	_ = json.NewDecoder(res.Body).Decode(&lr)
	if lr.AccessToken == "" || lr.RefreshToken == "" {
		t.Fatalf("login did not return tokens")
	}

	// Call /auth/me with Authorization header
	res, err = getWithAuth(client, "/api/v1/auth/me", lr.AccessToken)
	if err != nil {
		t.Fatalf("me request failed: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("me unexpected status: %d", res.StatusCode)
	}

	// Refresh tokens using cookies (client jar holds cookies)
	res, err = postWithAuth(client, "/api/v1/auth/refresh", "", nil)
	if err != nil {
		t.Fatalf("refresh request failed: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("refresh unexpected status: %d", res.StatusCode)
	}
	var rr loginResp
	_ = json.NewDecoder(res.Body).Decode(&rr)
	if rr.AccessToken == "" || rr.RefreshToken == "" {
		t.Fatalf("refresh did not return tokens")
	}

	// Logout
	res, err = postWithAuth(client, "/api/v1/auth/logout", "", nil)
	if err != nil {
		t.Fatalf("logout request failed: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("logout unexpected status: %d", res.StatusCode)
	}

	// After logout, calling /auth/me with old token may or may not be allowed depending on invalidation.
	// But refresh should fail because cookies were cleared.
	res, err = postWithAuth(client, "/api/v1/auth/refresh", "", nil)
	if err != nil {
		t.Fatalf("refresh after logout failed: %v", err)
	}
	if res.StatusCode == http.StatusOK {
		t.Fatalf("expected refresh after logout to fail, got 200")
	}
}
