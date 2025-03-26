package util

import (
	"testing"
	"time"
)

// TestStopwatch_Start 测试 Stopwatch 的 Start 方法
func TestStopwatch_Start(t *testing.T) {
	sw := NewStopwatch()
	sw.Start()

	// 验证 stopwatch 是否已开始计时
	if !sw.running {
		t.Errorf("Expected stopwatch to be running, but it is not.")
	}
}

// TestStopwatch_Stop 测试 Stopwatch 的 Stop 方法
func TestStopwatch_Stop(t *testing.T) {
	sw := NewStopwatch()
	sw.Start()
	time.Sleep(1 * time.Second)
	sw.Stop()

	// 检查 stopwatch 是否停止了
	if sw.running {
		t.Errorf("Expected stopwatch to be stopped, but it is still running.")
	}

	// 确认已计时的时间大于0
	if sw.Elapsed() <= 0 {
		t.Errorf("Expected elapsed time to be greater than 0, but it is %v", sw.Elapsed())
	}
}

// TestStopwatch_Elapsed 测试 Stopwatch 获取经过时间的方法
func TestStopwatch_Elapsed(t *testing.T) {
	sw := NewStopwatch()
	sw.Start()
	time.Sleep(500 * time.Millisecond)
	sw.Stop()

	// 确保经过的时间大于或等于500毫秒
	if sw.Elapsed() < 500*time.Millisecond {
		t.Errorf("Expected elapsed time to be at least 500ms, but it was %v", sw.Elapsed())
	}
}

// TestStopwatch_ElapsedMilliseconds 测试获取经过时间（毫秒）
func TestStopwatch_ElapsedMilliseconds(t *testing.T) {
	sw := NewStopwatch()
	sw.Start()
	time.Sleep(100 * time.Millisecond)
	sw.Stop()

	// 确保返回的时间至少为100毫秒
	elapsed := sw.ElapsedMilliseconds()
	if elapsed < 100 {
		t.Errorf("Expected elapsed time in milliseconds to be at least 100, but got %d", elapsed)
	}
}

// TestStopwatch_Reset 测试 Stopwatch 的 Reset 方法
func TestStopwatch_Reset(t *testing.T) {
	sw := NewStopwatch()
	sw.Start()
	time.Sleep(1 * time.Second)
	sw.Stop()

	// 重置计时器
	sw.Reset()

	// 确认计时器已重置
	if sw.Elapsed() != 0 {
		t.Errorf("Expected stopwatch elapsed time to be 0 after reset, but got %v", sw.Elapsed())
	}
	if sw.running {
		t.Errorf("Expected stopwatch to be not running after reset, but it is running.")
	}
}

// TestStopwatch_String 测试 Stopwatch 的 String 方法
func TestStopwatch_String(t *testing.T) {
	sw := NewStopwatch()
	sw.Start()
	time.Sleep(200 * time.Millisecond)
	sw.Stop()

	result := sw.String()
	// 确认返回的字符串格式是有效的
	if result == "" {
		t.Errorf("Expected stopwatch string to be a non-empty string, but got an empty string.")
	}
}
