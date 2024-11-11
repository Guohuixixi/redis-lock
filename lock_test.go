package redis_lock

import (
	"context"
	"errors"
	"github.com/Guohuixixi/redis-lock/mocks"
	"github.com/golang/mock/gomock"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestClient_TryLock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	testCases := []struct {
		name string
		//输入
		key        string
		expiration time.Duration

		mock func() redis.Cmdable

		wantErr  error
		wantLock *Lock
	}{
		{
			name:       "locked",
			key:        "locked_key",
			expiration: time.Minute,
			mock: func() redis.Cmdable {
				cmdable := mocks.NewMockCmdable(ctrl)
				res := redis.NewBoolResult(true, nil)
				cmdable.EXPECT().
					SetNX(gomock.Any(), "locked_key", gomock.Any(), time.Minute).
					Return(res)
				return cmdable
			},
			wantLock: &Lock{
				key: "locked_key",
			},
		},
		{
			name:       "network error",
			key:        "network_key",
			expiration: time.Minute,
			mock: func() redis.Cmdable {
				cmdable := mocks.NewMockCmdable(ctrl)
				res := redis.NewBoolResult(false, errors.New("network error"))
				cmdable.EXPECT().
					SetNX(gomock.Any(), "network_key", gomock.Any(), time.Minute).
					Return(res)
				return cmdable
			},
			wantErr: errors.New("network error"),
		},
		{
			name:       "failed lock",
			key:        "failed_key",
			expiration: time.Minute,
			mock: func() redis.Cmdable {
				cmdable := mocks.NewMockCmdable(ctrl)
				res := redis.NewBoolResult(false, nil)
				cmdable.EXPECT().
					SetNX(gomock.Any(), "failed_key", gomock.Any(), time.Minute).
					Return(res)
				return cmdable
			},
			wantErr: errors.New("rlock: 抢锁失败"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := NewClient(tc.mock())
			l, err := c.TryLock(context.Background(), tc.key, tc.expiration)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.NotNil(t, l.client)
			assert.Equal(t, tc.wantLock.key, l.key)
			assert.NotEmpty(t, l.value)
		})
	}
}
