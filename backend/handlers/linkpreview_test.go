package handlers

import (
	"net/url"
	"testing"

	"prep-dashboard/backend/models"
)

func TestParseMeta(t *testing.T) {
	body := `
	<html><head>
		<title>Two Sum - LeetCode</title>
		<meta property="og:title" content="Two Sum" />
		<meta property="og:description" content="Given an array of integers..." />
		<meta property="og:image" content="/static/cover.png" />
		<meta property="og:site_name" content="LeetCode" />
		<link rel="icon" href="/favicon.ico" />
	</head><body>hi</body></html>`

	base, _ := url.Parse("https://leetcode.com/problems/two-sum/")
	p := models.LinkPreview{URL: base.String(), SiteName: base.Hostname()}
	parseMeta(body, &p, base)

	if p.Title != "Two Sum" {
		t.Errorf("title = %q, want Two Sum", p.Title)
	}
	if p.Description != "Given an array of integers..." {
		t.Errorf("description = %q", p.Description)
	}
	if p.Image != "https://leetcode.com/static/cover.png" {
		t.Errorf("image not resolved to absolute: %q", p.Image)
	}
	if p.SiteName != "LeetCode" {
		t.Errorf("siteName = %q, want LeetCode", p.SiteName)
	}
	if p.Favicon != "https://leetcode.com/favicon.ico" {
		t.Errorf("favicon = %q", p.Favicon)
	}
}

func TestParseMetaFallbackTitle(t *testing.T) {
	body := `<html><head><title>Plain Page</title></head></html>`
	base, _ := url.Parse("https://example.com/x")
	p := models.LinkPreview{URL: base.String()}
	parseMeta(body, &p, base)
	if p.Title != "Plain Page" {
		t.Errorf("title = %q, want Plain Page", p.Title)
	}
}
