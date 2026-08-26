package handlers

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/html"
	"prep-dashboard/backend/models"
)

// previewClient has a sane timeout so a slow site can't hang the server.
var previewClient = &http.Client{Timeout: 8 * time.Second}

// LinkPreview fetches a URL and extracts OpenGraph / meta tags for a rich card.
// Body: { "url": "https://leetcode.com/problems/two-sum/" }
func LinkPreview(c *gin.Context) {
	var req struct {
		URL string `json:"url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.URL) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url is required"})
		return
	}

	target := req.URL
	if !strings.HasPrefix(target, "http://") && !strings.HasPrefix(target, "https://") {
		target = "https://" + target
	}
	parsed, err := url.Parse(target)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid url"})
		return
	}

	preview := models.LinkPreview{URL: target, SiteName: parsed.Hostname()}

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	httpReq, _ := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	// Pretend to be a browser — some sites (LeetCode) block default Go UA.
	httpReq.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0 Safari/537.36")
	httpReq.Header.Set("Accept", "text/html,application/xhtml+xml")

	resp, err := previewClient.Do(httpReq)
	if err != nil {
		// Still return a minimal preview so the UI can show *something*.
		preview.Title = parsed.Hostname()
		c.JSON(http.StatusOK, preview)
		return
	}
	defer resp.Body.Close()

	// Cap the read so a huge page can't blow up memory.
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024*1024))
	parseMeta(string(body), &preview, parsed)

	if preview.Title == "" {
		preview.Title = parsed.Hostname()
	}
	if preview.Favicon == "" {
		preview.Favicon = parsed.Scheme + "://" + parsed.Host + "/favicon.ico"
	}
	c.JSON(http.StatusOK, preview)
}

// parseMeta walks the HTML tree pulling <title>, og:*, twitter:*, description, favicon.
func parseMeta(body string, p *models.LinkPreview, base *url.URL) {
	doc, err := html.Parse(strings.NewReader(body))
	if err != nil {
		return
	}

	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "title":
				if p.Title == "" && n.FirstChild != nil {
					p.Title = strings.TrimSpace(n.FirstChild.Data)
				}
			case "meta":
				key, content := "", ""
				for _, a := range n.Attr {
					switch a.Key {
					case "property", "name":
						key = a.Val
					case "content":
						content = a.Val
					}
				}
				content = strings.TrimSpace(content)
				if content == "" {
					break
				}
				switch key {
				case "og:title", "twitter:title":
					p.Title = content
				case "og:description", "twitter:description", "description":
					if p.Description == "" {
						p.Description = content
					}
				case "og:image", "twitter:image", "twitter:image:src":
					if p.Image == "" {
						p.Image = absURL(content, base)
					}
				case "og:site_name":
					p.SiteName = content
				}
			case "link":
				rel, href := "", ""
				for _, a := range n.Attr {
					switch a.Key {
					case "rel":
						rel = a.Val
					case "href":
						href = a.Val
					}
				}
				if strings.Contains(rel, "icon") && p.Favicon == "" && href != "" {
					p.Favicon = absURL(href, base)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)
}

// absURL resolves a possibly-relative asset URL against the page's base.
func absURL(ref string, base *url.URL) string {
	u, err := url.Parse(ref)
	if err != nil {
		return ref
	}
	return base.ResolveReference(u).String()
}
