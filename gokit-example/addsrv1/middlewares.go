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

type withTrimMiddleware struct {
	next AddService
	trimService endpoint.Endpoint	// 通过它调用其他的服务
}

func NewServiceWithTrim(trimEndpoint endpoint.Endpoint, svc AddService) AddService {
	return &withTrimMiddleware{
		trimService: trimEndpoint,
		next:        svc,
	}
}

// 为 withTrimMiddleware 实现 AddService 接口
func (tm withTrimMiddleware) Sum(ctx context.Context, a, b int) (res int, err error) {
	return tm.next.Sum(ctx, a, b) // 复用之前的逻辑
}

func (tm withTrimMiddleware) Concat(ctx context.Context, a, b string) (res string, err error) {
	// 需要新的逻辑处理
	// 外部调用我们的Concat方法时
	// 1. 发起RPC调用 trim_service 对数据进行处理 （调用其他服务/依赖其他的服务）
	respA, err := tm.trimService(ctx, trimRequest{s: a}) // 执行，其实是作为客户端对外发起请求
	if err != nil {
		return "", err
	}
	respB, err := tm.trimService(ctx, trimRequest{s: b}) // 执行，其实是作为客户端对外发起请求
	if err != nil {
		return "", err
	}
	trimA := respA.(trimResponse) // 拿到处理后的响应
	trimB := respB.(trimResponse) // 拿到处理后的响应

	// 2. 拿到处理后的数据再拼接
	return tm.next.Concat(ctx, trimA.s, trimB.s)
}