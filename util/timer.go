package util

import (
	"fmt"
	"time"
)

func Timer(startCh, stopCh, finishedCh chan bool, duration time.Duration) {

	timer := time.NewTimer(time.Duration(duration))

	if !timer.Stop() {
		<-timer.C
	}

	for {
		select {
		case <-startCh:
			timer.Reset(time.Duration(duration))
		case <-stopCh:
			fmt.Printf("Timer stopped\n")
			timer.Stop()
		case <-timer.C:
			finishedCh <- true
		}
	}

}
