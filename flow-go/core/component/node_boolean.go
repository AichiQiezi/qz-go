package component

type INodeBooleanComponent interface {
	ProcessBoolean() (bool, error)
}

type NodeBooleanComponent struct {
	*nodeCommon
}
