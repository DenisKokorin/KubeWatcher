package e2e

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"
)

func TestUploadAndGetDocument(t *testing.T) {
	if err := ResetDatabase(); err != nil {
		t.Fatalf("reset db failed: %v", err)
	}
	defer ResetDatabase()

	adminClient := newCookieClient()
	adminEmail := makeUniqueEmail("adminfile")
	_, _ = postJSON(adminClient, "/api/v1/auth/register", map[string]string{"name": "AdminF", "email": adminEmail, "position": "Lead", "team": "Ops", "password": "AdminPass", "role": "admin"})
	res, _ := postJSON(adminClient, "/api/v1/auth/login", map[string]string{"email": adminEmail, "password": "AdminPass"})
	var adminLR struct {
		AccessToken string `json:"AccessToken"`
	}
	_ = json.NewDecoder(res.Body).Decode(&adminLR)

	// Create a user to attach document
	userEmail := makeUniqueEmail("docuser")
	res, _ = postWithAuth(adminClient, "/api/v1/admin/users", adminLR.AccessToken, map[string]string{"name": "DocUser", "email": userEmail, "position": "P", "team": "X", "password": "p", "role": "user"})
	var uid struct {
		UID string `json:"uid"`
	}
	_ = json.NewDecoder(res.Body).Decode(&uid)

	// Prepare a temp file
	tmp := os.TempDir() + string(os.PathSeparator) + "e2e-sample.txt"
	_ = os.WriteFile(tmp, []byte("hello e2e"), 0644)
	defer os.Remove(tmp)

	// Upload
	uploadPath := "/api/v1/admin/users/" + uid.UID + "/document"
	res, err := uploadFile(adminClient, uploadPath, adminLR.AccessToken, "file", tmp)
	if err != nil {
		t.Fatalf("upload failed: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		t.Fatalf("upload bad status: %d %s", res.StatusCode, string(body))
	}

	// Get presigned URL
	res, err = getWithAuth(adminClient, "/api/v1/admin/users/"+uid.UID+"/document", adminLR.AccessToken)
	if err != nil {
		t.Fatalf("get doc url failed: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		t.Fatalf("get doc url bad status: %d %s", res.StatusCode, string(body))
	}
	var resp struct {
		Success bool `json:"success"`
		Data    struct {
			URL string `json:"url"`
		} `json:"data"`
	}
	_ = json.NewDecoder(res.Body).Decode(&resp)
	if !resp.Success || resp.Data.URL == "" {
		t.Fatalf("presigned url missing")
	}
}
