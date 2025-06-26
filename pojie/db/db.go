package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
)

// 日志文件
var logFile *os.File
var logger *log.Logger

// 初始化日志
func init() {
	dir, _ := os.Getwd()
	fmt.Println("当前执行目录:", dir)

	var err error
	logFile, err = os.OpenFile("./db_errors.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("无法创建日志文件:", err)
		os.Exit(1)
	}
	logger = log.New(logFile, "", log.LstdFlags)
}

// SaveAttempt 记录已尝试的密码
func SaveAttempt(qrRoute, password string, status int) {
	_, err := DB.Exec(`
        INSERT INTO attempts (qrcode_route, password, status, created_at, updated_at) 
        VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) 
        ON CONFLICT(qrcode_route, password) DO UPDATE 
        SET status = excluded.status, updated_at = CURRENT_TIMESTAMP`,
		qrRoute, password, status)
	if err != nil {
		LogDBError(fmt.Sprintf("写入数据库失败 (attempts): %s - %s - Status: %d - 错误: %v", qrRoute, password, status, err))
	}
}

// SaveCrackedQR 记录成功破解的 `qrroute`
func SaveCrackedQR(qrRoute, password string) {
	_, err := DB.Exec(`
        INSERT INTO cracked_qr (qrcode_route, password, created_at, updated_at) 
        VALUES (?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) 
        ON CONFLICT(qrcode_route) DO UPDATE 
        SET password = excluded.password, updated_at = CURRENT_TIMESTAMP`,
		qrRoute, password)
	if err != nil {
		LogDBError(fmt.Sprintf("写入数据库失败 (cracked_qr): %s - %s - 错误: %v", qrRoute, password, err))
	}
}

// LogDBError 记录数据库错误日志
func LogDBError(message string) {
	logMsg := fmt.Sprintf("[%s] %s", time.Now().Format("2006-01-02 15:04:05"), message)
	fmt.Println(logMsg)
	logger.Println(logMsg)
}

// IsCracked 检查是否已破解成功
func IsCracked(qrRoute string) bool {
	var exists string
	err := DB.QueryRow(`SELECT password FROM cracked_qr WHERE qrcode_route = ?`, qrRoute).Scan(&exists)
	return err == nil
}

// IsAttempted 检查是否已尝试过某个密码
func IsAttempted(qrRoute, password string) bool {
	var exists string
	err := DB.QueryRow(`SELECT password FROM attempts WHERE qrcode_route = ? AND password = ? and status != 0`, qrRoute, password).Scan(&exists)
	return err == nil
}

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./database.db")
	if err != nil {
		fmt.Println("数据库连接失败:", err)
		return
	}

	createTables()
}

func createTables() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS attempts (
			qrcode_route TEXT,
			password TEXT,
			status INTEGER,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (qrcode_route, password)
		);`,
		`CREATE TABLE IF NOT EXISTS cracked_qr (
			qrcode_route TEXT PRIMARY KEY,
			password TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
	}

	for _, query := range queries {
		_, err := DB.Exec(query)
		if err != nil {
			fmt.Println("创建表失败:", err)
		}
	}
}
