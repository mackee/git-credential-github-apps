package githubapps

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v68/github"
)

// Auther provides API Token with authentication on GitHub Apps.
type Auther interface {
	FetchToken(context.Context) (string, error)
}

type auther struct {
	atr            *ghinstallation.AppsTransport
	installationID int64
	store          Store
}

// NewAutherFromFile is constructor with filename that private key for Auther.
func NewAutherFromFile(privateKey string, appID int64, options ...AutherOption) (Auther, error) {
	tr := http.DefaultTransport
	atr, err := ghinstallation.NewAppsTransportKeyFromFile(tr, appID, privateKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create apps transport: %w", err)
	}

	a := &auther{atr: atr}
	for _, o := range options {
		err := o(a)
		if err != nil {
			return nil, fmt.Errorf("fail to initialize by option: %w", err)
		}
	}

	return a, nil
}

func (a *auther) FetchToken(ctx context.Context) (string, error) {
	if a.store != nil && !a.store.Expired() {
		return a.store.Token(), nil
	}

	t, err := a.fetchToken(ctx)
	if err != nil {
		return "", err
	}

	if a.store != nil {
		err := a.store.Save(t.GetToken(), t.GetExpiresAt().Time)
		if err != nil {
			return "", err
		}
	}

	return t.GetToken(), nil
}

func (a *auther) fetchToken(ctx context.Context) (*github.InstallationToken, error) {
	client := github.NewClient(&http.Client{Transport: a.atr})

	token, _, err := client.Apps.CreateInstallationToken(ctx, a.installationID, nil)
	if err != nil {
		return nil, fmt.Errorf("fail to create installation token: %w", err)
	}
	return token, nil
}

// AutherOption is set parameter. This is like a Functional Option Pattern.
type AutherOption func(*auther) error

// WithLogin is optional parameter for Auther.
// This provides to search a installation ID from login name.
func WithLogin(ctx context.Context, login string) AutherOption {
	return func(a *auther) error {
		client := github.NewClient(&http.Client{Transport: a.atr})
		lo := &github.ListOptions{
			PerPage: 100,
			Page:    1,
		}
	OUTER:
		for {
			ins, resp, err := client.Apps.ListInstallations(ctx, lo)
			if err != nil {
				return fmt.Errorf("fail to fetch installations: %w", err)
			}
			for _, in := range ins {
				if login == in.GetAccount().GetLogin() {
					a.installationID = in.GetID()
					break OUTER
				}
			}
			if resp.LastPage == 0 || lo.Page == resp.LastPage {
				return fmt.Errorf("%s is not found in installations", login)
			}
			lo.Page++
		}
		return nil
	}
}

// WithInstallationID is option parameter with installation ID for Auther.
func WithInstallationID(id int64) AutherOption {
	return func(a *auther) error {
		a.installationID = id
		return nil
	}
}

// WithBaseURL provides change api endpoint. For GitHub Enterprise
func WithBaseURL(base string) AutherOption {
	return func(a *auther) error {
		a.atr.BaseURL = base
		return nil
	}
}

// WithStore provides injection token store to Auther.
func WithStore(store Store) AutherOption {
	return func(a *auther) error {
		a.store = store
		return nil
	}
}
