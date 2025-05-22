package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	subdomain string
	email     string
	apiToken  string
	http      *http.Client
}

type OAuthClient struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Identifier  string `json:"identifier"`
	Kind        string `json:"kind"`
	Description string `json:"description,omitempty"`
}

type OAuthToken struct {
	ID        int64    `json:"id"`
	ClientID  int64    `json:"client_id"`
	UserID    int64    `json:"user_id"`
	Scopes    []string `json:"scopes"`
	FullToken string   `json:"full_token,omitempty"`
	ExpiresAt string   `json:"expires_at,omitempty"`
}

type oauthClientWrapper struct {
	Client OAuthClient `json:"client"`
}

type oauthTokenWrapper struct {
	Token OAuthToken `json:"token"`
}

func NewClient(subdomain, email, apiToken string) *Client {
	return &Client{
		subdomain: subdomain,
		email:     email,
		apiToken:  apiToken,
		http:      &http.Client{},
	}
}

func (c *Client) CreateOAuthClient(name, identifier, kind, description string) (*OAuthClient, error) {
	url := fmt.Sprintf("https://%s.zendesk.com/api/v2/oauth/clients.json", c.subdomain)
	
	payload := oauthClientWrapper{
		Client: OAuthClient{
			Name:        name,
			Identifier:  identifier,
			Kind:       kind,
			Description: description,
		},
	}
	
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(fmt.Sprintf("%s/token", c.email), c.apiToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create OAuth client: %s", string(body))
	}

	var result oauthClientWrapper
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result.Client, nil
}

func (c *Client) ReadOAuthClient(id int64) (*OAuthClient, error) {
	url := fmt.Sprintf("https://%s.zendesk.com/api/v2/oauth/clients/%d.json", c.subdomain, id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(fmt.Sprintf("%s/token", c.email), c.apiToken)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to read OAuth client: %s", string(body))
	}

	var result oauthClientWrapper
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result.Client, nil
}

func (c *Client) DeleteOAuthClient(id int64) error {
	url := fmt.Sprintf("https://%s.zendesk.com/api/v2/oauth/clients/%d.json", c.subdomain, id)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(fmt.Sprintf("%s/token", c.email), c.apiToken)

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete OAuth client: %s", string(body))
	}

	return nil
}

func (c *Client) CreateOAuthToken(clientID int64, scopes []string, expiresAt string) (*OAuthToken, error) {
	url := fmt.Sprintf("https://%s.zendesk.com/api/v2/oauth/tokens.json", c.subdomain)
	
	payload := oauthTokenWrapper{
		Token: OAuthToken{
			ClientID:  clientID,
			Scopes:    scopes,
			ExpiresAt: expiresAt,
		},
	}
	
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(fmt.Sprintf("%s/token", c.email), c.apiToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create OAuth token: %s", string(body))
	}

	var result oauthTokenWrapper
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result.Token, nil
}

func (c *Client) ReadOAuthToken(id int64) (*OAuthToken, error) {
	url := fmt.Sprintf("https://%s.zendesk.com/api/v2/oauth/tokens/%d.json", c.subdomain, id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(fmt.Sprintf("%s/token", c.email), c.apiToken)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to read OAuth token: %s", string(body))
	}

	var result oauthTokenWrapper
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result.Token, nil
}

func (c *Client) DeleteOAuthToken(id int64) error {
	url := fmt.Sprintf("https://%s.zendesk.com/api/v2/oauth/tokens/%d.json", c.subdomain, id)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(fmt.Sprintf("%s/token", c.email), c.apiToken)

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete OAuth token: %s", string(body))
	}

	return nil
} 