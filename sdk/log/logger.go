package log

import (
	"fmt"
	"github.com/op/go-logging"
	"io"
	"os"
	"path/filepath"
)

var (
	ApiLogger *logging.Logger
	Logger    *logging.Logger
	logLevel  = map[string]logging.Level{
		"DEBUG":    logging.DEBUG,
		"INFO":     logging.INFO,
		"WARNING":  logging.WARNING,
		"ERROR":    logging.ERROR,
		"CRITICAL": logging.CRITICAL,
		"NOTICE":   logging.NOTICE,
	}
)

type LoggerConfig struct {
	ErrorLogPah   string
	AccessLogPath string
	Level         string
	DevMode       bool
	StdErr        bool
	NoRotate      bool
	RotateEnable  bool
}

func Init(cfg LoggerConfig) {
	ApiLogger = cfg.getLogger("access")
	Logger = cfg.getLogger("error")
}

func (cfg *LoggerConfig) getLevelBackend(logFile io.Writer, level string, format logging.Formatter) logging.LeveledBackend {
	logBackend := logging.NewLogBackend(logFile, "", 0)
	levelBackend := logging.AddModuleLevel(logging.NewBackendFormatter(logBackend, format))
	levelBackend.SetLevel(logLevel[level], "")
	return levelBackend
}

func (cfg *LoggerConfig) getLoggerBackend(logpath, level string, stderr bool) logging.LeveledBackend {
	loggerFmtOut := logging.MustStringFormatter(`%{color}%{time:2006-01-02 15:04:05} %{shortfile} %{shortfunc} [%{level:.4s} %{color:reset}] > %{message}`)
	loggerFmtFile := logging.MustStringFormatter(`%{color}%{time:2006-01-02 15:04:05} %{shortfile} %{shortfunc} [%{level:.4s} %{color:reset}] > %{message}`,
	)
	var err error
	var logWriter io.Writer
	if !cfg.RotateEnable {
		dir := filepath.Dir(logpath)
		err = os.MkdirAll(dir, 0755)
		if err == nil {
			var file *os.File
			file, err := os.OpenFile(logpath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
			if err == nil {
				if _, err = file.Seek(0, io.SeekEnd); err == nil {
					logWriter = file
				} else {
					fmt.Printf("seek file:%s to end err:%s\n", logpath, err)
				}
			} else {
				fmt.Printf("open file :%s err :%s\n", logpath, err)
			}
		} else {
			fmt.Printf("mkdifall err:%s\n", err)
		}
	}

	loggerBackend := logging.SetBackend()
	fileLevelBackend := cfg.getLevelBackend(logWriter, level, loggerFmtFile)
	stdErrLevelBackend := cfg.getLevelBackend(os.Stderr, level, loggerFmtOut)

	if err == nil && stderr {
		loggerBackend = logging.SetBackend(fileLevelBackend, stdErrLevelBackend)
	} else if err == nil && !stderr {
		loggerBackend = logging.SetBackend(fileLevelBackend)
	} else {
		loggerBackend = logging.SetBackend(stdErrLevelBackend)
	}

	return loggerBackend
}

func (cfg *LoggerConfig) getLogger(logger string) *logging.Logger {
	switch logger {
	case "access":
		var logger = logging.MustGetLogger("access")

		if cfg.AccessLogPath == "" {
			cfg.AccessLogPath = "/var/log/messages"
		}
		loggerBackend := cfg.getLoggerBackend(cfg.AccessLogPath, cfg.Level, cfg.StdErr || cfg.DevMode)
		logger.SetBackend(loggerBackend)
		return logger
	default:
		var logger = logging.MustGetLogger("error")
		if cfg.AccessLogPath == "" {
			cfg.AccessLogPath = "/var/log/messages"
		}
		loggerBackend := cfg.getLoggerBackend(cfg.AccessLogPath, cfg.Level, cfg.StdErr || cfg.DevMode)
		logger.SetBackend(loggerBackend)
		return logger
	}
	return nil

}
