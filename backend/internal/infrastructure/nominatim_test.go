package infrastructure

import (
	"basekarya-backend/internal/config"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNominatimFetcher_GetAddressFromCoords_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := nominatimResponse{DisplayName: "123 Main St, Springfield, IL"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	fetcher := NewNominatimFetcher(
		&config.ExternalServiceConfig{NominatimUrl: server.URL + "?lat=%f&lon=%f"},
		server.Client(),
	)

	result := fetcher.GetAddressFromCoords(39.7817, -89.6501)
	assert.Equal(t, "123 Main St, Springfield, IL", result)
}

func TestNominatimFetcher_GetAddressFromCoords_EmptyDisplayName(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := nominatimResponse{DisplayName: ""}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	fetcher := NewNominatimFetcher(
		&config.ExternalServiceConfig{NominatimUrl: server.URL + "?lat=%f&lon=%f"},
		server.Client(),
	)

	result := fetcher.GetAddressFromCoords(39.7817, -89.6501)
	assert.Contains(t, result, "39.781700")
	assert.Contains(t, result, "-89.650100")
}

func TestNominatimFetcher_GetAddressFromCoords_AllRetriesFail(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	url := server.URL + "?lat=%f&lon=%f"
	client := server.Client()
	server.Close()

	fetcher := NewNominatimFetcher(
		&config.ExternalServiceConfig{NominatimUrl: url},
		client,
	)

	result := fetcher.GetAddressFromCoords(39.7817, -89.6501)
	assert.Contains(t, result, "Unknown Location")
}
