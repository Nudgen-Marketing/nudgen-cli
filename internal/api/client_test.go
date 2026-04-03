package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetCampaigns(t *testing.T) {
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		if req.URL.Path != "/campaigns" {
			t.Errorf("Expected path /campaigns, got %s", req.URL.Path)
		}
		if req.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("Expected Authorization test-token, got %s", req.Header.Get("Authorization"))
		}

		// Send response wrapped in campaigns field to match real backend structure
		rw.WriteHeader(http.StatusOK)
		json.NewEncoder(rw).Encode(struct {
			Campaigns []Campaign `json:"campaigns"`
		}{
			Campaigns: []Campaign{
				{ID: "1", Name: "Campaign 1", Status: "active"},
			},
		})
	}))
	// Close the server when test finishes
	defer server.Close()

	// Use Client & confirm it works
	client := NewClientWithToken("test-token")
	client.BaseURL = server.URL // Override for testing

	campaigns, err := client.GetCampaigns()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(campaigns) != 1 {
		t.Errorf("Expected 1 campaign, got %d", len(campaigns))
	}

	if campaigns[0].Name != "Campaign 1" {
		t.Errorf("Expected Campaign 1, got %s", campaigns[0].Name)
	}
}

func TestGetContacts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
		json.NewEncoder(rw).Encode(struct {
			Contacts []Contact `json:"contacts"`
		}{
			Contacts: []Contact{
				{ID: "c1", Name: "User 1", Email: "user1@example.com"},
			},
		})
	}))
	defer server.Close()

	client := NewClientWithToken("test-token")
	client.BaseURL = server.URL

	contacts, err := client.GetContacts()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(contacts) != 1 {
		t.Errorf("Expected 1 contact, got %d", len(contacts))
	}
}
