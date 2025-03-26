package util

import (
	"fmt"
	"time"
)

// Stopwatch 用于计算操作的耗时
type Stopwatch struct {
	startTime   time.Time
	stopTime    time.Time
	running     bool
	elapsedTime time.Duration
}

// NewStopwatch 创建一个新的 Stopwatch 实例
func NewStopwatch() *Stopwatch {
	return &Stopwatch{}
}

// Start 开始计时
func (s *Stopwatch) Start() {
	if !s.running {
		s.startTime = time.Now().UTC()
		s.running = true
	}
}

// Stop 停止计时并计算总耗时
func (s *Stopwatch) Stop() {
	if s.running {
		s.stopTime = time.Now().UTC()
		s.elapsedTime = s.stopTime.Sub(s.startTime)
		s.running = false
	}
}

// Reset 重置 Stopwatch，清除已经计时的数据
func (s *Stopwatch) Reset() {
	s.startTime = time.Time{}
	s.stopTime = time.Time{}
	s.elapsedTime = 0
	s.running = false
}

// Elapsed 返回已计时的时间，单位为毫秒
func (s *Stopwatch) Elapsed() time.Duration {
	if s.running {
		// 如果计时器仍在运行，返回从开始到当前时间的差值
		return time.Since(s.startTime)
	}
	return s.elapsedTime
}

// ElapsedMilliseconds 返回已计时的时间，单位为毫秒
func (s *Stopwatch) ElapsedMilliseconds() int64 {
	return int64(s.Elapsed() / time.Millisecond)
}

// ElapsedSeconds 返回已计时的时间，单位为秒
func (s *Stopwatch) ElapsedSeconds() int64 {
	return int64(s.Elapsed() / time.Second)
}

// String 返回已计时的时间，格式化为 "xx ms"
func (s *Stopwatch) String() string {
	return fmt.Sprintf("%d ms", s.ElapsedMilliseconds())
}
