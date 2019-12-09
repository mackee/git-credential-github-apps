package main

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultAPIBaseURL    = "api.github.com"
	defaultCacheFilename = "git-credential-github-apps-token-cache"
)

func main() {
	var (
		privateKey     string
		appID          int64
		installationID int64
		login          string
		hostname       string
		apibase        string
		cachefile      string
	)
	cachedir, err := os.UserCacheDir()
	if err != nil {
		fmt.Printf("[ERROR] fail to detect cache dir: %s\n", err)
		os.Exit(1)
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

	flag.Parse()

	if privateKey == "" || appID == 0 {
		fmt.Println("[ERROR] must be set -privatekey and -appid")
		os.Exit(1)
	}
	if installationID == 0 && login == "" {
		fmt.Println("[ERROR] must be set -installationid or -login")
		os.Exit(1)
	}

	if flag.NArg() != 1 || flag.Arg(0) != "get" {
		os.Exit(0)
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
			fmt.Printf("[ERROR] %s\n", err)
			os.Exit(1)
		}
		options = append(options, WithStore(s))
	}

	input := &credentialInput{}
	if _, err := input.ReadFrom(os.Stdin); err != nil {
		fmt.Printf("[ERROR] %s\n", err)
		os.Exit(1)
	}
	if input.host != hostname || !strings.HasPrefix(input.protocol, "http") {
		fmt.Print(input)
		os.Exit(0)
	}

	if err := printCredential(ctx, privateKey, appID, options); err != nil {
		fmt.Printf("[ERROR] %s\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}

func printCredential(ctx context.Context, privateKey string, appID int64, options []AutherOption) error {
	auther, err := NewAutherFromFile(privateKey, appID, options...)
	if err != nil {
		return err
	}
	token, err := auther.FetchToken(ctx)

	fmt.Println("protocol=https")
	fmt.Println("host=github.com")
	fmt.Println("username=x-access-token")
	fmt.Printf("password=%s\n", token)

	return nil
}

type credentialInput struct {
	host     string
	protocol string
	username string
	password string
}

func (c *credentialInput) ReadFrom(r io.Reader) (int64, error) {
	input := bufio.NewScanner(r)
	for input.Scan() && input.Text() != "" {
		text := input.Text()
		kv := strings.SplitN(text, "=", 2)
		if len(kv) != 2 {
			return 0, fmt.Errorf("input text is invalid: input line=%s", text)
		}
		switch kv[0] {
		case "host":
			c.host = kv[1]
		case "protocol":
			c.protocol = kv[1]
		case "username":
			c.username = kv[1]
		case "password":
			c.password = kv[1]
		default:
			return 0, fmt.Errorf("input text is invalid: input line=%s", text)
		}
	}
	if err := input.Err(); err != nil {
		return 0, fmt.Errorf("fail to scan from reader: %w", err)
	}
	return 0, nil
}

func (c *credentialInput) String() string {
	out := &bytes.Buffer{}
	fmt.Fprintf(out, "host=%s\n", c.host)
	fmt.Fprintf(out, "protocol=%s\n", c.protocol)
	fmt.Fprintf(out, "username=%s\n", c.username)
	fmt.Fprintf(out, "password=%s\n", c.password)

	return out.String()
}
