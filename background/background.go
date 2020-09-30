package background

import (
	"context"
	"time"
)

// 后台运行工具(ctx 结束后自动退出)

// 自定义后台方法签名
type userTaskType func(context.Context)

// 任务方法签名(return: 是否停止循环, ture停止)
type taskType func(context.Context) bool

// TimerTask 定时条件后台任务
type TimerTask struct {
	userTask userTaskType
	timer    time.Duration
}

// NewTimerTask 创建定时条件触发任务
func NewTimerTask(ctx context.Context, timer time.Duration, task userTaskType) (m *TimerTask) {
	defer func() { RunLoopTask(ctx, m.buildRun(ctx)) }()
	return &TimerTask{userTask: task, timer: timer}
}

// 构造run任务
func (m *TimerTask) buildRun(ctx context.Context) taskType {
	if m.timer > 0 {
		return m.runWithTimer
	}
	return m.runWithoutTimer
}

// 不带定时器
func (m *TimerTask) runWithoutTimer(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		m.userTask(ctx)
		return true
	default:
		m.userTask(ctx)
		return false
	}
}

// 带定时器
func (m *TimerTask) runWithTimer(ctx context.Context) bool {
	timer := time.NewTimer(m.timer)
	select {
	case <-ctx.Done():
		m.userTask(ctx)
		return true
	case <-timer.C:
		m.userTask(ctx)
		return false
	}
}

// RunTask 后台执行
func RunTask(ctx context.Context, task userTaskType) {
	go func() {
		task(ctx)
	}()
}

// RunLoopTask 后台循环执行任务
func RunLoopTask(ctx context.Context, task taskType) {
	RunTask(ctx, func(ctx context.Context) {
		for !task(ctx) {
		}
	})
}
