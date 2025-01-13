package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/log"
	"golang.org/x/time/rate"
)

// -------------------------middlewares-------------------------------
// Go kit 试图通过使用中间件（或装饰器）模式来执行严格的关注分离（separation of concerns）。
// 中间件可以包装端点或服务以添加功能，比如日志记录、速率限制、负载平衡或分布式跟踪。

// 应用层中间件
// 如果要在应用程序层面添加日志，例如需要记录下详细的请求参数，那么就需要为我们的服务来定义中间件
type logMiddleware struct {
	logger log.Logger // go-kit自带的log
	// logger zap.Logger	// 集成zap库
	next AddService
}

// metrics记录
type instrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	countResult    metrics.Histogram
	next           AddService
}

func (mw logMiddleware) Sum(ctx context.Context, a, b int) (res int, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "sum",
			"a", a,
			"b", b,
			"output", res,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	res, err = mw.next.Sum(ctx, a, b)
	return
}

func (mw logMiddleware) Concat(ctx context.Context, a, b string) (res string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "sum",
			"a", a,
			"b", b,
			"output", res,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	res, err = mw.next.Concat(ctx, a, b)
	return
}

func (mw instrumentingMiddleware) Sum(ctx context.Context, a, b int) (res int, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "sum", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
		mw.countResult.Observe(float64(res))
	}(time.Now())

	res, err = mw.next.Sum(ctx, a, b)
	return
}

func (mw instrumentingMiddleware) Concat(ctx context.Context, a, b string) (res string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "concat", "error", "false"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	res, err = mw.next.Concat(ctx, a, b)
	return
}

// NewLogMiddleware 创建一个带日志的add service
func NewLogMiddleware(logger log.Logger, svc AddService) AddService {
	return &logMiddleware{
		logger: logger,
		next:   svc,
	}
}

// type Middleware func(Endpoint) Endpoint
// 在中间件接收Endpoint参数和返回Endpoint之间，你可以做任何事
func loggingMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			logger.Log("msg", "calling endpoint")
			defer logger.Log("msg", "called endpoint")
			return next(ctx, request)
		}
	}
}

// rateMiddleware 限流中间件
func rateMiddleware(limit *rate.Limiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			if !limit.Allow() {
				return nil, ErrRateLimit
			}
			return next(ctx, request)
		}
	}
}
