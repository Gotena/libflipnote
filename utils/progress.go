package utils

import (
	"time"

	"github.com/schollz/progressbar/v3"
)

func IncrementBar(progress *progressbar.ProgressBar, amount int) {
	go func() {
		for i := 0; i < amount; i++ {
			progress.Add(1)
			time.Sleep(time.Millisecond * 10)
		}
	}()
}
