package nativeFunctions

import "time"

func timeInSeconds() float64 {
	return float64(time.Now().UnixNano()) / 1e9
}

var GlobalFunctions = map[string]interface{}{
	"clock": timeInSeconds(),
}
