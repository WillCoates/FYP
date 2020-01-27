package main

import (
	"context"

	"github.com/WillCoates/FYP/auth/business"
	proto "github.com/WillCoates/FYP/common/protocol/auth"
)

type AuthService struct {
	logic *business.Logic
}

func NewAuthService(logic *business.Logic) *AuthService {
	service := new(AuthService)
	service.logic = logic
	return service
}

func (service *AuthService) Authenticate(ctx context.Context, request *proto.AuthenticateRequest) (*proto.AuthenticateResponse, error) {
	token, err := service.logic.Authenticate(ctx, request.EmailAddress, request.Password, request.Audience, request.Duration)
	if err != nil {
		return nil, err
	}

	protoToken := new(proto.Token)

	protoToken.Token, err = service.logic.EncodeTokenStr(token)

	if err != nil {
		return nil, err
	}

	result := new(proto.AuthenticateResponse)
	result.Success = true
	result.Token = protoToken

	return result, nil
}

func (service *AuthService) GetTokenPermissions(protoToken *proto.Token, stream proto.AuthService_GetTokenPermissionsServer) error {
	token, err := service.logic.DecodeTokenStr(protoToken.Token)
	if err != nil {
		return err
	}

	channel, err := service.logic.GetTokenPermissions(stream.Context(), token)
	if err != nil {
		return err
	}

	for perm := range channel {
		protoPerm := new(proto.Permission)
		protoPerm.Permission = perm
		protoPerm.For = token.Payload.Subject
		stream.Send(protoPerm)
	}

	return nil
}

func (service *AuthService) IsTokenValid(ctx context.Context, protoToken *proto.Token) (*proto.TokenValidResponse, error) {
	token, err := service.logic.DecodeTokenStr(protoToken.Token)
	if err != nil {
		return nil, err
	}

	valid := true

	err = service.logic.IsTokenValid(ctx, token)
	if err == business.ErrTokenInvalid {
		valid = false
	} else if err != nil {
		return nil, err
	}

	response := new(proto.TokenValidResponse)
	response.Valid = valid
	return response, nil
}

func (service *AuthService) InvalidateToken(ctx context.Context, protoToken *proto.Token) (*proto.InvalidateTokenResponse, error) {
	token, err := service.logic.DecodeTokenStr(protoToken.Token)

	if err != nil {
		return nil, err
	}

	service.logic.InvalidateToken(ctx, token)

	result := new(proto.InvalidateTokenResponse)
	result.Success = true
	return result, nil
}

func (service *AuthService) Register(ctx context.Context, request *proto.RegistrationRequest) (*proto.RegistrationResponse, error) {
	err := service.logic.Register(ctx, request.EmailAddress, request.Name, request.Password)
	if err != nil {
		return nil, err
	}

	result := new(proto.RegistrationResponse)
	result.Success = true

	return result, err
}

var _ proto.AuthServiceServer = (*AuthService)(nil)
