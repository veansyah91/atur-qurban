package whatsapp

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_SendMessage(t *testing.T) {
	tests := []struct {
		name       string
		handler    http.HandlerFunc
		baseURL    string
		wantErr    bool
		errContain string
	}{
		{
			name: "happy path — server return 200 dengan code SUCCESS",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(SendMessageResponse{Code: "SUCCESS", Message: "ok"})
			},
			wantErr: false,
		},
		{
			name: "error HTTP non-2xx — server return 500",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("internal server error"))
			},
			wantErr:    true,
			errContain: "status 500",
		},
		{
			name:       "error baseURL kosong",
			handler:    nil,
			baseURL:    "",
			wantErr:    true,
			errContain: "tidak dikonfigurasi",
		},
		{
			name: "error response body bukan JSON valid",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("bukan json{{{"))
			},
			wantErr:    true,
			errContain: "gagal parse response",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var client *Client

			if tt.baseURL == "" && tt.handler == nil {
				// Tes baseURL kosong — tidak perlu server
				client = NewClient("", "user", "pass")
			} else {
				server := httptest.NewServer(tt.handler)
				defer server.Close()
				baseURL := tt.baseURL
				if baseURL == "" {
					baseURL = server.URL
				}
				client = NewClient(baseURL, "user", "pass")
			}

			err := client.SendMessage(context.Background(), "081234567890", "halo")
			if tt.wantErr {
				require.Error(t, err)
				if tt.errContain != "" {
					assert.Contains(t, err.Error(), tt.errContain)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}
