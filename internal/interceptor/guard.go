package interceptor

import (
	"context"
	"strings"

	"github.com/AnhBigBrother/enlighten-backend/cfg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryAuthGuard(services []string, methods []string) grpc.UnaryServerInterceptor {
	serviceMap := map[string]bool{}
	for _, s := range services {
		serviceMap[s] = true
	}
	methodMap := map[string]bool{}
	for _, r := range methods {
		methodMap[r] = true
	}

	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		arr := strings.Split(info.FullMethod, "/")
		if serviceMap[arr[1]] || methodMap[info.FullMethod] {
			_, ok := ctx.Value(cfg.CtxKeys.User).(map[string]interface{})
			if !ok {
				return nil, status.Error(codes.Unauthenticated, "access_token failed")
			}
		}

		return handler(ctx, req)
	}
}
