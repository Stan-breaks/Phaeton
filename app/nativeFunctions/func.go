package nativefunctions

import "time"

func timeInSeconds() float32 {
	now := time.Now()
	unixTime := float32(now.UnixNano())
	return unixTime
}

var GlobalFunctions = map[string]interface{}{
	"clock": timeInSeconds(),
}
