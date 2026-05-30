package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockWhatsAppSender adalah mock untuk interface WhatsAppSender
type MockWhatsAppSender struct {
	mock.Mock
}

func (m *MockWhatsAppSender) SendMessage(ctx context.Context, phone, message string) error {
	args := m.Called(ctx, phone, message)
	return args.Error(0)
}
