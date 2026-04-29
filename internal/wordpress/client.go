package wordpress

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client to instancja połączenia z WP.
type Client struct {
	baseURL     string
	username    string
	password    string
	bearerToken string
	http        *http.Client
	linkCache   map[string]string // Cache: attachmentID -> SourceURL
}

// NewClient to konstruktor
func NewClient(domain, user, password, bearerToken string) *Client {
	domain = strings.TrimSuffix(domain, "/")
	domain = strings.TrimSuffix(domain, "/wp-json")
	apiURL := fmt.Sprintf("%s/wp-json", domain)
	return &Client{
		baseURL:     apiURL,
		username:    user,
		password:    password,
		bearerToken: bearerToken,
		http: &http.Client{
			Timeout: time.Second * 30,
		},
		linkCache: make(map[string]string),
	}
}

// newRequest tworzy nowe zapytanie HTTP z ustawionymi nagłówkami bezpieczeństwa.
func (c *Client) newRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	// Ustawiamy User-Agent, aby uniknąć błędów 406 i blokad firewalli.
	req.Header.Set("User-Agent", "GoPress/1.0 (WordPress Image Optimizer)")
	req.Header.Set("Accept", "application/json, */*")

	return req, nil
}

// CheckConnection sprawdza, czy dane logowania są poprawne.
// Uderza w endpoint /users/me, który wymaga autoryzacji.
func (c *Client) CheckConnection() error {
	endpoint := fmt.Sprintf("%s/wp/v2/users/me", c.baseURL)

	req, err := c.newRequest("GET", endpoint, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.username, c.password)
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("Bład sieci: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("błąd autoryzacji: serwer zwrócił %d", resp.StatusCode)
	}
	return nil
}

// CheckFileBirdConnection sprawdza, czy token Bearer jest poprawny.
func (c *Client) CheckFileBirdConnection() error {
	if c.bearerToken == "" {
		return fmt.Errorf("brak tokenu")
	}
	endpoint := fmt.Sprintf("%s/filebird/public/v1/folders", c.baseURL)

	req, err := c.newRequest("GET", endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.bearerToken))

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("błąd sieci: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("nieprawidłowy token (kod %d)", resp.StatusCode)
	}

	return nil
}
