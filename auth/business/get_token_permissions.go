package business

import (
	"context"

	"github.com/WillCoates/FYP/common/auth"
)

func (logic *Logic) GetTokenPermissions(ctx context.Context, token *auth.Token) (chan string, error) {
	audience, err := logic.GetAudience(ctx, token.Payload.Audience)
	if err != nil {
		return nil, err
	}

	channel := make(chan string)

	go func() {
		for _, perm := range audience.Perms {
			channel <- perm
		}
		close(channel)
	}()

	return channel, nil
}
