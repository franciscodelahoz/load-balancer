package health

import "time"

type Config struct {
	Interval           time.Duration
	Timeout            time.Duration
	HealthyThreshold   int
	UnhealthyThreshold int
	Path               string
	Method             string
}
