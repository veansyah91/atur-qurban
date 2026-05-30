package worker

import (
	"context"
	"log"
	"time"

	"github.com/username/qurban-app/pkg/queue"
	"github.com/username/qurban-app/pkg/whatsapp"
)

const (
	maxRetries    = 3
	retryBaseWait = 2 * time.Second
)

// NotificationWorker memproses job notifikasi dari queue
type NotificationWorker struct {
	notifQueue queue.NotificationQueuer
	waClient   whatsapp.WhatsAppSender
	logger     *log.Logger
	done       chan struct{}
}

// NewNotificationWorker membuat instance baru NotificationWorker
func NewNotificationWorker(notifQueue queue.NotificationQueuer, waClient whatsapp.WhatsAppSender, logger *log.Logger) *NotificationWorker {
	if logger == nil {
		logger = log.New(log.Writer(), "[NotificationWorker] ", log.LstdFlags)
	}
	return &NotificationWorker{
		notifQueue: notifQueue,
		waClient:   waClient,
		logger:     logger,
		done:       make(chan struct{}),
	}
}

// Start menjalankan worker untuk memproses job secara terus-menerus
func (w *NotificationWorker) Start(ctx context.Context) {
	go w.run(ctx)
	w.logger.Println("Worker dimulai dan siap memproses notifikasi")
}

// Done mengembalikan channel yang ditutup saat worker selesai berjalan.
// Gunakan ini di test untuk menunggu worker berhenti secara deterministik.
func (w *NotificationWorker) Done() <-chan struct{} {
	return w.done
}

// run adalah loop utama worker yang memproses queue
func (w *NotificationWorker) run(ctx context.Context) {
	defer close(w.done)

	for {
		select {
		case <-ctx.Done():
			w.logger.Println("Worker dihentikan")
			return
		default:
		}

		// Dequeue dengan timeout 5 detik
		job, err := w.notifQueue.Dequeue(ctx, 5*time.Second)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			w.logger.Printf("Error dequeue: %v\n", err)
			continue
		}

		if job == nil {
			continue
		}

		w.sendWithRetry(ctx, job)
	}
}

// sendWithRetry mencoba mengirim pesan dengan retry maksimal maxRetries kali
func (w *NotificationWorker) sendWithRetry(ctx context.Context, job *queue.NotificationJob) {
	var lastErr error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		sendCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		lastErr = w.waClient.SendMessage(sendCtx, job.Phone, job.Message)
		cancel()

		if lastErr == nil {
			w.logger.Printf("Successfully sent message to %s (attempt %d)\n", job.Phone, attempt)
			return
		}

		w.logger.Printf("Failed to send message to %s (attempt %d/%d): %v\n", job.Phone, attempt, maxRetries, lastErr)

		if attempt < maxRetries {
			wait := retryBaseWait * time.Duration(attempt)
			select {
			case <-ctx.Done():
				return
			case <-time.After(wait):
			}
		}
	}
	w.logger.Printf("Giving up sending message to %s after %d attempts: %v\n", job.Phone, maxRetries, lastErr)
}

// GetQueueLength mengembalikan jumlah job yang masih dalam queue
func (w *NotificationWorker) GetQueueLength(ctx context.Context) (int64, error) {
	return w.notifQueue.Length(ctx)
}

