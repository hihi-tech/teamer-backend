package main

import "time"

type TimeRange struct {
	Start *time.Time `json:"start"`
	End *time.Time `json:"end"`
}
