package enums

// NodeType 定义节点类型
type NodeType int

const (
	NodeCommon NodeType = iota
	NodeSwitch
	NodeBoolean
	NodeFor
	NodeIterator
	NodeScript
	NodeSwitchScript
	NodeBooleanScript
	NodeForScript
	NodeFallback
)

// NodeTypeInfo 存储每种节点类型的额外信息
type NodeTypeInfo struct {
	Name        string
	Description string
	IsScript    bool
}

// NodeTypeMap 定义所有节点类型的信息
var NodeTypeMap = map[NodeType]NodeTypeInfo{
	NodeCommon:        {"common", "普通", false},
	NodeSwitch:        {"switch", "选择", false},
	NodeBoolean:       {"boolean", "布尔", false},
	NodeFor:           {"for", "循环次数", false},
	NodeIterator:      {"iterator", "循环迭代", false},
	NodeScript:        {"script", "脚本", true},
	NodeSwitchScript:  {"switch_script", "选择脚本", true},
	NodeBooleanScript: {"boolean_script", "布尔脚本", true},
	NodeForScript:     {"for_script", "循环次数脚本", true},
	NodeFallback:      {"fallback", "降级", false},
}

// GetNodeTypeInfo 根据 NodeType 获取信息
func GetNodeTypeInfo(nodeType NodeType) (NodeTypeInfo, bool) {
	info, exists := NodeTypeMap[nodeType]
	return info, exists
}
