package main

import (
	"fmt"
	"sync"
	"time"
)

// 定义一个结构体类型
type Item struct {
	ID    int
	Title string
}

func processItem(item Item) {
	fmt.Printf("开始处理 ID=%d, Title=%s\n", item.ID, item.Title)
	time.Sleep(500 * time.Millisecond) // 模拟耗时操作
}

func main() {
	// 创建结构体数组
	items := []Item{
		{ID: 1, Title: "苹果"},
		{ID: 2, Title: "香蕉"},
		{ID: 3, Title: "橘子"},
	}

	var wg sync.WaitGroup

	for _, item := range items {
		//显式拷贝循环变量，避免闭包捕获陷阱
		//itemCopy := item

		wg.Add(1)
		go func() {
			defer wg.Done()
			processItem(item)
		}()
	}

	wg.Wait()
	fmt.Println("所有结构体元素处理完成")
}
