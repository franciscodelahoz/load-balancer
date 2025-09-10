package health

import "time"

type Config struct {
	Interval         time.Duration
	Timeout          time.Duration
	Path             string
	Method           string
	SuccessThreshold int
	FailureThreshold int
}
