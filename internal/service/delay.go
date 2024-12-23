package service

import (
	"time"

	"github.com/charmbracelet/log"
	"golang.org/x/exp/rand"
)

type Delay [2]time.Duration

func (d *Delay) Take() {
	if d == nil || d[0] == 0 && d[1] == 0 {
		return
	}
	// random duration between d[0] and d[1]
	duration := d[0] + time.Duration(rand.Int63n(int64(d[1]-d[0])))
	log.Info("🕰 delaying request", "duration", duration)
	time.Sleep(duration)
}

func (d *Delay) String() string {
	return d[0].String() + "-" + d[1].String()
}
