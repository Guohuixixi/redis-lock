package redis_lock

import "time"

type RetryStrategy interface {
	Next() (time.Duration, bool)
}
type FixedIntervalRetry struct {
	//重试间隔
	Interval time.Duration
	//重试最大次数
	Max int
	//已重试
	Cnt int
}

func (f *FixedIntervalRetry) Next() (time.Duration, bool) {
	f.Cnt++
	return f.Interval, f.Cnt <= f.Max
}
