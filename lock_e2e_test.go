package redis_lock

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestClient_TryLock_e2e(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	testCases := []struct {
		name string
		//输入
		key        string
		expiration time.Duration

		before func()
		after  func()

		wantErr  error
		wantLock *Lock
	}{
		{
			name:       "locked",
			key:        "locked_key",
			expiration: time.Minute,
			before: func() {

			},
			after: func() {
				res, err := rdb.Del(context.Background(), "locked_key").Result()
				require.NoError(t, err)
				require.Equal(t, int64(1), res)
			},
			wantLock: &Lock{
				key: "locked_key",
			},
		},
		{
			name:       "failed lock",
			key:        "failed_key",
			expiration: time.Minute,
			before: func() {
				res, err := rdb.Set(context.Background(), "failed_key", "12345", time.Minute).Result()
				require.NoError(t, err)
				require.Equal(t, "OK", res)
			},
			after: func() {
				res, err := rdb.Get(context.Background(), "failed_key").Result()
				require.NoError(t, err)
				require.Equal(t, "12345", res)
			},
			wantErr: ErrFailedToPreemptLock,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before()
			c := NewClient(rdb)
			l, err := c.TryLock(context.Background(), tc.key, tc.expiration)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			tc.after()
			assert.NotNil(t, l.client)
			assert.Equal(t, tc.wantLock.key, l.key)
			assert.NotEmpty(t, l.value)
		})
	}
}
