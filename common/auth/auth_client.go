package auth

import (
	"context"
	"errors"
	"io"

	proto "github.com/WillCoates/FYP/common/protocol/auth"
	"google.golang.org/grpc"
)

type Client struct {
	clientConn *grpc.ClientConn
	authClient proto.AuthServiceClient
	keys       *KeyBundle
}

var ErrTokenExpired = errors.New("Token expired")

func NewAuthClient(authURL string, keys string) (*Client, error) {
	var err error
	client := new(Client)

	client.clientConn, err = grpc.Dial(authURL)
	if err != nil {
		return nil, err
	}

	client.authClient = proto.NewAuthServiceClient(client.clientConn)

	client.keys, err = LoadBundleHTTP(keys, MasterKey())

	if err != nil {
		return nil, err
	}

	return client, nil
}

func (client *Client) GetToken(rawToken string) (*Token, []proto.Permission, error) {
	token, err := ParseToken([]byte(rawToken), client.keys.Keys)
	if err != nil {
		return nil, nil, err
	}

	var protoToken proto.Token
	protoToken.Token = rawToken

	resp, err := client.authClient.IsTokenValid(context.Background(), &protoToken)

	if err != nil {
		return nil, nil, err
	} else if !resp.Valid {
		return nil, nil, ErrTokenExpired
	}

	permsClient, err := client.authClient.GetTokenPermissions(context.Background(), &protoToken)
	if err != nil {
		return nil, nil, err
	}

	var perms []proto.Permission

	for {
		perm, err := permsClient.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, err
		}
		perms = append(perms, *perm)
	}

	return token, perms, nil
}
