package e2e

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestReviewsFilteringPagination(t *testing.T) {
	if err := ResetDatabase(); err != nil {
		t.Fatalf("reset db failed: %v", err)
	}
	defer ResetDatabase()

	// Create admin to be able to create employees
	adminClient := newCookieClient()
	adminEmail := makeUniqueEmail("admin")
	_, _ = postJSON(adminClient, "/api/v1/auth/register", map[string]string{"name": "Admin", "email": adminEmail, "position": "Lead", "team": "Ops", "password": "AdminPass", "role": "admin"})
	res, _ := postJSON(adminClient, "/api/v1/auth/login", map[string]string{"email": adminEmail, "password": "AdminPass"})
	var adminLR struct {
		AccessToken string `json:"AccessToken"`
	}
	_ = json.NewDecoder(res.Body).Decode(&adminLR)

	// Create two employees to attach reviews (using admin token)
	emp1Email := makeUniqueEmail("emp1")
	emp2Email := makeUniqueEmail("emp2")
	res, _ = postWithAuth(adminClient, "/api/v1/admin/users", adminLR.AccessToken, map[string]string{"name": "E1", "email": emp1Email, "position": "P", "team": "X", "password": "p", "role": "user"})
	var uid1 struct {
		UID string `json:"uid"`
	}
	_ = json.NewDecoder(res.Body).Decode(&uid1)
	res, _ = postWithAuth(adminClient, "/api/v1/admin/users", adminLR.AccessToken, map[string]string{"name": "E2", "email": emp2Email, "position": "P", "team": "X", "password": "p", "role": "user"})
	var uid2 struct {
		UID string `json:"uid"`
	}
	_ = json.NewDecoder(res.Body).Decode(&uid2)

	// Create a reviewer (normal user) and login
	reviewerClient := newCookieClient()
	reviewerEmail := makeUniqueEmail("reviewer")
	_, _ = postJSON(reviewerClient, "/api/v1/auth/register", map[string]string{"name": "R", "email": reviewerEmail, "position": "Dev", "team": "T", "password": "Pass"})
	res, _ = postJSON(reviewerClient, "/api/v1/auth/login", map[string]string{"email": reviewerEmail, "password": "Pass"})
	var lr struct {
		AccessToken string `json:"AccessToken"`
	}
	_ = json.NewDecoder(res.Body).Decode(&lr)

	// Create reviews for both employees using reviewer token
	for i := 0; i < 5; i++ {
		_, _ = postWithAuth(reviewerClient, "/api/v1/reviews", lr.AccessToken, map[string]interface{}{"employee_id": uid1.UID, "reviewer_id": uid2.UID, "period": "2025", "rating": 4})
	}
	for i := 0; i < 3; i++ {
		_, _ = postWithAuth(reviewerClient, "/api/v1/reviews", lr.AccessToken, map[string]interface{}{"employee_id": uid2.UID, "reviewer_id": uid1.UID, "period": "2025", "rating": 3})
	}

	// Query employees with pagination
	url := queryURL("/api/v1/reviews/employees", map[string]string{"page": "1", "size": "2"})
	res, err := adminClient.Get(url)
	if err != nil {
		t.Fatalf("list request failed: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		t.Fatalf("list bad status: %d %s", res.StatusCode, string(body))
	}
	var resp struct {
		Success bool        `json:"success"`
		Data    interface{} `json:"data"`
	}
	_ = json.NewDecoder(res.Body).Decode(&resp)
	if !resp.Success {
		t.Fatalf("list not successful")
	}
}
