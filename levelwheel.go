package timer

import (
	"container/list"
	"time"
)

const (
	SlotCnt  = 256
	LevelCnt = 7
)

type levelWheel struct {
	curPos uint8
	slots  []*list.List
}

type Job struct {
	levelSlotIdx []uint8 // level: slot
	f            func()
	CreateTime   time.Time
}

func newLevelWheel() *levelWheel {
	lw := &levelWheel{
		curPos: 0,
		slots:  make([]*list.List, SlotCnt),
	}

	for i := 0; i < SlotCnt; i++ {
		lw.slots[i] = list.New()
	}

	return lw
}

func (lw *levelWheel) addJob(slotIdx int, job *Job) error {
	if slotIdx < 0 || slotIdx >= SlotCnt {
		return ErrSlotInvalid
	}
	lw.slots[slotIdx].PushBack(job)
	return nil
}

func (lw *levelWheel) fireEvents() {
	l := lw.slots[lw.curPos]
	if l == nil {
		return
	}
	for e := l.Front(); e != nil; e = e.Next() {
		job := e.Value.(*Job)
		go job.f()
	}

	lw.slots[lw.curPos].Init()
}
