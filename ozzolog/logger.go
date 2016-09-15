package ozzolog

import (
	"github.com/gamexg/ozzo-log"
)

var logger = log.NewLogger()

type Config struct {
	LogConsoleLevel int

	LogFileName        string
	LogFileBackupCount int
	LogFileLevel       int
	LogEmailHost       string
	LogEmailUsername   string
	LogEmailPassword   string
	LogEmailSender     string
	LogEmailSubject    string
	LogEmailRecipients []string
	LogEmailLevel      int
}

// 非多线程安全
func Open(conf *Config) {
	ct := log.NewConsoleTarget()
	ct.MaxLevel = log.Level(conf.LogConsoleLevel)
	logger.Targets = append(logger.Targets, ct)

	if conf.LogFileName != "" {
		ft := log.NewFileTarget()
		ft.FileName = conf.LogFileName
		ft.MaxLevel = log.Level(conf.LogFileLevel)
		ft.MaxBytes = 1 * 1024 * 1024 //2M
		ft.BackupCount = conf.LogFileBackupCount
		ft.Rotate = true //自动切割旋转日志文件
		logger.Targets = append(logger.Targets, ft)
	}

	if conf.LogEmailHost != "" {
		et := log.NewMailTarget()
		et.Host = conf.LogEmailHost
		et.Password = conf.LogEmailPassword
		et.Username = conf.LogEmailUsername
		et.Subject = conf.LogEmailSubject
		et.Recipients = conf.LogEmailRecipients
		et.Sender = conf.LogEmailSender
		et.MaxLevel = log.Level(conf.LogEmailLevel)

		logger.Targets = append(logger.Targets, et)
	}

	logger.Open()

	logger.Debug("日志完成初始化...")
}

func GetLogger(category string, formatter ...log.Formatter) *log.Logger {
	return logger.GetLogger(category, formatter...)
}

func Close() {
	logger.Close()
}
