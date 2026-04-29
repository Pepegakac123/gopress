package wordpress

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRequest_Headers(t *testing.T) {
	client := NewClient("https://example.com", "user", "pass", "token")
	
	req, err := client.newRequest("GET", "https://example.com/test", nil)
	if err != nil {
		t.Fatalf("newRequest failed: %v", err)
	}

	ua := req.Header.Get("User-Agent")
	if ua != "GoPress/1.0 (WordPress Image Optimizer)" {
		t.Errorf("Expected User-Agent 'GoPress/1.0 (WordPress Image Optimizer)', got '%s'", ua)
	}

	accept := req.Header.Get("Accept")
	if accept != "application/json, */*" {
		t.Errorf("Expected Accept 'application/json, */*', got '%s'", accept)
	}
}

func TestCheckConnection_HeadersSent(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ua := r.Header.Get("User-Agent")
		if ua != "GoPress/1.0 (WordPress Image Optimizer)" {
			t.Errorf("Server received wrong User-Agent: %s", ua)
		}
		
		accept := r.Header.Get("Accept")
		if accept != "application/json, */*" {
			t.Errorf("Server received wrong Accept: %s", accept)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := NewClient(ts.URL, "user", "pass", "token")
	err := client.CheckConnection()
	if err != nil {
		t.Errorf("CheckConnection failed: %v", err)
	}
}
