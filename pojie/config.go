package main

//
//import (
//	"database/sql"
//	"fmt"
//	_ "github.com/mattn/go-sqlite3"
//)
//
//var DB *sql.DB
//
//func InitDB() {
//	var err error
//	DB, err = sql.Open("sqlite3", "./pojie/database.db")
//	if err != nil {
//		fmt.Println("数据库连接失败:", err)
//		return
//	}
//
//	createTables()
//}
//
//func createTables() {
//	queries := []string{
//		`CREATE TABLE IF NOT EXISTS attempts (
//			qrcode_route TEXT,
//			password TEXT,
//			status INTEGER,
//			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
//			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
//			PRIMARY KEY (qrcode_route, password)
//		);`,
//		`CREATE TABLE IF NOT EXISTS cracked_qr (
//			qrcode_route TEXT PRIMARY KEY,
//			password TEXT,
//			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
//			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
//		);`,
//	}
//
//	for _, query := range queries {
//		_, err := DB.Exec(query)
//		if err != nil {
//			fmt.Println("创建表失败:", err)
//		}
//	}
//}
