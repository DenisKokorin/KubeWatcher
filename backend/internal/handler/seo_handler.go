package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type SEOHandler struct {
	baseURL string
}

func NewSEOHandler(baseURL string) *SEOHandler {
	baseURL = strings.TrimSpace(baseURL)
	if baseURL == "" {
		baseURL = "http://localhost:3000"
	}
	baseURL = strings.TrimRight(baseURL, "/")
	return &SEOHandler{baseURL: baseURL}
}

func (s *SEOHandler) Robots(c echo.Context) error {
	content := strings.Builder{}
	content.WriteString("User-agent: *\n")
	content.WriteString("Allow: /\n")
	content.WriteString("Disallow: /api/\n")
	content.WriteString("Disallow: /admin\n")
	content.WriteString("Disallow: /dashboard\n")
	content.WriteString("Disallow: /nodes\n")
	content.WriteString("Disallow: /pods\n")
	content.WriteString("Disallow: /employees\n")
	content.WriteString("Disallow: /reviews\n")
	content.WriteString("Crawl-delay: 10\n")
	content.WriteString(fmt.Sprintf("Sitemap: %s/sitemap.xml\n", s.baseURL))

	return c.Blob(http.StatusOK, "text/plain; charset=utf-8", []byte(content.String()))
}

func (s *SEOHandler) Sitemap(c echo.Context) error {
	xml := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url>
    <loc>%s/</loc>
    <changefreq>weekly</changefreq>
    <priority>0.8</priority>
  </url>
</urlset>`, s.baseURL)

	return c.Blob(http.StatusOK, "application/xml; charset=utf-8", []byte(xml))
}
