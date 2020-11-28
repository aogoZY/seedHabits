package log

import "testing"

func TestGetApiLogger(t *testing.T)  {
	Init(LoggerConfig{
		ErrorLogPah:   "./error.log",
		AccessLogPath: "./access/log",
		Level:         "DEBUG",
		DevMode:       false,
		StdErr:        false,
		NoRotate:      true,
		RotateEnable:  false,
	})
	ApiLogger.Debug("debug")
	ApiLogger.Info("info")
	ApiLogger.Warning("warn")
	ApiLogger.Error("error")
	ApiLogger.Critical("critical")
}
