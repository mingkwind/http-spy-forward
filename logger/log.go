package logger

import (
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var (
	Log *logrus.Entry
)

func init() {
	logger := logrus.New()
	logger.Formatter = new(prefixed.TextFormatter)
	logger.Level = logrus.InfoLevel
	// 添加前缀
	Log = logger.WithFields(logrus.Fields{"prefix": "http-spy-forward"})
}
