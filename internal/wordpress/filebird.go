package wordpress

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// FbFolder reprezentuje strukturę folderu w FileBird
type FbFolder struct {
	ID       int         `json:"id"`
	Text     string      `json:"text"`
	Children []FbFolder  `json:"children"`
	Parent   interface{} `json:"parent"`
}

// FbFoldersResponse reprezentuje odpowiedź z listą folderów
type FbFoldersResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Folders []FbFolder `json:"folders"`
	} `json:"data"`
}

// FbAttachmentIDsResponse reprezentuje odpowiedź z listą ID załączników
type FbAttachmentIDsResponse struct {
	Success bool `json:"success"`
	Data    struct {
		AttachmentIDs []string `json:"attachment_ids"`
	} `json:"data"`
}

// GetFileBirdFolders pobiera listę wszystkich folderów z FileBird
func (c *Client) GetFileBirdFolders() ([]FbFolder, error) {
	if c.bearerToken == "" {
		return nil, fmt.Errorf("brak tokenu FileBird")
	}

	endpoint := fmt.Sprintf("%s/filebird/public/v1/folders", c.baseURL)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.bearerToken))

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("błąd sieci: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("serwer zwrócił kod %d", resp.StatusCode)
	}

	var res FbFoldersResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("błąd dekodowania JSON: %w", err)
	}

	return res.Data.Folders, nil
}

// GetFileBirdAttachmentIDs pobiera listę ID załączników dla danego folderu
func (c *Client) GetFileBirdAttachmentIDs(folderID int) ([]string, error) {
	if c.bearerToken == "" {
		return nil, fmt.Errorf("brak tokenu FileBird")
	}

	endpoint := fmt.Sprintf("%s/filebird/public/v1/attachment-id/?folder_id=%d", c.baseURL, folderID)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.bearerToken))

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("błąd sieci: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("serwer zwrócił kod %d", resp.StatusCode)
	}

	var res FbAttachmentIDsResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("błąd dekodowania JSON: %w", err)
	}

	return res.Data.AttachmentIDs, nil
}

// GetMediaSourceURL pobiera bezpośredni URL załącznika na podstawie jego ID
func (c *Client) GetMediaSourceURL(attachmentID string) (string, error) {
	endpoint := fmt.Sprintf("%s/wp/v2/media/%s", c.baseURL, attachmentID)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(c.username, c.password)

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("błąd sieci: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("serwer zwrócił kod %d dla ID %s", resp.StatusCode, attachmentID)
	}

	var media MediaResponse
	if err := json.NewDecoder(resp.Body).Decode(&media); err != nil {
		return "", fmt.Errorf("błąd dekodowania JSON: %w", err)
	}

	return media.SourceURL, nil
}
