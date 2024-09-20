package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

type PlainFormatter struct{}

func (p *PlainFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var colorStart, colorEnd string
	var levelName string
	switch entry.Level {
	case logrus.TraceLevel:
		colorStart = "\033[37;2m"
		levelName = "TRACE"
	case logrus.DebugLevel:
		colorStart = "\033[37;2m"
		levelName = "DEBUG"
	case logrus.InfoLevel:
		colorStart = "\033[32m"
		levelName = "INFO"
	case logrus.WarnLevel:
		colorStart = "\033[33m"
		levelName = "WARN"
	case logrus.ErrorLevel:
		colorStart = "\033[31m"
		levelName = "ERROR"
	case logrus.FatalLevel, logrus.PanicLevel:
		colorStart = "\033[31m"
		levelName = "FATAL"
	}
	colorEnd = "\033[0m"

	return []byte(colorStart + levelName + ": " + entry.Message + colorEnd + "\n"), nil
}

func init() {
	logrus.SetFormatter(new(PlainFormatter))
	logrus.SetOutput(os.Stdout)
}
