package element

import (
	"qz-go/flow-go/common/enums"
	"qz-go/flow-go/common/errors"
	"qz-go/flow-go/core/component"
	"qz-go/flow-go/core/executor"
	"sync/atomic"
)

// Node 节点
type Node struct {
	Id       string
	Name     string
	Type     enums.NodeType
	instance component.NodeComponent

	access atomic.Bool
}

func (n *Node) execute(slowIndex int) (err error) {
	if n.access.Load() || n.instance.IsAccess() {
		// todo 要改成单例
		ne := &executor.DefaultNodeExecutor{}
		err = ne.Execute(n.instance)
	}

	if n.instance.IsEnd() {
		return errors.NewErrChainEnd("chain to the end")
	}

	return err
}
