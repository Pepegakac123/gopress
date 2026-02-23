package wordpress

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_GetFileBirdFolders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		resp := FbFoldersResponse{
			Success: true,
		}
		resp.Data.Folders = []FbFolder{
			{ID: 1, Text: "Folder 1", Children: []FbFolder{{ID: 2, Text: "Subfolder 1"}}},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "user", "pass", "test-token")
	folders, err := client.GetFileBirdFolders()
	if err != nil {
		t.Fatalf("Błąd: %v", err)
	}

	if len(folders) != 1 {
		t.Errorf("Oczekiwano 1 folderu, otrzymano %d", len(folders))
	}

	if folders[0].Text != "Folder 1" {
		t.Errorf("Oczekiwano 'Folder 1', otrzymano '%s'", folders[0].Text)
	}

	if len(folders[0].Children) != 1 || folders[0].Children[0].Text != "Subfolder 1" {
		t.Errorf("Błąd w strukturze dzieci")
	}
}

func TestClient_GetFileBirdAttachmentIDs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := FbAttachmentIDsResponse{
			Success: true,
		}
		resp.Data.AttachmentIDs = []string{"10", "20"}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "user", "pass", "token")
	ids, err := client.GetFileBirdAttachmentIDs(5)
	if err != nil {
		t.Fatalf("Błąd: %v", err)
	}

	if len(ids) != 2 || ids[0] != "10" || ids[1] != "20" {
		t.Errorf("Nieprawidłowe ID załączników: %v", ids)
	}
}

func TestClient_GetMediaSourceURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := MediaResponse{
			ID:        10,
			SourceURL: "https://example.com/img.jpg",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "user", "pass", "token")
	url, err := client.GetMediaSourceURL("10")
	if err != nil {
		t.Fatalf("Błąd: %v", err)
	}

	if url != "https://example.com/img.jpg" {
		t.Errorf("Oczekiwano https://example.com/img.jpg, otrzymano %s", url)
	}
}

func TestFindFolderRecursive(t *testing.T) {
	folders := []FbFolder{
		{ID: 1, Text: "Materiały", Children: []FbFolder{}},
		{ID: 2, Text: "Zdjęcia", Children: []FbFolder{
			{ID: 3, Text: "Ogrzewanie", Children: []FbFolder{}},
			{ID: 4, Text: "5 zdjęć do każdej usługi", Children: []FbFolder{
				{ID: 5, Text: "Wylewki", Children: []FbFolder{}},
			}},
		}},
	}

	found := findFolderRecursive(folders, "5 zdjęć do każdej usługi", "5 zdjęć")
	if found == nil || found.ID != 4 {
		t.Errorf("Nie znaleziono odpowiedniego folderu")
	}

	found = findFolderRecursive(folders, "nie-istnieje", "")
	if found != nil {
		t.Errorf("Znaleziono folder, który nie powinien istnieć")
	}

	found = findFolderRecursive(folders, "5 zdjęć", "")
	if found == nil || found.ID != 4 {
		t.Errorf("Nie znaleziono odpowiedniego folderu po skróconej nazwie")
	}
}
