package executor

import (
	"qz-go/flow-go/core/component"
)

type NodeExecutor interface {
	Execute(nc component.NodeComponent) error
	Retry(nc component.NodeComponent, currentRetryCount int) error
}

type DefaultNodeExecutor struct {
}

func (ne *DefaultNodeExecutor) Execute(nc component.NodeComponent) error {
	retryCount := nc.GetRetryCount()

	for i := range retryCount {
		var err error
		if i != 0 {
			err = ne.Retry(nc, i) // 重试逻辑
		} else {
			err = nc.Execute()
		}

		if err == nil {
			break
		}

		retryCount++
	}
	return nil
}

func (ne *DefaultNodeExecutor) Retry(nc component.NodeComponent, currentRetryCount int) error {
	return nc.Execute()
}
