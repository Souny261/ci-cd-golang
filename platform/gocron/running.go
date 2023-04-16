package gocron

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
)

func RunningService() {
	s := gocron.NewScheduler(time.Local)
	fmt.Println("Start Running Service")
	s.Every(3).Seconds().Do(func() {
		fmt.Println("Running.....")
	})
}
