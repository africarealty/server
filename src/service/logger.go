package service

import (
	"github.com/africarealty/server/src/kit/log"
)

var Logger = log.Init(&log.Config{Level: log.TraceLevel, Format: log.FormatterJson})

func LF() log.CLoggerFunc {
	return func() log.CLogger {
		return log.L(Logger).Srv("africarealty").Node("africarealty")
	}
}

func L() log.CLogger {
	return LF()()
}
