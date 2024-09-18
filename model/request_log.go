package model

import "time"

type Requesting struct {
	StartTime  time.Time
	EndTime    time.Time
	Duration   time.Duration
	StatusCode int
	ClientIp   string
	Method     string
	Path       string
	UserAgent  string
}
