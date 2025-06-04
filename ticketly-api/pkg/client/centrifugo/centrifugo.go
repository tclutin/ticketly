package centrifugo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
	apiURL     string
	apiKey     string
}

type Message struct {
	Data any `json:"data"`
}

func New(apiURL, apikey, secret string) *Client {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	return &Client{
		httpClient: client,
		apiURL:     apiURL,
		apiKey:     apikey,
	}
}

func (c *Client) Publish(channel, data any) error {
	request := map[string]interface{}{
		"channel": channel,
		"data":    data,
	}

	body, _ := json.Marshal(request)

	req, _ := http.NewRequest("POST", c.apiURL+"/api/publish", bytes.NewBuffer(body))

	req.Header.Set("X-API-Key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result struct {
		Error any `json:"error,omitempty"`
	}

	var x any
	json.NewDecoder(resp.Body).Decode(&x)

	fmt.Println(x)
	if result.Error != nil {
		return fmt.Errorf("centrifugo error: %s", result.Error)
	}

	return nil
}
