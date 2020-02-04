package auth

import (
	"context"
	"strings"

	proto "github.com/WillCoates/FYP/common/protocol/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type key int

var metaToken key = 0
var metaPerms key = 1

var ErrMissingMetadata error = status.Errorf(codes.InvalidArgument, "Missing metadata")
var ErrBadToken error = status.Errorf(codes.Unauthenticated, "Bad token")

type AuthServerStream struct {
	token *Token
	perms []proto.Permission
	grpc.ServerStream
}

func (stream *AuthServerStream) Context() context.Context {
	ctx := stream.ServerStream.Context()
	ctx = context.WithValue(ctx, metaToken, stream.token)
	ctx = context.WithValue(ctx, metaPerms, stream.perms)
	return ctx
}

func NewAuthServerStream(token *Token, perms []proto.Permission, serverStream grpc.ServerStream) grpc.ServerStream {
	stream := new(AuthServerStream)
	stream.token = token
	stream.perms = perms
	stream.ServerStream = serverStream
	return stream
}

func UnaryServerInteceptor(next grpc.UnaryServerInterceptor, authClient *Client) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		meta, good := metadata.FromIncomingContext(ctx)

		if !good {
			return nil, ErrMissingMetadata
		}

		auth := meta.Get("Authorization")
		if len(auth) < 1 {
			return nil, ErrBadToken
		}

		// Most clients should perfix with Bearer, but if they don't assume header is only token
		rawToken := auth[0]
		rawToken = strings.TrimPrefix(rawToken, "Bearer ")

		token, perms, err := authClient.GetToken(rawToken)

		if err != nil {
			return nil, err
		}

		ctx = context.WithValue(ctx, metaToken, token)
		ctx = context.WithValue(ctx, metaPerms, perms)

		if next != nil {
			return next(ctx, req, info, handler)
		}
		return handler(ctx, req)
	}
}

func StreamServerInteceptor(next grpc.StreamServerInterceptor, authClient *Client) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		meta, good := metadata.FromIncomingContext(ss.Context())

		if !good {
			return ErrMissingMetadata
		}

		auth := meta.Get("Authorization")
		if len(auth) < 1 {
			return ErrBadToken
		}

		// Most clients should perfix with Bearer, but if they don't assume header is only token
		rawToken := auth[0]
		rawToken = strings.TrimPrefix(rawToken, "Bearer ")

		token, perms, err := authClient.GetToken(rawToken)

		if err != nil {
			return err
		}

		ss = NewAuthServerStream(token, perms, ss)

		if next != nil {
			return next(srv, ss, info, handler)
		}
		handler(srv, ss)
		return nil
	}
}

func FromContext(ctx context.Context) (*Token, []proto.Permission, bool) {
	token, ok := ctx.Value(metaToken).(*Token)
	if !ok {
		return nil, nil, false
	}

	perms, ok := ctx.Value(metaPerms).([]proto.Permission)
	if !ok {
		return nil, nil, false
	}

	return token, perms, true
}
