package main

import (
	"time"
)

type animData struct {
	next    time.Time
	forward bool
}

var animDataMap map[any]*animData

func init() {
	animDataMap = map[any]*animData{}
}

func anim(key any, duration time.Duration, from float32, to float32) float32 {
	if _, exist := animDataMap[key]; !exist {
		animDataMap[key] = &animData{
			next:    time.Now(),
			forward: false,
		}
	}

	data := animDataMap[key]

	if time.Now().After(data.next) {
		data.next = time.Now().Add(duration)
		data.forward = !data.forward
	}

	progress := float32(data.next.Sub(time.Now())) / float32(duration)
	if data.forward {
		progress = 1 - progress
	}

	return from + ((to - from) * progress)
}
