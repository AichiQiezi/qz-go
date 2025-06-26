package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func TestMultiGoFunc(t *testing.T) {
	go func() {
		go func() {
			time.Sleep(3 * time.Second)
			fmt.Println("halo ")
		}()

		time.Sleep(1 * time.Second)
	}()

	time.Sleep(5 * time.Second)
	fmt.Println("world")
}

func TestUpdateHotRow(t *testing.T) {
	// 连接数据库
	dsn := "root:aishuishui@tcp(127.0.0.1:3306)/qz"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("failed to connect to db: %v", err)
	}
	defer db.Close()

	// 初始化计数器
	_, _ = db.Exec(`UPDATE hot SET counter = 0 WHERE id = 1`)

	concurrency := 100
	var wg sync.WaitGroup
	wg.Add(concurrency)

	start := time.Now()

	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()

			// 热点更新语句
			_, err := db.Exec(`UPDATE hot SET counter = counter + 1 WHERE id = 1`)
			if err != nil {
				log.Printf("update failed: %v", err)
			}
		}()
	}

	wg.Wait()

	duration := time.Since(start)
	t.Logf("Total time: %v", duration)

	// 查询最终结果
	var finalCount int
	_ = db.QueryRow(`SELECT counter FROM hot WHERE id = 1`).Scan(&finalCount)
	t.Logf("Final counter: %d", finalCount)
}
