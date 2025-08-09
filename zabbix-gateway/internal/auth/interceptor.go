package auth

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	validToken string
}

func NewAuthInterceptor(token string) *AuthInterceptor {
	return &AuthInterceptor{validToken: token}
}

func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}
		authHeaders := md.Get("authorization")
		if len(authHeaders) < 1 {
			return nil, status.Error(codes.Unauthenticated, "missing authorization header")
		}
		if authHeaders[0] != "Bearer "+interceptor.validToken {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}
		return handler(ctx, req)
	}
}
