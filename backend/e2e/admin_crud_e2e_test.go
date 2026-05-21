package e2e

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

type createResp struct {
	UID string `json:"uid"`
}

func TestAdminCRUDAndRoleProtection(t *testing.T) {
	if err := ResetDatabase(); err != nil {
		t.Fatalf("reset db failed: %v", err)
	}
	defer ResetDatabase()

	// Create admin user
	adminClient := newCookieClient()
	adminEmail := makeUniqueEmail("admin")
	reg := map[string]string{"name": "Admin E2E", "email": adminEmail, "position": "Lead", "team": "Ops", "password": "AdminPass", "role": "admin"}
	res, err := postJSON(adminClient, "/api/v1/auth/register", reg)
	if err != nil {
		t.Fatalf("register admin failed: %v", err)
	}
	if res.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(res.Body)
		t.Fatalf("admin register bad: %d %s", res.StatusCode, string(body))
	}

	// Login admin
	res, err = postJSON(adminClient, "/api/v1/auth/login", map[string]string{"email": adminEmail, "password": "AdminPass"})
	if err != nil {
		t.Fatalf("admin login failed: %v", err)
	}
	if res.StatusCode != http.StatusAccepted {
		t.Fatalf("admin login bad status: %d", res.StatusCode)
	}
	var lr struct {
		AccessToken string `json:"AccessToken"`
	}
	_ = json.NewDecoder(res.Body).Decode(&lr)
	if lr.AccessToken == "" {
		t.Fatalf("no admin access token")
	}

	// Use admin to create a user via admin endpoint
	userBody := map[string]string{"name": "Created", "email": makeUniqueEmail("created"), "position": "Dev", "team": "A", "password": "UserPass", "role": "user"}
	res, err = postWithAuth(adminClient, "/api/v1/admin/users", lr.AccessToken, userBody)
	if err != nil {
		t.Fatalf("create user failed: %v", err)
	}
	if res.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(res.Body)
		t.Fatalf("create user bad: %d %s", res.StatusCode, string(body))
	}
	var cr createResp
	_ = json.NewDecoder(res.Body).Decode(&cr)
	if cr.UID == "" {
		t.Fatalf("create response missing uid")
	}

	// Get users list
	res, err = getWithAuth(adminClient, "/api/v1/admin/users", lr.AccessToken)
	if err != nil {
		t.Fatalf("get users failed: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("get users bad status: %d", res.StatusCode)
	}

	// Now create a normal user and ensure they cannot access admin routes
	userClient := newCookieClient()
	userEmail := makeUniqueEmail("user")
	_, _ = postJSON(userClient, "/api/v1/auth/register", map[string]string{"name": "User E2E", "email": userEmail, "position": "Dev", "team": "A", "password": "UserPass"})
	res, _ = postJSON(userClient, "/api/v1/auth/login", map[string]string{"email": userEmail, "password": "UserPass"})
	var ulr struct {
		AccessToken string `json:"AccessToken"`
	}
	_ = json.NewDecoder(res.Body).Decode(&ulr)

	// Attempt admin access
	res, err = getWithAuth(userClient, "/api/v1/admin/users", ulr.AccessToken)
	if err != nil {
		t.Fatalf("user get admin failed: %v", err)
	}
	if res.StatusCode != http.StatusForbidden && res.StatusCode != http.StatusUnauthorized {
		body, _ := io.ReadAll(res.Body)
		t.Fatalf("expected forbidden/unauthorized for non-admin, got %d body: %s", res.StatusCode, string(body))
	}
}
