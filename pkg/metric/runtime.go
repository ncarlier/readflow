package metric

import (
	"expvar"
	"runtime"
	"time"
)

var startTime = time.Now().UTC()

func goroutines() interface{} {
	return runtime.NumGoroutine()
}

func uptime() interface{} {
	uptime := time.Since(startTime)
	return int64(uptime)
}

func init() {
	expvar.Publish("goroutines", expvar.Func(goroutines))
	expvar.Publish("uptime", expvar.Func(uptime))
}
