package requestId

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/logger"
	"github.com/labstack/echo/v4"
	"net/http"
)

func RequestId(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		reqId := make([]byte, 8)
		_, err := rand.Read(reqId)
		if err != nil {
			logger.CustomLogger.LogErrorInfo(logger.BadRequestId, err.Error())
			return context.NoContent(http.StatusInternalServerError)
		}

		base64ID := base64.URLEncoding.EncodeToString(reqId)
		context.Set(logger.RequestId, base64ID)
		return next(context)
	}
}

func GetRequestIdFromContext(ctx echo.Context) string {
	intReqId := ctx.Get(logger.RequestId)
	if intReqId == nil {
		logger.CustomLogger.LogErrorInfo(logger.BadRequestId, "error get ReqID")
		return logger.BadRequestId
	}

	return fmt.Sprintf("%v", intReqId)
}
