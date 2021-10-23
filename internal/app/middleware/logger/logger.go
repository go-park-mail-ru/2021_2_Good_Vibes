package logger

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"time"
)

func AccessLog(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error{
		start := time.Now()

		defer func() {
			logrus.WithFields(logrus.Fields{
				"method":         context.Request().Method,
				"remote_address": context.Request().RemoteAddr,
				"work_time":      time.Since(start),
			}).Info(context.Request().RequestURI)
		}()

		return next(context)
	}
}
