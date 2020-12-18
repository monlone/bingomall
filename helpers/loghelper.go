package helper

import (
	"github.com/sirupsen/logrus"
	"os"
)

var (
	SQLLogger     *logrus.Logger
	AccessLogger  *logrus.Logger
	ServiceLogger *logrus.Logger
	WorkLogger    *logrus.Logger
	ErrorLogger   *logrus.Logger
)

func Logger(outPath string) *logrus.Logger {
	logger := logrus.New()
	_, err := os.Stat(outPath)
	if os.IsNotExist(err) {
		// 文件不存在,创建
		_, err = os.Create(outPath)
	}
	file, err := os.OpenFile(outPath, os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		logger.Out = file
	} else {
		logger.Info("打开 " + outPath + " 下的日志文件失败, 使用默认方式显示日志！")
	}
	return logger
}

func init() {
	WorkLogger = Logger("logs/work.log")
	AccessLogger = Logger("logs/access.log")
	ErrorLogger = Logger("logs/error.log")
	SQLLogger = Logger("logs/sql.log")
	ServiceLogger = Logger("logs/service.log")
}
