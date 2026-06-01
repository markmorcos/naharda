// Package sources fetches data from upstream providers with an honest
// User-Agent and contact link (project.md §9.6).
package sources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// userAgent identifies Naharda and a contact for abuse reports (§9.6, §12).
const userAgent = "Naharda/1.0 (+https://naharda.com; abuse@naharda.com)"

var client = &http.Client{Timeout: 8 * time.Second}

// getJSON performs a GET and decodes a JSON body into target.
func getJSON(ctx context.Context, url string, target any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("upstream GET %s: status %d", url, resp.StatusCode)
	}
	return json.NewDecoder(resp.Body).Decode(target)
}
