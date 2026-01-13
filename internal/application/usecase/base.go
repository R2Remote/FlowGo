package usecase

import (
	"context"
)

// UseCase 用例接口
type UseCase[Req, Resp any] interface {
	Execute(ctx context.Context, req Req) (Resp, error)
}

// BaseUseCase 基础用例
type BaseUseCase struct {
	// 可以在这里注入依赖
}
