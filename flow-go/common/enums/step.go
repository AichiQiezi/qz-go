package enums

// CmpStepTypeEnum 表示步骤的类型
type CmpStepTypeEnum int

const (
	START CmpStepTypeEnum = iota // iota 用于生成连续的整数值
	END
	SINGLE
)
