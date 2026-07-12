package webhook

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
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
	url        string
	httpClient *http.Client
}

func NewClient(url string) *Client {
	return &Client{
		url:        url,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// Send posts the event payload, retrying up to maxAttempts times (with a
// fixed delay between attempts) until a 200 response is received. A no-op
// if no webhook URL is configured. Runs synchronously — callers that don't
// want to block should call this in its own goroutine.
func (c *Client) Send(event string, codes []string) {
	if c.url == "" {
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

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		req, err := http.NewRequest(http.MethodPost, c.url, bytes.NewReader(body))
		if err != nil {
			log.Printf("webhook: failed to build request: %v", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			log.Printf("webhook: attempt %d/%d for event %q failed: %v", attempt, maxAttempts, event, err)
		} else {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return
			}
			log.Printf("webhook: attempt %d/%d for event %q got status %d", attempt, maxAttempts, event, resp.StatusCode)
		}

		if attempt < maxAttempts {
			time.Sleep(retryDelay)
		}
	}

	log.Printf("webhook: giving up after %d attempts for event %q (codes=%v)", maxAttempts, event, codes)
}
