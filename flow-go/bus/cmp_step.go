package bus

import (
	"fmt"
	"qz-go/flow-go/common/enums"
	"time"
)

// CmpStep 是步骤的结构体
type CmpStep struct {
	NodeId            string
	NodeName          string
	Tag               string
	StepType          enums.CmpStepTypeEnum
	StartTime         time.Time
	EndTime           time.Time
	TimeSpent         int64 // 毫秒
	Success           bool
	Err               error
	RollbackTimeSpent int64 // 回滚消耗的时间，毫秒
}

// NewCmpStep 构造函数，用于初始化 CmpStep 对象
func NewCmpStep(nodeId, nodeName string, stepType enums.CmpStepTypeEnum) *CmpStep {
	return &CmpStep{
		NodeId:   nodeId,
		NodeName: nodeName,
		StepType: stepType,
	}
}

// BuildString 根据类型构建字符串
func (c *CmpStep) BuildString() string {
	if c.StepType == enums.SINGLE {
		if c.NodeName == "" {
			return fmt.Sprintf("{%s}", c.NodeId)
		}
		return fmt.Sprintf("{%s}[%s]", c.NodeId, c.NodeName)
	}
	// 目前没有其他的类型
	return ""
}

// BuildStringWithTime 根据类型和时间构建字符串
func (c *CmpStep) BuildStringWithTime() string {
	if c.StepType == enums.SINGLE {
		if c.NodeName == "" {
			if c.TimeSpent > 0 {
				return fmt.Sprintf("{%s}<%d>", c.NodeId, c.TimeSpent)
			}
			return fmt.Sprintf("{%s}", c.NodeId)
		}
		if c.TimeSpent > 0 {
			return fmt.Sprintf("{%s}[%s]<%d>", c.NodeId, c.NodeName, c.TimeSpent)
		}
		return fmt.Sprintf("{%s}[%s]", c.NodeId, c.NodeName)
	}
	// 目前没有其他的类型
	return ""
}

// BuildRollbackStringWithTime 回滚时间字符串
func (c *CmpStep) BuildRollbackStringWithTime() string {
	if c.StepType == enums.SINGLE {
		if c.NodeName == "" {
			if c.RollbackTimeSpent > 0 {
				return fmt.Sprintf("{%s}<%d>", c.NodeId, c.RollbackTimeSpent)
			}
			return fmt.Sprintf("{%s}", c.NodeId)
		}
		if c.RollbackTimeSpent > 0 {
			return fmt.Sprintf("{%s}[%s]<%d>", c.NodeId, c.NodeName, c.RollbackTimeSpent)
		}
		return fmt.Sprintf("{%s}[%s]", c.NodeId, c.NodeName)
	}
	// 目前没有其他的类型
	return ""
}

// Equals 判断两个 CmpStep 是否相等
func (c *CmpStep) Equals(other *CmpStep) bool {
	if other == nil {
		return false
	}
	return c.NodeId == other.NodeId
}
