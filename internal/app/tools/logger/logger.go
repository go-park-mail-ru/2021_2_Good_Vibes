package logger

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/t-tomalak/logrus-easy-formatter"
	"os"
)

type Logger struct {
	LogrusLoggerAccess  *logrus.Logger
	LogrusLoggerHandler *logrus.Logger
}

var CustomLogger Logger

const (
	RequestId     = "request_id"
	Method        = "method"
	WorkTime      = "work_time"
	RemoteAddress = "remote_address"
	BadRequestId  = "-1"
)

func InitLogger() {
	logrus.SetLevel(logrus.TraceLevel)

	CustomLogger.LogrusLoggerAccess = &logrus.Logger{
		Out:   os.Stderr,
		Level: logrus.DebugLevel,
		Formatter: &easy.Formatter{
			TimestampFormat: "2006-01-02 15:04:05",
			LogFormat: "[%lvl%] [%request_id%] %time%:" +
				" %method% | %msg% | remote_address - %remote_address% | work_time - %work_time%\n",
		},
	}

	CustomLogger.LogrusLoggerHandler = &logrus.Logger{
		Out:   os.Stderr,
		Level: logrus.TraceLevel,
		Formatter: &easy.Formatter{
			TimestampFormat: "2006-01-02 15:04:05",
			LogFormat: "[%lvl%] [%request_id%] %time%: " +
				"%msg%\n",
		},
	}
}

func (l *Logger) LogErrorInfo(requestId_ string, err string) {
	l.LogrusLoggerHandler.WithFields(logrus.Fields{
		RequestId: requestId_,
	}).Error(err)
}

func (l *Logger) LogAccessLog(requestId_ string, method string, remoteAddress string,
	workTime string, requestURI string) {
	CustomLogger.LogrusLoggerAccess.WithFields(logrus.Fields{
		RequestId:     requestId_,
		Method:        method,
		RemoteAddress: remoteAddress,
		WorkTime:      workTime,
	}).Info(requestURI)
}

func TryGetLoggerFromContext(ctx echo.Context) *logrus.Entry {
	logger, ok := ctx.Get("logger").(*logrus.Entry)
	if !ok {
		logger = logrus.WithField("req_id", "2299")
	}
	return logger
}
