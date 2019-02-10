package cronjobs

import (
	"time"
)

func init() {
	ticker := time.NewTicker(time.Minute)

	go func() {
		for {
			select {
			case <-ticker.C:
				checkDataSourceNodes()
			}
		}
	}()
}
