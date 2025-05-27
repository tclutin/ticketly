package ticketly

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client interface {
	CreateTicket(request CreateTicketRequest) (uint64, error)
}

type Ticketly struct {
	client *http.Client
	apiURL string
}

func NewTicketly() *Ticketly {
	return &Ticketly{
		client: &http.Client{},
		apiURL: "http://localhost:8090/api/v1",
	}
}

func (t *Ticketly) Register(request RegisterUserRequest) (uint64, error) {
	if err := request.Validate(); err != nil {
		return 0, err
	}

	data, err := json.Marshal(request)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal request: %w", err)
	}

	endpoint := fmt.Sprintf("%s/users", t.apiURL)

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(data))
	if err != nil {
		return 0, fmt.Errorf("failed to send request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := t.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Println(string(bodyBytes))
		return 0, fmt.Errorf("failed to register user: invalid status code: %d", resp.StatusCode)
	}

	var response struct {
		UserId uint64 `json:"user_id"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, fmt.Errorf("failed to decode response body: %w", err)
	}

	return response.UserId, nil
}

func (t *Ticketly) CreateTicket(request CreateTicketRequest) (uint64, error) {
	if err := request.Validate(); err != nil {
		return 0, err
	}

	data, err := json.Marshal(request)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal request: %w", err)
	}

	endpoint := fmt.Sprintf("%s/tickets", t.apiURL)

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(data))
	if err != nil {
		return 0, fmt.Errorf("failed to send request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := t.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return 0, fmt.Errorf("failed to create ticket: invalid status code: %d", resp.StatusCode)
	}

	var response struct {
		TicketId uint64 `json:"ticket_id"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, fmt.Errorf("failed to decode response body: %w", err)
	}

	return response.TicketId, nil
}
