package timer

import (
	"errors"
	"log"
	"sync"

	"time"
)

var (
	ErrWheelLevelInvalid = errors.New("timing wheel level invalid, should between [1, 7]")
	ErrSlotInvalid       = errors.New("slot index invalid")
	ErrTickDurInvalid    = errors.New("tick duration invalid")
	ErrDelayOverflow     = errors.New("time delay overflow")
)

type TimingWheel struct {
	sync.Mutex
	levelCnt int
	tickDur  time.Duration
	wheels   []*levelWheel
	exit     chan struct{}
}

func NewTimingWheel(levelCnt int, tickDur time.Duration) (*TimingWheel, error) {
	if levelCnt < 1 || levelCnt > 7 {
		return nil, ErrWheelLevelInvalid
	}
	if tickDur <= 0*time.Second {
		return nil, ErrTickDurInvalid
	}
	tw := &TimingWheel{
		levelCnt: levelCnt,
		tickDur:  tickDur,
		wheels:   make([]*levelWheel, levelCnt),
		exit:     make(chan struct{}),
	}
	for i := 0; i < levelCnt; i++ {
		tw.wheels[i] = newLevelWheel()
	}
	return tw, nil
}

func (tw *TimingWheel) AddFunc(delay time.Duration, f func()) error {
	tw.Lock()
	defer tw.Unlock()

	tc := uint64(delay / tw.tickDur)
	if tc > (1 << (uint(tw.levelCnt) * 8)) {
		return ErrDelayOverflow
	}

	var level int
	var slotIdx int
	levelSlotIdx := make([]uint8, tw.levelCnt)
	for i := 0; i < tw.levelCnt-1; i++ {
		levelSlotIdx[i] = tw.wheels[i].curPos + uint8(tc&0xff)
		tc = tc >> 8
	}

	job := &Job{
		levelSlotIdx: levelSlotIdx,
		f:            f,
		CreateTime:   time.Now(),
	}

	for i := tw.levelCnt - 1; i >= 0; i-- {
		if levelSlotIdx[i] != 0 {
			level = i
			slotIdx = int(levelSlotIdx[i])
			break
		}
	}
	log.Printf("levelSlotIdx:tc:%d %+v", uint64(delay/tw.tickDur), levelSlotIdx)

	return tw.addJob(level, slotIdx, job)
}

func (tw *TimingWheel) addJob(level, slotIdx int, job *Job) error {
	if level < 0 || level >= len(tw.wheels) {
		return ErrWheelLevelInvalid
	}
	if slotIdx < 0 || slotIdx >= SlotCnt {
		return ErrSlotInvalid
	}
	tw.wheels[level].addJob(slotIdx, job)
	return nil
}

func (tw *TimingWheel) Run() {
	tick := time.Tick(tw.tickDur)
	for {
		select {
		case <-tick:
			tw.tick()
		case <-tw.exit:
			log.Println("timingwheel exit")
			return
		}
	}
}

func (tw *TimingWheel) tick() {
	tw.Lock()
	defer tw.Unlock()

	// move position
	tw.wheels[0].curPos++

	// fire events
	tw.wheels[0].fireEvents()

	if tw.wheels[0].curPos != 0 {
		return
	}

	// turn wheels
	for i := 1; i < tw.levelCnt; i++ {
		tw.wheels[i].curPos++
		// move list to lower level
		l := tw.wheels[i].slots[tw.wheels[i].curPos]
		for e := l.Front(); e != nil; e = e.Next() {
			job := e.Value.(*Job)
			slotIdx := int(job.levelSlotIdx[i-1])
			tw.addJob(i-1, slotIdx, job)
		}
		l.Init()

		if tw.wheels[i].curPos != 0 {
			break
		}
	}
}

func (tw *TimingWheel) Close() {
	tw.exit <- struct{}{}
}
