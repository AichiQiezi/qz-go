package main

import (
	"fmt"
	"time"
)

// Msg 是传递给每个 Action 的消息
type Msg struct {
	Content string
}

// Action 接口，每个 Action 必须实现 Call 方法
type Action interface {
	Call(queue string, message *Msg, next func() bool) bool
}

// PrintAction 是一个实现了 Action 接口的结构
type PrintAction struct {
	Name   string
	Cancel bool // 模拟某个步骤希望中断链
}

func (p PrintAction) Call(queue string, message *Msg, next func() bool) bool {
	fmt.Printf("Action: %s | Queue: %s | Msg: %s\n", p.Name, queue, message.Content)

	if p.Cancel {
		fmt.Println("-> 中断链条: ", p.Name)
		return false
	}

	return next()
}

// continuation 构建责任链，支持中断
func continuation(actions []Action, queue string, message *Msg, final func()) func() bool {
	next := func() bool {
		final()
		return true
	}

	for i := len(actions) - 1; i >= 0; i-- {
		act := actions[i]
		currentNext := next
		next = func(a Action, cn func() bool) func() bool {
			return func() bool {
				return a.Call(queue, message, cn)
			}
		}(act, currentNext)
	}

	return next
}

func main() {
	var formattedExpiration string
	parsedExpTime, err := time.Parse("2006-01-02T15:04:05Z", "2025-05-13T00:00:00Z")
	if err == nil {
		formattedExpiration = parsedExpTime.Format("2006-01-02")
	} else {
		formattedExpiration = "cnm"
	}

	fmt.Println(formattedExpiration)
}
