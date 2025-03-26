package executor

import (
	"qz-go/flow-go/bus"
	"qz-go/flow-go/common/errors"
	"qz-go/flow-go/property"
)

// FlowExecutor 流程规则主要执行器
type FlowExecutor struct {
	property.FlowConfig
}

func (fe *FlowExecutor) execute(chainId string, requestId string, param interface{}) error {
	slot, index := bus.NewDataBus().OfferIndex(bus.NewSlot())
	if slot == nil {
		return errors.NewSlotNotFound("slot not found")
	}

	slot.PutRequestId(requestId)
	slot.PutParams(param)

	chain := bus.NewFlowBus().GetChain(chainId)
	if chain == nil {
		return errors.NewChainNotFound("chain not found")
	}

	err := chain.Execute(index)
	if err != nil {
	}
	return nil
}
