package model

import "time"

type JobInfo struct {
	ID                     int       // 主键ID
	JobGroup               int       // 执行器主键ID
	JobDesc                string    // 任务描述
	AddTime                time.Time // 添加时间
	UpdateTime             time.Time // 更新时间
	Author                 string    // 负责人
	AlarmEmail             string    // 报警邮件
	ScheduleType           string    // 调度类型
	ScheduleConf           string    // 调度配置（依赖调度类型）
	MisfireStrategy        string    // 调度过期策略
	ExecutorRouteStrategy  string    // 执行器路由策略
	ExecutorHandler        string    // 执行器任务 Handler 名称
	ExecutorParam          string    // 执行器任务参数
	ExecutorBlockStrategy  string    // 阻塞处理策略
	ExecutorTimeout        int       // 任务执行超时时间（秒）
	ExecutorFailRetryCount int       // 失败重试次数
	GlueType               string    // GLUE 类型
	GlueSource             string    // GLUE 源代码
	GlueRemark             string    // GLUE 备注
	GlueUpdateTime         time.Time // GLUE 更新时间
	ChildJobID             string    // 子任务 ID（多个逗号分隔）
	TriggerStatus          int       // 调度状态：0-停止，1-运行
	TriggerLastTime        int64     // 上次调度时间（时间戳）
	TriggerNextTime        int64     // 下次调度时间（时间戳）
}

// TableName 表名
func (g *JobInfo) TableName() string {
	return "xxl_job_info"
}
