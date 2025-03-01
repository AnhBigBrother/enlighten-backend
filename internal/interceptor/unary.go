package interceptor

import (
	"context"
	"log"

	"github.com/AnhBigBrother/enlighten-backend/cfg"
	jwtToken "github.com/AnhBigBrother/enlighten-backend/internal/pkg/jwt-token"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func AuthInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	access_token, ok := md["access_token"]
	if !ok {
		auth_token, ok := md["authorization"]
		if !ok {
			return handler(ctx, req)
		}
		access_token = auth_token
	}
	if len(access_token) == 0 {
		return handler(ctx, req)
	}

	userClaim, err := jwtToken.ParseAndValidate(access_token[0])
	if err != nil {
		return handler(ctx, req)
	}

	ctx_with_user := context.WithValue(ctx, cfg.CtxKeys.User, userClaim)
	return handler(ctx_with_user, req)
}

func LoggingInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	log.Println("unary interceptor:", info.FullMethod)
	return handler(ctx, req)
}
