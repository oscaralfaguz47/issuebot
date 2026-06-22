package github

import (
	"context"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation/v2"
	gh "github.com/google/go-github/v79/github"
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

func (c *Client) clientForInstallation(installationID int64) (*gh.Client, error) {
	itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, c.appID, installationID, c.privateKeyPath)
	if err != nil {
		return nil, err
	}
	return gh.NewClient(&http.Client{Transport: itr}), nil
}

func (c *Client) PostComment(ctx context.Context, installationID int64, owner, repo string, issueNumber int, body string) error {
	client, err := c.clientForInstallation(installationID)
	if err != nil {
		return err
	}

	comment := &gh.IssueComment{Body: gh.String(body)}
	_, _, err = client.Issues.CreateComment(ctx, owner, repo, issueNumber, comment)
	return err
}
