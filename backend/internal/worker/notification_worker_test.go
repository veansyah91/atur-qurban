package worker

import (
	"context"
	"errors"
	"io"
	"log"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/username/qurban-app/internal/mocks"
	"github.com/username/qurban-app/pkg/queue"
)

// silentLogger membuat logger yang membuang semua output agar test output bersih
func silentLogger() *log.Logger {
	return log.New(io.Discard, "", 0)
}

func TestNotificationWorker_HappyPath(t *testing.T) {
	mockQueue := new(mocks.MockNotificationQueuer)
	mockSender := new(mocks.MockWhatsAppSender)

	job := &queue.NotificationJob{Phone: "6281234567890", Message: "halo"}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Dequeue pertama return job, lalu cancel context sehingga worker berhenti
	mockQueue.On("Dequeue", mock.Anything, mock.Anything).
		Return(job, nil).Once()
	mockQueue.On("Dequeue", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) { cancel() }).
		Return((*queue.NotificationJob)(nil), nil)

	mockSender.On("SendMessage", mock.Anything, job.Phone, job.Message).
		Return(nil)

	w := NewNotificationWorker(mockQueue, mockSender, silentLogger())
	w.Start(ctx)

	// Tunggu worker berhenti secara deterministik
	<-w.Done()

	mockSender.AssertExpectations(t)
	mockQueue.AssertExpectations(t)
}

func TestNotificationWorker_QueueKosong(t *testing.T) {
	mockQueue := new(mocks.MockNotificationQueuer)
	mockSender := new(mocks.MockWhatsAppSender)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	callCount := 0
	mockQueue.On("Dequeue", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			callCount++
			if callCount >= 2 {
				cancel()
			}
		}).
		Return((*queue.NotificationJob)(nil), nil)

	w := NewNotificationWorker(mockQueue, mockSender, silentLogger())
	w.Start(ctx)

	<-w.Done()

	mockSender.AssertNotCalled(t, "SendMessage", mock.Anything, mock.Anything, mock.Anything)
}

func TestNotificationWorker_SendMessageGagal(t *testing.T) {
	mockQueue := new(mocks.MockNotificationQueuer)
	mockSender := new(mocks.MockWhatsAppSender)

	job := &queue.NotificationJob{Phone: "6281234567890", Message: "halo"}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Dequeue return job sekali, lalu cancel
	mockQueue.On("Dequeue", mock.Anything, mock.Anything).
		Return(job, nil).Once()
	mockQueue.On("Dequeue", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) { cancel() }).
		Return((*queue.NotificationJob)(nil), nil)

	// SendMessage selalu gagal (semua 3 attempt)
	mockSender.On("SendMessage", mock.Anything, job.Phone, job.Message).
		Return(errors.New("connection timeout"))

	w := NewNotificationWorker(mockQueue, mockSender, silentLogger())
	w.Start(ctx)

	<-w.Done()

	// SendMessage harus dipanggil maxRetries (3) kali sebelum menyerah
	mockSender.AssertNumberOfCalls(t, "SendMessage", maxRetries)
}

func TestNotificationWorker_ContextCancel(t *testing.T) {
	mockQueue := new(mocks.MockNotificationQueuer)
	mockSender := new(mocks.MockWhatsAppSender)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Langsung cancel sebelum Start

	mockQueue.On("Dequeue", mock.Anything, mock.Anything).
		Return((*queue.NotificationJob)(nil), context.Canceled).Maybe()

	w := NewNotificationWorker(mockQueue, mockSender, silentLogger())
	w.Start(ctx)

	// Tunggu worker berhenti secara deterministik
	<-w.Done()

	mockSender.AssertNotCalled(t, "SendMessage", mock.Anything, mock.Anything, mock.Anything)
}

