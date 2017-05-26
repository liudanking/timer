package timer

import (
	"log"
	"testing"
	"time"

	"github.com/liudanking/timer"
)

func TestTimingWheel(t *testing.T) {
	tw, err := NewTimingWheel(1 * time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	go tw.Run()
	time.Sleep(100 * time.Millisecond)
	start := time.Now()
	for i := 15; i < 20; i++ {
		delay := time.Duration(i) * time.Second
		tw.AddFunc(delay, func() {
			log.Printf("job func triggerd, time diff %dms", time.Now().Sub(start)/1e6)
		})
		log.Printf("%d added", i)
	}

	time.Sleep(20 * time.Second)

	for i := 5; i < 10; i++ {
		delay := time.Duration(i) * time.Second
		tw.AddFunc(delay, func() {
			log.Printf("job func triggerd, time diff %dms", time.Now().Sub(start)/1e6)
		})
		log.Printf("%d added", i)
	}

	time.Sleep(7 * time.Second)

	tw.Close()

}

func TestTimingWheel2(t *testing.T) {
	tw, err := timer.NewTimingWheel(1 * time.Millisecond)
	if err != nil {
		log.Fatal(err)
	}
	go tw.Run()

	for i := 0; i < 10; i++ {
		delay := time.Duration(i+1) * time.Second
		tw.AddFunc(delay,
			func(createTime time.Time, i int, delay time.Duration) func() {
				return func() {
					log.Printf("job %d with delay %v triggerd, time diff: %v", i, delay, time.Now().Sub(createTime))
				}
			}(time.Now(), i, delay),
		)
	}

	time.Sleep(15 * time.Second)
	tw.Close()
}
