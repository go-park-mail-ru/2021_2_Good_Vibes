package logger

import (
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/middleware/requestId"
	customLogger "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/logger"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"time"
)

const LoggerFieldName = "logger"

func AccessLog(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		start := time.Now()
		reqId_ := requestId.GetRequestIdFromContext(context)
		method := context.Request().Method
		remoteAddr := context.Request().RemoteAddr
		reqURI := context.Request().RequestURI
		customLogger.CustomLogger.LogAccessLog(reqId_, method, remoteAddr, "start request", reqURI)

		logger := customLogger.CustomLogger.LogrusLoggerHandler.WithFields(logrus.Fields{
			customLogger.RequestId: reqId_,
		})
		context.Set(LoggerFieldName, logger)

		defer func() {
			customLogger.CustomLogger.LogAccessLog(reqId_, method, remoteAddr, time.Since(start).String(), reqURI)
		}()

		return next(context)
	}
}
