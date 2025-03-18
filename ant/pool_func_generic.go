package ant

// PoolWithFuncGeneric is the generic version of PoolWithFunc.
type PoolWithFuncGeneric[T any] struct {
	*poolCommon

	// fn is the unified function for processing tasks.
	fn func(T)
}
