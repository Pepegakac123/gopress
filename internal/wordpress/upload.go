package wordpress

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type MediaResponse struct {
	ID        int    `json:"id"`
	SourceURL string `json:"source_url"`
	Link      string `json:"link"`
}

func (c *Client) UploadFile(filePath string) (*MediaResponse, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("Nie można otworzyć pliki: %w", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, fmt.Errorf("Nie można utworzyć formularza: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("Nie można przesłać pliku: %w", err)
	}
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("Nie można zamknąć formularza: %w", err)
	}

	endpoint := fmt.Sprintf("%s/wp/v2/media", c.baseURL)
	req, err := http.NewRequest("POST", endpoint, body)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filepath.Base(filePath)))
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("błąd sieci: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("WP zwrócił błąd %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var media MediaResponse
	if err := json.NewDecoder(resp.Body).Decode(&media); err != nil {
		return nil, fmt.Errorf("błąd dekodowania JSON: %w", err)
	}

	return &media, nil
}
