package requestId

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/logger"
	"github.com/labstack/echo/v4"
	"net/http"
)

const (
	ReqId        = "req_id"
	BadRequestId = "-1"
)

func RequestId(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		reqId := make([]byte, 8)
		_, err := rand.Read(reqId)
		if err != nil {
			logger.CustomLogger.LogErrorInfo(BadRequestId, err.Error())
			return context.NoContent(http.StatusInternalServerError)
		}

		base64ID := base64.URLEncoding.EncodeToString(reqId)
		context.Set(ReqId, base64ID)
		return next(context)
	}
}

func GetRequestIdFromContext(ctx echo.Context) string {
	intReqId := ctx.Get(ReqId)
	if intReqId == nil {
		logger.CustomLogger.LogErrorInfo(BadRequestId, "error get ReqID")
		return BadRequestId
	}

	return fmt.Sprintf("%v", intReqId)
}
