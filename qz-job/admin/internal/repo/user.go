package model

import (
	"strconv"
	"strings"
)

type User struct {
	ID         int
	Username   string
	Password   string
	Role       int    // 0-普通用户、1-管理员
	Permission string // 逗号分割的执行器ID列表
}

// ValidPermission 判断当前用户是否有权限访问指定 jobGroup
func (u *User) ValidPermission(jobGroup int) bool {
	if u.Role == 1 {
		return true
	}
	if strings.TrimSpace(u.Permission) != "" {
		permissions := strings.Split(u.Permission, ",")
		jobGroupStr := strconv.Itoa(jobGroup)
		for _, p := range permissions {
			if strings.TrimSpace(p) == jobGroupStr {
				return true
			}
		}
	}
	return false
}

// TableName 表名
func (u *User) TableName() string {
	return "xxl_job_user"
}
