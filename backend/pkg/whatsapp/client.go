package whatsapp

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client mengelola komunikasi dengan GOWA WhatsApp Gateway
type Client struct {
	baseURL  string
	username string
	password string
	client   *http.Client
}

// SendMessageRequest adalah struktur request untuk mengirim pesan
type SendMessageRequest struct {
	Phone   string `json:"phone"`
	Message string `json:"message"`
}

// SendMessageResponse adalah struktur response dari GOWA
type SendMessageResponse struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Results interface{} `json:"results"`
}

// NewClient membuat instance baru GOWA client
// WhatsAppSender interface untuk mengirim pesan WhatsApp
type WhatsAppSender interface {
	SendMessage(ctx context.Context, phone, message string) error
}

func NewClient(baseURL, username, password string) *Client {
	return &Client{
		baseURL:  baseURL,
		username: username,
		password: password,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SendMessage mengirim pesan WhatsApp melalui GOWA API
func (c *Client) SendMessage(ctx context.Context, phone, message string) error {
	if c.baseURL == "" {
		return fmt.Errorf("GOWA URL tidak dikonfigurasi")
	}

	// Format phone ke format GOWA: 628xxx@s.whatsapp.net
	formattedPhone := FormatPhoneNumber(phone)

	reqBody := SendMessageRequest{
		Phone:   formattedPhone,
		Message: message,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("SendMessage: gagal marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("SendMessage: gagal membuat request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.getBasicAuth())

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("SendMessage: gagal mengirim request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("SendMessage: gagal membaca response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("SendMessage: GOWA return status %d: %s", resp.StatusCode, string(body))
	}

	var respData SendMessageResponse
	if err := json.Unmarshal(body, &respData); err != nil {
		return fmt.Errorf("SendMessage: gagal parse response: %w", err)
	}

	if respData.Code != "SUCCESS" && respData.Code != "200" && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("SendMessage: GOWA error - %s: %s", respData.Code, respData.Message)
	}

	return nil
}

// getBasicAuth mengembalikan header Basic Auth yang sudah di-encode
func (c *Client) getBasicAuth() string {
	credentials := c.username + ":" + c.password
	encoded := base64.StdEncoding.EncodeToString([]byte(credentials))
	return "Basic " + encoded
}
