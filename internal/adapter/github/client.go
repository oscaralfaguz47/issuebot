package github

import (
	"context"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation/v2"
)

type Client struct {
	appID          int64
	privateKeyPath string
}

func NewClient(appID int64, privateKeyPath string) *Client {
	return &Client{appID: appID, privateKeyPath: privateKeyPath}
}

func (c *Client) InstallationToken(ctx context.Context, installationID int64) (string, error) {
	itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, c.appID, installationID, c.privateKeyPath)
	if err != nil {
		return "", err
	}
	return itr.Token(ctx)
}
