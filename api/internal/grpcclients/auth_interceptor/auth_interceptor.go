package authinterceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type AuthInterceptor struct{ AuthToken string }

func NewAuthInterceptor(authToken string) *AuthInterceptor {
	return &AuthInterceptor{AuthToken: authToken}
}

func (interceptor *AuthInterceptor) Unary() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctxWithToken := metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+interceptor.AuthToken)
		return invoker(ctxWithToken, method, req, reply, cc, opts...)
	}
}
