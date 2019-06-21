package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

var AppLog = logrus.New()
var SysLog = logrus.New()

func InitLogs() {
	cfg := GetConfig()
	if cfg.Environment == "production" {
		appFile, err := os.OpenFile(cfg.Logs.FilePath+"app.log", os.O_CREATE|os.O_WRONLY, 0666)
		if err == nil {
			AppLog.Out = appFile
		} else {
			SysLog.Info("初始化日志失败,日志输出到控制台")
		}

		SysFile, err := os.OpenFile(cfg.Logs.FilePath+"sys.log", os.O_CREATE|os.O_WRONLY, 0666)
		if err == nil {
			SysLog.Out = SysFile
		} else {
			SysLog.Info("初始化日志失败,日志输出到控制台")
		}

	} else {
		SysLog.Out = os.Stdout

	}

	AppLog.SetLevel(logrus.Level(cfg.Logs.Level))
	SysLog.SetLevel(logrus.Level(cfg.Logs.Level))

	switch cfg.Logs.Formatter {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{})
	case "text":
		logrus.SetFormatter(&logrus.TextFormatter{})
	default:
		SysLog.Error("不支持的Formatter格式,只支持json和text")
	}

}
