package background

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func Test_TimerTaskWithTimer(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	NewTimerTask(ctx, 1*time.Second, func(ctx context.Context) {
		fmt.Println(time.Now())
	})

	go func() {
		time.Sleep(5 * time.Second)
		cancel()
	}()
	time.Sleep(10 * time.Second)
}

func Test_TimerTaskWithoutTimer(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	NewTimerTask(ctx, 0, func(ctx context.Context) {
		fmt.Println(time.Now())
	})

	go func() {
		time.Sleep(5 * time.Second)
		cancel()
	}()
	time.Sleep(10 * time.Second)
}
