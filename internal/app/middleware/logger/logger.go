package logger

import (
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/middleware/requestId"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/logger"
	"github.com/labstack/echo/v4"
	"time"
)

func AccessLog(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		start := time.Now()
		reqId_ := requestId.GetRequestIdFromContext(context)
		method := context.Request().Method
		remoteAddr := context.Request().RemoteAddr
		reqURI := context.Request().RequestURI

		logger.CustomLogger.LogAccessLog(reqId_, method, remoteAddr, "start request", reqURI)

		defer func() {
			logger.CustomLogger.LogAccessLog(reqId_, method, remoteAddr, time.Since(start).String(), reqURI)
		}()

		return next(context)
	}
}
