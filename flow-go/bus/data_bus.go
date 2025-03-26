package bus

import (
	"sync"
	"sync/atomic"
)

// DataBus 管理 bus
type DataBus struct {
	slots *sync.Map
	queue *SafeQueue[int]

	occupyCount atomic.Int32 // bus 数量
}

func NewDataBus() *DataBus {
	//if size <= 0 {
	//	size = 1024
	//}

	queue := NewSafeQueue[int]()
	for i := 0; i < 1024; i++ {
		queue.Enqueue(i) // 预填充队列
	}

	return &DataBus{
		slots:       &sync.Map{},
		queue:       queue,
		occupyCount: atomic.Int32{},
	}
}

func (db *DataBus) GetSlot(slotIndex int) *Slot {
	value, ok := db.slots.Load(slotIndex)
	if !ok {
		return nil // 如果不存在，返回 nil
	}

	// 类型断言，将 `any` 转换为 `*Slot`
	slot, _ := value.(*Slot)
	return slot
}

func (db *DataBus) OfferIndex(slot *Slot) (*Slot, int) {
	slotIndex, ok := db.queue.TryDequeue() // 快速尝试
	if ok {
		db.slots.Store(slotIndex, slot)
		db.occupyCount.Add(1)

		return slot, slotIndex
	}

	slotIndex = db.queue.Dequeue()
	db.slots.Store(slotIndex, slot)
	db.occupyCount.Add(1)

	return slot, slotIndex
}

func (db *DataBus) ReleaseIndex(slotIndex int) {
	db.slots.Delete(slotIndex)
	db.queue.Enqueue(slotIndex)
	db.occupyCount.Add(-1)
}
