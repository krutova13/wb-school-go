package l25

import "time"

func Sleep(duration time.Duration) {
	ch := make(chan struct{})
	go func() {
		timer := time.NewTimer(duration)
		<-timer.C
		close(ch)
	}()
	<-ch
}
