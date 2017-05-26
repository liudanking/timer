# timer

A high precision timer base on timing-wheel written in Go


## Usage

```go
package main

import (
	"log"
	"time"

	"github.com/liudanking/timer"
)

func main() {
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
					log.Printf("job %d with delay %v triggerd, time diff: %v",
						i, delay, time.Now().Sub(createTime))
				}
			}(time.Now(), i, delay),
		)
	}

	time.Sleep(15 * time.Second)
	tw.Close()
}


```

```
2017/05/26 18:15:03 job 0 with delay 1s triggerd, time diff: 1.006292097s
2017/05/26 18:15:04 job 1 with delay 2s triggerd, time diff: 2.011096993s
2017/05/26 18:15:05 job 2 with delay 3s triggerd, time diff: 3.011222287s
2017/05/26 18:15:06 job 3 with delay 4s triggerd, time diff: 4.011262349s
2017/05/26 18:15:07 job 4 with delay 5s triggerd, time diff: 5.011292894s
2017/05/26 18:15:08 job 5 with delay 6s triggerd, time diff: 6.013259413s
2017/05/26 18:15:09 job 6 with delay 7s triggerd, time diff: 7.013191768s
2017/05/26 18:15:10 job 7 with delay 8s triggerd, time diff: 8.014250186s
2017/05/26 18:15:11 job 8 with delay 9s triggerd, time diff: 9.014169874s
2017/05/26 18:15:12 job 9 with delay 10s triggerd, time diff: 10.014251052s
2017/05/26 18:15:17 timingwheel exit
```
