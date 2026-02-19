package infrastructure

import (
	"net/http"
	"time"
)

type HttpClientProvider struct {
	client *http.Client
}

func NewHttpClientProvider() *HttpClientProvider {
	t := &http.Transport{
		MaxIdleConns:      10,
		IdleConnTimeout:   30 * time.Second,
		DisableKeepAlives: false,
	}

	client := &http.Client{
		Transport: t,
		Timeout:   15 * time.Second,
	}

	return &HttpClientProvider{client}
}

func (h *HttpClientProvider) GetClient() *http.Client {
	return h.client
}
