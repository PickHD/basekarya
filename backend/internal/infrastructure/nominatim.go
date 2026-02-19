package infrastructure

import (
	"encoding/json"
	"fmt"
	"hris-backend/internal/config"
	"hris-backend/pkg/logger"
	"net/http"
	"time"
)

type nominatimResponse struct {
	DisplayName string `json:"display_name"`
}

type NominatimFetcher struct {
	client *http.Client
	url    string
}

func NewNominatimFetcher(cfg *config.ExternalServiceConfig) *NominatimFetcher {
	t := &http.Transport{
		MaxIdleConns:      10,
		IdleConnTimeout:   30 * time.Second,
		DisableKeepAlives: false,
	}

	client := &http.Client{
		Transport: t,
		Timeout:   15 * time.Second,
	}

	return &NominatimFetcher{
		client: client,
		url:    cfg.NominatimUrl,
	}
}

func (n *NominatimFetcher) GetAddressFromCoords(lat, long float64) string {
	url := fmt.Sprintf(n.url, lat, long)

	maxRetries := 3
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		if i > 0 {
			time.Sleep(time.Duration(i*2) * time.Second)
			logger.Infof("Retrying nominatim fetch... attempt %d", i+1)
		}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return fmt.Sprintf("%f, %f", lat, long)
		}

		req.Header.Set("User-Agent", "HRIS-App-Backend/1.0 (taufik@januar35@gmail.com)")

		resp, err := n.client.Do(req)
		if err != nil {
			lastErr = err
			continue
		}

		var result nominatimResponse
		decodeErr := json.NewDecoder(resp.Body).Decode(&result)
		resp.Body.Close()

		if decodeErr != nil {
			logger.Errorw("failed decode JSON", decodeErr)
			return fmt.Sprintf("%f, %f", lat, long)
		}

		if result.DisplayName != "" {
			return result.DisplayName
		}
	}

	logger.Warnf("failed to fetch nominatim after retries: %v", lastErr)
	return fmt.Sprintf("Unknown Location (%f, %f)", lat, long)
}
