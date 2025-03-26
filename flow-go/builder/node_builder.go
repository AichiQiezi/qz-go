package builder

import (
	"qz-go/flow-go/common/enums"
	"qz-go/flow-go/core/element"
)

type NodeBuilder struct {
	node     *element.Node
	nodeType enums.NodeType
}

func newNodeBuilder(nodeType enums.NodeType) *NodeBuilder {
	return &NodeBuilder{
		node:     new(element.Node),
		nodeType: nodeType,
	}
}

// NewCommonNode 创建通用节点
func NewCommonNode() *NodeBuilder {
	return newNodeBuilder(enums.NodeCommon)
}

// NewSwitchNode 创建 Switch 类型节点
func NewSwitchNode() *NodeBuilder {
	return newNodeBuilder(enums.NodeSwitch)
}

// NewBooleanNode 创建 Boolean 类型节点
func NewBooleanNode() *NodeBuilder {
	return newNodeBuilder(enums.NodeBoolean)
}

// NewForNode 创建 For 类型节点
func NewForNode() *NodeBuilder {
	return newNodeBuilder(enums.NodeFor)
}

// NewIteratorNode 创建 Iterator 类型节点
func NewIteratorNode() *NodeBuilder {
	return newNodeBuilder(enums.NodeIterator)
}

// Build 返回构造的 Node
func (b *NodeBuilder) Build() *element.Node {
	return b.node
}
