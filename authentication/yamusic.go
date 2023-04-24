package authentication

import (
	"context"
	"github.com/ndrewnee/go-yamusic/yamusic"
	"os"
)

func YaMusic() *yamusic.Client {
	client := yaMusicAuth()

	return client
}

func yaMusicAuth() *yamusic.Client {
	ctx := context.Background()
	ya := yamusic.NewClient(yamusic.AccessToken(241139439, os.Getenv("YANDEX_SECRET")))
	status, _, err := ya.Account().GetStatus(ctx)

	if err != nil {
		panic(err)
	}

	ya.SetUserID(status.Result.Account.UID)

	return ya
}
