package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

var (
	baseURL = func() string {
		if v := os.Getenv("BACKEND_BASE_URL"); v != "" {
			return v
		}
		return "http://localhost:7979"
	}()
)

func newCookieClient() *http.Client {
	jar, _ := cookiejar.New(nil)
	return &http.Client{Jar: jar, Timeout: 30 * time.Second}
}

func postJSON(client *http.Client, path string, body interface{}) (*http.Response, error) {
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, baseURL+path, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	return client.Do(req)
}

func getWithAuth(client *http.Client, path, accessToken string) (*http.Response, error) {
	req, _ := http.NewRequest(http.MethodGet, baseURL+path, nil)
	if accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}
	return client.Do(req)
}

func postWithAuth(client *http.Client, path, accessToken string, body interface{}) (*http.Response, error) {
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, baseURL+path, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	if accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}
	return client.Do(req)
}

func putWithAuth(client *http.Client, path, accessToken string, body interface{}) (*http.Response, error) {
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPut, baseURL+path, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	if accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}
	return client.Do(req)
}

func deleteWithAuth(client *http.Client, path, accessToken string) (*http.Response, error) {
	req, _ := http.NewRequest(http.MethodDelete, baseURL+path, nil)
	if accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}
	return client.Do(req)
}

func uploadFile(client *http.Client, path, accessToken, fieldName, filePath string) (*http.Response, error) {
	buf := &bytes.Buffer{}
	w := multipart.NewWriter(buf)
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	fw, err := w.CreateFormFile(fieldName, filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(fw, f); err != nil {
		return nil, err
	}
	w.Close()

	req, _ := http.NewRequest(http.MethodPost, baseURL+path, buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	if accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}
	return client.Do(req)
}

func makeUniqueEmail(prefix string) string {
	return fmt.Sprintf("%s-%d@example.com", prefix, time.Now().UnixNano())
}

func queryURL(path string, params map[string]string) string {
	u, _ := url.Parse(baseURL + path)
	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String()
}
