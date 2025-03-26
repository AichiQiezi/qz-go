package component

import (
	"qz-go/flow-go/bus"
	"qz-go/flow-go/common/enums"
	"qz-go/flow-go/util"
	"time"
)

// NodeComponent 接口
type NodeComponent interface {
	Precess() error
	IsAccess() bool
	IsEnd() bool
	GetSlot() bus.Slot
	GetSlotIndex()
	BeforeProcess()
	AfterProcess()
	OnSuccess()
	OnError()
	GetRetryCount() int
	Execute() error
}

type nodeCommon struct {
	NodeComponent

	nodeId   string
	name     string
	nodeType enums.NodeType

	retryCount int
}

func (n *nodeCommon) GetName() string {
	return n.name
}

func (n *nodeCommon) Execute() error {
	stopwatch := util.NewStopwatch()
	stopwatch.Start()

	cmpStep := bus.NewCmpStep(n.nodeId, n.name, enums.CmpStepTypeEnum(2))
	defer func() {
		// 后置处理
		n.AfterProcess()

		stopwatch.Stop()
		cmpStep.EndTime = time.Now().UTC()
		cmpStep.TimeSpent = stopwatch.ElapsedMilliseconds()
	}()

	st := n.GetSlot()
	cmpStep.StartTime = time.Now().UTC()
	st.AddStep(cmpStep)

	n.BeforeProcess()
	err := n.Precess()

	if err != nil {
		n.OnError()
	} else {
		cmpStep.Success = false
		cmpStep.Err = err
		n.OnSuccess()
	}

	return err
}

// finalizeExecution 负责结束时的逻辑，包括异常处理、计时停止、后置处理等
//func (n *nodeCommon) finalizeExecution(stopwatch *util.Stopwatch, cmpStep *entity.CmpStep) {
//	// 捕获异常并处理成功或失败
//	if err := recover(); err == nil {
//		n.onSuccess()
//	} else {
//		n.onError()
//	}
//
//	// 后置处理
//	n.afterProcess()
//
//	// 结束计时
//	stopwatch.Stop()
//	cmpStep.EndTime = time.Now().UTC()
//	cmpStep.TimeSpent = stopwatch.ElapsedMilliseconds()
//}
