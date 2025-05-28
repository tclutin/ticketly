package ticketly

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Client interface {
	CreateTicket(request CreateTicketRequest) (uint64, error)
	Register(request RegisterUserRequest) (uint64, error)
	GetUserByExternalId(externalId string) (UserResponse, error)
}

type Ticketly struct {
	client *http.Client
	apiURL string
}

func NewClient() *Ticketly {
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
		return 0, fmt.Errorf("invalid status code: %d", resp.StatusCode)
	}

	var response struct {
		UserId uint64 `json:"user_id"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, fmt.Errorf("failed to decode response body: %w", err)
	}

	return response.UserId, nil
}

func (t *Ticketly) GetUserByExternalId(externalId string) (UserResponse, error) {
	endpoint := fmt.Sprintf("%s/users/%s", t.apiURL, url.PathEscape(externalId))

	var response UserResponse

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return response, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := t.client.Do(req)
	if err != nil {
		return response, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return response, fmt.Errorf("invalid status code: %d", resp.StatusCode)
	}

	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return response, fmt.Errorf("failed to decode response body: %w", err)
	}

	return response, nil
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
