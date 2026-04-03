package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientListCurrencies(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "adds api key header and query params",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if got := r.Header.Get("X-API-Key"); got != "test-key" {
					t.Fatalf("X-API-Key = %q, want %q", got, "test-key")
				}

				if got := r.URL.Query().Get("limit"); got != "5" {
					t.Fatalf("limit = %q, want %q", got, "5")
				}

				if got := r.URL.Query().Get("offset"); got != "10" {
					t.Fatalf("offset = %q, want %q", got, "10")
				}

				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`[{"code":"USD","name":"US Dollar","symbol":"$","decimalPlaces":2}]`))
			}))
			defer server.Close()

			client := NewClient(Config{
				BaseURL: server.URL,
				APIKey:  "test-key",
			})

			currencies, err := client.ListCurrencies(context.Background(), ListCurrenciesParams{
				Limit:  5,
				Offset: 10,
			})
			if err != nil {
				t.Fatalf("ListCurrencies() error = %v", err)
			}

			if len(currencies) != 1 {
				t.Fatalf("len(currencies) = %d, want %d", len(currencies), 1)
			}

			if got := currencies[0].Code; got != "USD" {
				t.Fatalf("currencies[0].Code = %q, want %q", got, "USD")
			}
		})
	}
}

func TestClientListCurrenciesUnauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("unauthorized"))
	}))
	defer server.Close()

	client := NewClient(Config{
		BaseURL: server.URL,
		APIKey:  "bad-key",
	})

	_, err := client.ListCurrencies(context.Background(), ListCurrenciesParams{})
	if err == nil {
		t.Fatal("ListCurrencies() error = nil, want error")
	}
}
