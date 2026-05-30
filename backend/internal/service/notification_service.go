package service

import (
	"context"
	"fmt"

	"github.com/username/qurban-app/pkg/queue"
)

// NotificationService interface untuk operasi notifikasi
type NotificationService interface {
	Send(ctx context.Context, phone, message string) error
}

// notificationService implementasi NotificationService
type notificationService struct {
	notifQueue queue.NotificationQueuer
}

// NewNotificationService membuat instance baru NotificationService
func NewNotificationService(notifQueue queue.NotificationQueuer) NotificationService {
	return &notificationService{notifQueue: notifQueue}
}

// Send mengirim notifikasi WhatsApp dengan menambahkannya ke queue
func (s *notificationService) Send(ctx context.Context, phone, message string) error {
	if phone == "" {
		return fmt.Errorf("Send: nomor telepon tidak boleh kosong")
	}
	if message == "" {
		return fmt.Errorf("Send: pesan tidak boleh kosong")
	}

	job := &queue.NotificationJob{
		Phone:   phone,
		Message: message,
	}

	if err := s.notifQueue.Enqueue(ctx, job); err != nil {
		return fmt.Errorf("Send: gagal enqueue job: %w", err)
	}

	return nil
}
