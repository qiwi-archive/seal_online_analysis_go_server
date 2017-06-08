package innotrio

import (
	"time"
	"fmt"
)

type Tracker struct {
	time time.Time
}

func NewTracker() (*Tracker) {
	return &Tracker{time.Now()}
}

func (self *Tracker) Log(name string, args ...interface{}) {
	elapsed := time.Since(self.time).String()
	fmt.Println(name, elapsed, args)
	//print("%s took %s", name, elapsed)
}

func (self *Tracker) GetSeconds() (float64) {
	return time.Since(self.time).Seconds()
}