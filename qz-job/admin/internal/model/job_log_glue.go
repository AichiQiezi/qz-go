package model

import "time"

// JobLogGlue represents glue code version history for a job
type JobLogGlue struct {
	ID         int       `json:"id"`          // 主键ID
	JobID      int       `json:"job_id"`      // 任务主键ID
	GlueType   string    `json:"glue_type"`   // GLUE类型（如 BEAN, GLUE_GROOVY 等）
	GlueSource string    `json:"glue_source"` // 脚本源码
	GlueRemark string    `json:"glue_remark"` // 修改备注
	AddTime    time.Time `json:"add_time"`    // 创建时间
	UpdateTime time.Time `json:"update_time"` // 更新时间
}

// TableName 表名
func (lg *JobLogGlue) TableName() string {
	return "xxl_job_log_glue"
}
