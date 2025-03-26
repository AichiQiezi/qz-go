package bus

import "qz-go/flow-go/core/element"

type FlowBus struct {
	chains map[string]*element.Chain
}

func NewFlowBus() *FlowBus {
	return &FlowBus{make(map[string]*element.Chain)}
}

func (fb *FlowBus) GetChain(chainId string) *element.Chain {
	return fb.chains[chainId]
}
