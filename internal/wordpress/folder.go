package wordpress

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// FbCreateFolderRequest - body do tworzenia folderu
type FbCreateFolderRequest struct {
	Name     string `json:"name"`
	ParentID int    `json:"parent_id"` // API: "parent_id"
}

// FbCreateFolderResponse - odpowiedź z ID nowego folderu
type FbCreateFolderResponse struct {
	Success bool `json:"success"`
	Data    struct {
		ID int `json:"id"`
	} `json:"data"`
}

// FbSetAttachmentRequest - body do przypisania plików
type FbSetAttachmentRequest struct {
	FolderID int   `json:"folder"` // API: "folder" (number)
	Ids      []int `json:"ids"`    // API: "ids" (number or array)
}

// FbBasicResponse - generyczna odpowiedź (np. dla set-attachment)
type FbBasicResponse struct {
	Success bool `json:"success"`
}

// CreateFolder tworzy nowy folder w FileBird i zwraca jego ID.
func (c *Client) CreateFolder(name string, parentID int) (int, error) {

	if c.bearerToken == "" {
		return 0, fmt.Errorf("brak tokenu FileBird - operacja niemożliwa")
	}

	endpoint := fmt.Sprintf("%s/filebird/public/v1/folders", c.baseURL)

	payload := FbCreateFolderRequest{
		Name:     name,
		ParentID: parentID,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.bearerToken))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return 0, fmt.Errorf("błąd sieci FileBird: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("FileBird CreateFolder błąd %d: %s", resp.StatusCode, string(body))
	}

	var res FbCreateFolderResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return 0, fmt.Errorf("błąd dekodowania odpowiedzi JSON: %w", err)
	}

	if !res.Success {
		return 0, fmt.Errorf("FileBird zwrócił success: false przy tworzeniu folderu")
	}

	return res.Data.ID, nil
}

// SetAttachmentFolder przypisuje pliki (IDs) do konkretnego folderu.
func (c *Client) SetAttachmentFolder(folderID int, attachmentIDs []int) error {

	if c.bearerToken == "" {
		return nil
	}

	endpoint := fmt.Sprintf("%s/filebird/public/v1/folder/set-attachment", c.baseURL)

	payload := FbSetAttachmentRequest{
		FolderID: folderID,
		Ids:      attachmentIDs,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.bearerToken))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("błąd sieci FileBird: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("FileBird SetAttachment błąd %d: %s", resp.StatusCode, string(body))
	}

	var res FbBasicResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return fmt.Errorf("błąd dekodowania odpowiedzi JSON: %w", err)
	}

	if !res.Success {
		return fmt.Errorf("FileBird nie udało się przypisać plików do folderu")
	}

	return nil
}
