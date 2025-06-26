package model

import "time"

type JobRegistry struct {
	ID            int       `json:"id"`
	RegistryGroup string    `json:"registry_group"`
	RegistryKey   string    `json:"registry_key"`
	RegistryValue string    `json:"registry_value"`
	UpdateTime    time.Time `json:"update_time"`
}

// TableName 表名
func (r *JobRegistry) TableName() string {
	return "xxl_job_registry"
}
