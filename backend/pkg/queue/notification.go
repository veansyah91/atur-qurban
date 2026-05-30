package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const NotificationQueueKey = "queue:notification"

// NotificationQueuer interface untuk operasi queue notifikasi
type NotificationQueuer interface {
	Enqueue(ctx context.Context, job *NotificationJob) error
	Dequeue(ctx context.Context, timeout time.Duration) (*NotificationJob, error)
	Length(ctx context.Context) (int64, error)
}

// NotificationJob merepresentasikan job notifikasi WhatsApp
type NotificationJob struct {
	Phone   string `json:"phone"`
	Message string `json:"message"`
}

// NotificationQueue mengelola queue notifikasi berbasis Redis
type NotificationQueue struct {
	rdb *redis.Client
}

// NewNotificationQueue membuat instance baru NotificationQueue
func NewNotificationQueue(rdb *redis.Client) *NotificationQueue {
	return &NotificationQueue{rdb: rdb}
}

// Enqueue menambahkan job ke queue
func (q *NotificationQueue) Enqueue(ctx context.Context, job *NotificationJob) error {
	if job == nil {
		return fmt.Errorf("Enqueue: job tidak boleh nil")
	}

	jsonData, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("Enqueue: gagal marshal job: %w", err)
	}

	// Gunakan context timeout jika ada
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
	}

	result := q.rdb.LPush(ctx, NotificationQueueKey, jsonData)
	if err := result.Err(); err != nil {
		return fmt.Errorf("Enqueue: gagal push ke queue: %w", err)
	}

	return nil
}

// Dequeue mengambil job dari queue (blocking operation dengan timeout)
func (q *NotificationQueue) Dequeue(ctx context.Context, timeout time.Duration) (*NotificationJob, error) {
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	result := q.rdb.BRPop(ctx, timeout, NotificationQueueKey)
	if err := result.Err(); err != nil {
		if err == redis.Nil {
			return nil, nil // Queue kosong, timeout
		}
		return nil, fmt.Errorf("Dequeue: gagal pop dari queue: %w", err)
	}

	values := result.Val()
	if len(values) < 2 {
		return nil, fmt.Errorf("Dequeue: response format tidak valid")
	}

	var job NotificationJob
	if err := json.Unmarshal([]byte(values[1]), &job); err != nil {
		return nil, fmt.Errorf("Dequeue: gagal unmarshal job: %w", err)
	}

	return &job, nil
}

// Length mengembalikan jumlah job dalam queue
func (q *NotificationQueue) Length(ctx context.Context) (int64, error) {
	result := q.rdb.LLen(ctx, NotificationQueueKey)
	if err := result.Err(); err != nil {
		return 0, fmt.Errorf("Length: gagal mendapatkan panjang queue: %w", err)
	}
	return result.Val(), nil
}
