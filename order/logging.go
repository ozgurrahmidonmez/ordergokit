package order

import (
	"github.com/go-kit/kit/log"
	"time"
)

type logmw struct {
	logger       log.Logger
	orderService OrderService
}

func LoggingMiddleware(logger log.Logger) ServiceMiddleware{
	return func(s OrderService) OrderService {
		return logmw{logger,s}
	}
}

func (mw logmw) Add(o OrderRequest) (err error) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "uppercase",
			"input", o,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = mw.orderService.Add(o)
	return
}


