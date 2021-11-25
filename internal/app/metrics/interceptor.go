package metrics

import (
	"context"
	customErrors "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/errors"
	"google.golang.org/grpc"
	"strconv"
	"time"
)

type Interceptor struct {
	m *Metric
}

func NewInterceptor(m *Metric) *Interceptor {
	return &Interceptor{
		m: m,
	}
}

func (i *Interceptor) Collect(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	t := time.Now()

	resp, err := handler(ctx, req)
	result := "100"
	if err != nil {
		result = strconv.Itoa(customErrors.ErrorStringToCode(err.Error()))
		i.m.Errors.WithLabelValues(result, info.FullMethod, info.FullMethod).Inc()
	} else {
		i.m.Hits.WithLabelValues(result, info.FullMethod, info.FullMethod).Inc()
	}
	i.m.TotalHits.Inc()
	i.m.Durations.WithLabelValues(result, info.FullMethod, info.FullMethod).Observe(time.Since(t).Seconds())
	return resp, err
}
