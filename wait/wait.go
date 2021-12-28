package wait

import (
	"time"

	"context"
)

var (
	NeverStop = context.TODO()
	NoThing   = func() { return }
)

// Until 不断的以间隔period运行f, 直到ctx.Done 返回.
// 间隔计算f函数运行时间.
func Until(ctx context.Context, f func(), period time.Duration) {
	BackoffUntil(ctx, f, NewSimpleBackoffManager(period), true)
}

// NonSlidingUntil 不断的以间隔period运行f, 直到ctx.Done 返回.
// 忽视f函数运行时间.
func NonSlidingUntil(ctx context.Context, f func(), period time.Duration) {
	BackoffUntil(ctx, f, NewSimpleBackoffManager(period), false)
}

// Forever 不断的以间隔period运行f.
// sliding 为true 则计算f函数运行时间.sliding 为false 则忽视f运行时间.
func Forever(f func(), period time.Duration, sliding bool) {
	BackoffUntil(NeverStop, f, NewSimpleBackoffManager(period), sliding)
}

// BackoffUntil 不断的运行 f 函数，运行间隔由BackoffManager 提供, 直到ctx.Done 返回.
// sliding 为true 则计算f函数运行时间.sliding 为false 则忽视f运行时间.
func BackoffUntil(ctx context.Context, f func(), backoff BackoffManager, sliding bool) {
	var t *time.Timer
	stopCh := ctx.Done()
	for {
		select {
		case <-stopCh:
			return
		default:
		}

		if !sliding {
			t = backoff.Backoff()
		}

		func() {
			f()
		}()

		if sliding {
			t = backoff.Backoff()
		}

		select {
		case <-stopCh:
			if !t.Stop() {
				<-t.C
			}
			return
		case <-t.C:
		}
	}
}

// BackoffManager 接口可以提供 Backoff 方法，
// Backoff 返回一个Timer 确定要等待的时间
type BackoffManager interface {
	Backoff() *time.Timer
}

// simpleBackoffManager 是一个简单的BackoffManager 以固定时间间隔返回Backoff
type simpleBackoffManager struct {
	duration time.Duration
	timer    *time.Timer
}

func (sb *simpleBackoffManager) getNextBackoff() time.Duration {
	return sb.duration
}

func (sb *simpleBackoffManager) Backoff() *time.Timer {
	backoff := sb.getNextBackoff()
	if sb.timer == nil {
		sb.timer = time.NewTimer(backoff)
	} else {
		sb.timer.Reset(backoff)
	}
	return sb.timer
}

// NewSimpleBackoffManager 返回一个简单的BackoffManager 以固定时间间隔返回Backoff
func NewSimpleBackoffManager(period time.Duration) BackoffManager {
	bf := &simpleBackoffManager{}
	bf.duration = period
	return bf
}
