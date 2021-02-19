package timeutil

import (
	"time"
)

func Now()time.Time  {
	return time.Now()
}

func Since(t time.Time)time.Duration{
	return time.Since(t)
}