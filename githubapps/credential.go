package githubapps

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

const (
	defaultAPIBaseURL    = "api.github.com"
	defaultCacheFilename = "git-credential-github-apps-token-cache"
)

var (
	ErrShowHelp = errors.New("receive show help args")
)

// Runner provides retrieving token
type Runner struct {
	hostname   string
	privateKey string
	appID      int64
	args       []string
	options    []AutherOption
}

// Hostname returns hostname from arguments
func (r *Runner) Hostname() string {
	return r.hostname
}

// Args returns command options not contains parsed.
func (r *Runner) Args() []string {
	return r.args
}

// Run returns credentials
func (r *Runner) Run(ctx context.Context) (string, error) {
	return readCredential(ctx, r.privateKey, r.appID, r.options)
}

// ParseArgs is retrieve github apps credentials with command line options.
func ParseArgs() (*Runner, error) {
	var (
		privateKey     string
		appID          int64
		installationID int64
		login          string
		hostname       string
		apibase        string
		cachefile      string
		showHelp       bool
	)
	cachedir, err := os.UserCacheDir()
	if err != nil {
		return nil, fmt.Errorf("fail to detect cache dir: %s", err)
	}

	flag.StringVar(&privateKey, "privatekey", "private_key.pem", "private key of GitHub Apps")
	flag.Int64Var(&appID, "appid", 0, "App ID of GitHub Apps")
	flag.Int64Var(&installationID, "installationid", 0, "Installation ID of organization or user on GitHub Apps")
	flag.StringVar(&login, "login", "", "login name of organization or user. if not set -installationid, search from Installation ID by use value.")
	flag.StringVar(&hostname, "hostname", "github.com", "hostname as using for an accessing in git")
	flag.StringVar(&apibase, "apibase", "api.github.com", "API hostname as using for a fetching GitHub APIs")
	flag.StringVar(
		&cachefile, "cachefile", filepath.Join(cachedir, defaultCacheFilename),
		"filename as save cached token",
	)
	flag.BoolVar(&showHelp, "h", false, "show this help")

	flag.Parse()

	if showHelp {
		flag.PrintDefaults()
		return nil, ErrShowHelp
	}

	if privateKey == "" || appID == 0 {
		return nil, fmt.Errorf("must be set -privatekey and -appid")
	}
	if installationID == 0 && login == "" {
		return nil, fmt.Errorf("must be set -installationid or -login")
	}

	ctx := context.Background()

	var options []AutherOption
	if defaultAPIBaseURL != apibase {
		options = append(options, WithBaseURL(apibase))
	}
	if installationID == 0 {
		options = append(options, WithLogin(ctx, login))
	} else {
		options = append(options, WithInstallationID(installationID))
	}
	if cachefile != "" {
		s, err := NewFileStore(cachefile)
		if err != nil {
			return nil, fmt.Errorf("%s", err)
		}
		options = append(options, WithStore(s))
	}

	return &Runner{hostname: hostname, options: options, privateKey: privateKey, appID: appID, args: flag.Args()}, nil
}

func readCredential(ctx context.Context, privateKey string, appID int64, options []AutherOption) (string, error) {
	auther, err := NewAutherFromFile(privateKey, appID, options...)
	if err != nil {
		return "", err
	}

	return auther.FetchToken(ctx)
}
