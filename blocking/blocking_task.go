package blocking

import (
	"fmt"
	"sync"

	"github.com/zhanglp92/go-utils/xmath"
)

// 阻塞任务

// LimitTask 最大任务数量阻塞
type LimitTask struct {
	// 扩容锁
	expansionMutex sync.Mutex

	// 通道数量和最大容量
	free, cap   int
	empty, full chan struct{}
}

// NewLimitTask size:最大任务容量
func NewLimitTask(size int, cap int) (m *LimitTask) {
	// 默认一个通道
	xmath.MaximinInt(&size, 1, 0, xmath.MinBorderLine)
	// 最小容量为size
	xmath.MaximinInt(&cap, size, 0, xmath.MinBorderLine)

	defer func() { m.Expansion(size) }()
	return &LimitTask{
		cap:   cap,
		free:  cap - size,
		empty: make(chan struct{}, cap),
		full:  make(chan struct{}, cap),
	}
}

// Lock 尝试阻塞
func (m *LimitTask) Lock() {
	m.full <- <-m.empty
}

// Unlock 释放资源
func (m *LimitTask) Unlock() {
	m.empty <- <-m.full
}

// Expansion 扩大容量
func (m *LimitTask) Expansion(n int) error {
	for i := 0; i < n; i++ {
		if err := m.expansionInc(); err != nil {
			return err
		}
	}
	return nil
}

// expansion 扩大容量
func (m *LimitTask) expansionInc() error {
	m.expansionMutex.Lock()
	defer m.expansionMutex.Unlock()

	if m.free <= 0 {
		return fmt.Errorf("more than max capacity[%v]", m.cap)
	}

	m.empty <- struct{}{}
	m.free--
	return nil
}
