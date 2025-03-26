package bus

// Slot 请求下上文的包装类
type Slot struct {
	steps       []*CmpStep
	routeResult bool
	//condition []*con
	metaData map[string]interface{}
}

func NewSlot() *Slot {
	return &Slot{}
}

func (s *Slot) AddStep(step *CmpStep) {
	s.steps = append(s.steps, step)
}

func (s *Slot) PutRequestId(id string) {
	s.metaData["_req_id"] = id
}

func (s *Slot) PutParams(params interface{}) {
	s.metaData["_params"] = params
}
