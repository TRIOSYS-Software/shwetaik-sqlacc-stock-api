package webhook

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	maxAttempts = 3
	retryDelay  = 2 * time.Second
)

type Payload struct {
	Event     string   `json:"event"`
	Timestamp int64    `json:"timestamp"`
	Codes     []string `json:"codes"`
	Count     int      `json:"count"`
}

type Client struct {
	urls       []string
	httpClient *http.Client
}

func NewClient(urls []string) *Client {
	return &Client{
		urls:       urls,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// Send posts the event payload to every configured URL concurrently, each
// retrying independently up to maxAttempts times (with a fixed delay
// between attempts) until a 200 response is received. A no-op if no
// webhook URLs are configured. Blocks until every URL has either succeeded
// or exhausted its retries — callers that don't want to block should call
// this in its own goroutine.
func (c *Client) Send(event string, codes []string) {
	if len(c.urls) == 0 {
		return
	}

	payload := Payload{
		Event:     event,
		Timestamp: time.Now().Unix(),
		Codes:     codes,
		Count:     len(codes),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("webhook: failed to marshal payload for event %q: %v", event, err)
		return
	}

	var wg sync.WaitGroup
	for _, url := range c.urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			c.sendTo(url, event, codes, body)
		}(url)
	}
	wg.Wait()
}

func (c *Client) sendTo(url, event string, codes []string, body []byte) {
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
		if err != nil {
			log.Printf("webhook: failed to build request for %s: %v", url, err)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			log.Printf("webhook: attempt %d/%d to %s for event %q failed: %v", attempt, maxAttempts, url, event, err)
		} else {
			// Drain the body before closing so the connection can be reused.
			_, _ = io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return
			}
			log.Printf("webhook: attempt %d/%d to %s for event %q got status %d", attempt, maxAttempts, url, event, resp.StatusCode)
		}

		if attempt < maxAttempts {
			time.Sleep(retryDelay)
		}
	}

	log.Printf("webhook: giving up after %d attempts to %s for event %q (codes=%v)", maxAttempts, url, event, codes)
}
