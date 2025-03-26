package element

type Condition struct {
}

func (c *Condition) execute(slotIndex int) error {
	//bus := bus.NewDataBus(100)

	executeCondition(slotIndex)
	return nil
}

func executeCondition(slotIndex int) {

}
