package main

import (
	"context"
	"errors"
)

// ---------------------------Services---------------------------------
// Go kit 服务层应该努力遵守整洁架构或六边形架构
// 服务（指Go kit中的service层）是实现所有业务逻辑的地方。服务层通常将多个端点粘合在一起。
// 在 Go kit 中，服务层通常被抽象为接口，这些接口的实现包含业务逻辑。
// AddService 把两个东西加在一起

type AddService interface {
	Sum(ctx context.Context, a, b int) (int, error)
	Concat(ctx context.Context, a, b string) (string, error)
}

type addService struct{}

var (
	// ErrEmptyString 两个参数都是空字符串的错误
	ErrEmptyString = errors.New("两个参数都是空字符串")
	ErrRateLimit   = errors.New("request rate limit")
)

// Sum 返回两个数的和
func (addService) Sum(_ context.Context, a, b int) (int, error) {
	// 业务逻辑
	return a + b, nil
}

// Concat 拼接两个字符串
func (addService) Concat(_ context.Context, a, b string) (string, error) {
	if a == "" && b == "" {
		return "", ErrEmptyString
	}
	return a + b, nil
}

// NewService 创建一个add service
func NewService() AddService {
	return &addService{}
}
