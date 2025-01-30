package nativeFunctions

import "time"

func timeInSeconds() float64 {
	return float64(time.Now().Unix())
}

var GlobalFunctions = map[string]interface{}{
	"clock": timeInSeconds(),
}
