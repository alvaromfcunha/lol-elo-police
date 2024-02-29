package police

import (
	"time"
)

func setInterval(cb func(), interval time.Duration) chan<- bool {
	ticker := time.NewTicker(interval)
	stopChan := make(chan bool)

	go func() {
		for {
			select {
			case <-stopChan:
				return
			case <-ticker.C:
				cb()
			}
		}
	}()

	return stopChan
}
