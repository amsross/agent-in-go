package main

import (
	"time"
)

func getCurrentTime() string {
	return time.Now().Format(time.RFC3339)
}
