package wordpress

import (
	"fmt"
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
	}
}

// CheckConnection sprawdza, czy dane logowania są poprawne.
// Uderza w endpoint /users/me, który wymaga autoryzacji.
func (c *Client) CheckConnection() error {
	endpoint := fmt.Sprintf("%s/wp/v2/users/me", c.baseURL)

	req, err := http.NewRequest("GET", endpoint, nil)
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
