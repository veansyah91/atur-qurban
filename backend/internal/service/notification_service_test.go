package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/username/qurban-app/internal/mocks"
	"github.com/username/qurban-app/pkg/queue"
)

func TestNotificationService_Send(t *testing.T) {
	tests := []struct {
		name       string
		phone      string
		message    string
		mockSetup  func(*mocks.MockNotificationQueuer)
		wantErr    bool
		errContain string
	}{
		{
			name:    "happy path — enqueue dipanggil dengan data benar",
			phone:   "+6281234567890",
			message: "halo peserta",
			mockSetup: func(m *mocks.MockNotificationQueuer) {
				m.On("Enqueue", mock.Anything, &queue.NotificationJob{
					Phone:   "+6281234567890",
					Message: "halo peserta",
				}).Return(nil)
			},
			wantErr: false,
		},
		{
			name:    "error phone kosong — enqueue tidak dipanggil",
			phone:   "",
			message: "halo",
			mockSetup: func(m *mocks.MockNotificationQueuer) {
				// Tidak ada expectation — Enqueue tidak boleh dipanggil
			},
			wantErr:    true,
			errContain: "nomor telepon tidak boleh kosong",
		},
		{
			name:    "error message kosong — enqueue tidak dipanggil",
			phone:   "+6281234567890",
			message: "",
			mockSetup: func(m *mocks.MockNotificationQueuer) {
				// Tidak ada expectation — Enqueue tidak boleh dipanggil
			},
			wantErr:    true,
			errContain: "pesan tidak boleh kosong",
		},
		{
			name:    "error enqueue gagal — send return wrapped error",
			phone:   "+6281234567890",
			message: "halo",
			mockSetup: func(m *mocks.MockNotificationQueuer) {
				m.On("Enqueue", mock.Anything, mock.Anything).
					Return(errors.New("redis connection refused"))
			},
			wantErr:    true,
			errContain: "gagal enqueue job",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQueue := new(mocks.MockNotificationQueuer)
			tt.mockSetup(mockQueue)

			svc := NewNotificationService(mockQueue)
			err := svc.Send(context.Background(), tt.phone, tt.message)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContain != "" {
					assert.Contains(t, err.Error(), tt.errContain)
				}
			} else {
				require.NoError(t, err)
			}

			mockQueue.AssertExpectations(t)
		})
	}
}
