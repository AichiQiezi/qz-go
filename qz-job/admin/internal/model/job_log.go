package model

import "time"

// JobLog represents the log entry of a scheduled job execution
type JobLog struct {
	ID                     int64     `json:"id"`                        // 主键ID
	JobGroup               int       `json:"job_group"`                 // 任务组ID
	JobID                  int       `json:"job_id"`                    // 任务ID
	ExecutorAddress        string    `json:"executor_address"`          // 执行器地址
	ExecutorHandler        string    `json:"executor_handler"`          // 执行器方法
	ExecutorParam          string    `json:"executor_param"`            // 执行参数
	ExecutorShardingParam  string    `json:"executor_sharding_param"`   // 分片参数
	ExecutorFailRetryCount int       `json:"executor_fail_retry_count"` // 失败重试次数
	TriggerTime            time.Time `json:"trigger_time"`              // 触发时间
	TriggerCode            int       `json:"trigger_code"`              // 触发状态码
	TriggerMsg             string    `json:"trigger_msg"`               // 触发日志
	HandleTime             time.Time `json:"handle_time"`               // 执行完成时间
	HandleCode             int       `json:"handle_code"`               // 执行状态码
	HandleMsg              string    `json:"handle_msg"`                // 执行日志
	AlarmStatus            int       `json:"alarm_status"`              // 告警状态
}

// TableName 表名
func (l *JobLog) TableName() string {
	return "xxl_job_log"
}
