package queue

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupMiniredis(t *testing.T) (*miniredis.Miniredis, *redis.Client) {
	mr, err := miniredis.Run()
	require.NoError(t, err)
	t.Cleanup(mr.Close)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	return mr, rdb
}

func TestNotificationQueue_Enqueue(t *testing.T) {
	t.Run("happy path — job masuk ke queue", func(t *testing.T) {
		_, rdb := setupMiniredis(t)
		q := NewNotificationQueue(rdb)

		job := &NotificationJob{Phone: "6281234567890", Message: "halo"}
		err := q.Enqueue(context.Background(), job)
		require.NoError(t, err)

		length, err := q.Length(context.Background())
		require.NoError(t, err)
		assert.Equal(t, int64(1), length)
	})

	t.Run("error job nil", func(t *testing.T) {
		_, rdb := setupMiniredis(t)
		q := NewNotificationQueue(rdb)

		err := q.Enqueue(context.Background(), nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "job tidak boleh nil")
	})
}

func TestNotificationQueue_Dequeue(t *testing.T) {
	t.Run("happy path — ambil job dari queue", func(t *testing.T) {
		_, rdb := setupMiniredis(t)
		q := NewNotificationQueue(rdb)

		job := &NotificationJob{Phone: "6281234567890", Message: "halo"}
		err := q.Enqueue(context.Background(), job)
		require.NoError(t, err)

		got, err := q.Dequeue(context.Background(), 1*time.Second)
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, job.Phone, got.Phone)
		assert.Equal(t, job.Message, got.Message)
	})

	t.Run("queue kosong — return nil saat timeout", func(t *testing.T) {
		_, rdb := setupMiniredis(t)
		q := NewNotificationQueue(rdb)

		// Gunakan context dengan timeout singkat agar test tidak lambat
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		got, err := q.Dequeue(ctx, 1*time.Second)
		// Bisa return nil job (timeout) atau context error
		if err != nil {
			assert.True(t, ctx.Err() != nil || err == redis.Nil)
		} else {
			assert.Nil(t, got)
		}
	})
}

func TestNotificationQueue_Length(t *testing.T) {
	t.Run("happy path — panjang queue sesuai jumlah job", func(t *testing.T) {
		_, rdb := setupMiniredis(t)
		q := NewNotificationQueue(rdb)

		length, err := q.Length(context.Background())
		require.NoError(t, err)
		assert.Equal(t, int64(0), length)

		require.NoError(t, q.Enqueue(context.Background(), &NotificationJob{Phone: "62811", Message: "a"}))
		require.NoError(t, q.Enqueue(context.Background(), &NotificationJob{Phone: "62812", Message: "b"}))

		length, err = q.Length(context.Background())
		require.NoError(t, err)
		assert.Equal(t, int64(2), length)
	})
}
