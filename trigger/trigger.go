package util

import (
	"context"
	"sync"
	"time"

	"github.com/zhanglp92/go-utils/background"
	"github.com/zhanglp92/go-utils/blocking"
	"github.com/zhanglp92/go-utils/xmath"
)

const (
	// 并发最大数量
	concurrenceConsumerMaxNum = 100
)

// ConsumerFunc 自定义消费方法
type ConsumerFunc func(ctx context.Context, buf []interface{})

// Trigger 延迟队列
type Trigger struct {
	sync.Mutex

	// 最大等待时间
	maxIdle time.Duration

	// 最大缓存消息个数
	maxMsgSize int

	// 消费方法(逐个都会消费)
	consumer []ConsumerFunc

	// 消息缓存(缓存复用)
	buf []interface{}

	// buf 当前消息长度
	bufSize int

	// 保证单线程写buf
	bufGuard chan interface{}

	// 并发消费协程最大数量
	blockConcurrence *blocking.LimitTask
}

// NewTrigger ...
func NewTrigger(ctx context.Context, maxIdle time.Duration, maxMsgSize int, consumers ...ConsumerFunc) *Trigger {
	// 纠正取值范围
	xmath.MaximinDuration(&maxIdle, 1*time.Second, 0, xmath.MinBorderLine)
	xmath.MaximinInt(&maxMsgSize, 10, 0, xmath.MinBorderLine)

	m := &Trigger{
		maxIdle:          maxIdle,
		maxMsgSize:       maxMsgSize,
		buf:              make([]interface{}, maxMsgSize),
		bufGuard:         make(chan interface{}),
		consumer:         consumers,
		blockConcurrence: blocking.NewLimitTask(1, concurrenceConsumerMaxNum),
	}

	// 启动后台任务
	m.runTask(ctx)

	return m
}

// AddConsumer 添加消费
func (m *Trigger) AddConsumer(consumer func(ctx context.Context, buf []interface{})) {
	if consumer != nil {
		m.consumer = append(m.consumer, consumer)
	}
}

// Enqueue 消息写队列
func (m *Trigger) Enqueue(o ...interface{}) {
	if o != nil {
		m.bufGuard <- o
	}
}

func (m *Trigger) runTask(ctx context.Context) {
	// 启动后台写任务
	m.writeTask(ctx)

	// 启动后台消费任务
	m.readTaks(ctx)
}

// 写入数据
func (m *Trigger) enqueue(ctx context.Context, o interface{}) {
	// 检查是否已经满了
	if m.bufSize >= m.maxMsgSize {
		m.dequenen(ctx)
	}

	m.Lock()
	defer m.Unlock()

	m.buf[m.bufSize] = o
	m.bufSize++
}

// 数据消费
func (m *Trigger) dequenen(ctx context.Context) {
	// 多写成消费
	m.blockConcurrence.Lock()
	go func(targetData []interface{}) {
		defer m.blockConcurrence.Unlock()
		for _, c := range m.consumer {
			c(ctx, targetData)
		}
	}(m.copyBuf())
}

func (m *Trigger) copyBuf() []interface{} {
	m.Lock()
	defer m.Unlock()

	// 复制消息使用多协程消费
	targetData := make([]interface{}, m.bufSize)
	copy(targetData, m.buf[:m.bufSize])

	// 缓存清空
	m.bufSize = 0

	return targetData
}

// 从chan写到buf
func (m *Trigger) writeTask(ctx context.Context) {
	background.NewTimerTask(ctx, 0, func(ctx context.Context) {
		select {
		case <-ctx.Done():
			break
		case o := <-m.bufGuard:
			m.enqueue(ctx, o)
		}
	})
}

// 运行后台任务
func (m *Trigger) readTaks(ctx context.Context) {
	background.NewTimerTask(ctx, m.maxIdle, m.dequenen)
}
