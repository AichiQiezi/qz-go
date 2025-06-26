package model

import (
	"strings"
	"time"
)

type JobGroup struct {
	ID          int       // 主键
	AppName     string    // 执行器 App 名称
	Title       string    // 执行器名称
	AddressType int       // 执行器地址类型：0=自动注册、1=手动录入
	AddressList string    // 执行器地址列表（逗号分隔，手动录入）
	UpdateTime  time.Time // 更新时间

	registryList []string // 执行器地址列表（系统注册），内部字段
}

// GetRegistryList 获取执行器注册地址列表（从 AddressList 字段解析）
func (g *JobGroup) GetRegistryList() []string {
	if strings.TrimSpace(g.AddressList) != "" {
		g.registryList = strings.Split(g.AddressList, ",")
	}
	return g.registryList
}

// TableName 表名
func (g *JobGroup) TableName() string {
	return "xxl_job_group"
}
