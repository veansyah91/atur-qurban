package mocks

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/username/qurban-app/pkg/queue"
)

// MockNotificationQueuer adalah mock untuk interface NotificationQueuer
type MockNotificationQueuer struct {
	mock.Mock
}

func (m *MockNotificationQueuer) Enqueue(ctx context.Context, job *queue.NotificationJob) error {
	args := m.Called(ctx, job)
	return args.Error(0)
}

func (m *MockNotificationQueuer) Dequeue(ctx context.Context, timeout time.Duration) (*queue.NotificationJob, error) {
	args := m.Called(ctx, timeout)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*queue.NotificationJob), args.Error(1)
}

func (m *MockNotificationQueuer) Length(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}
