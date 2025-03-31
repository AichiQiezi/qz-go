package main

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestTryPassword2(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("无法获取当前工作目录: %v", err)
	}
	fmt.Println("当前执行目录:", dir)

	//InitDB()

	password2, err := tryPassword2("https://qr71.cn/orHY8H/qIaTz4X", "afd363")
	fmt.Println(password2)
	require.NoError(t, err)
}

func TestRetry(t *testing.T) {
	InitDB()
	retryFailedAttempts()
}
