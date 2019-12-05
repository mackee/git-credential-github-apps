package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/v28/github"
)

func main() {
	var (
		privateKey     string
		appID          int64
		installationID int64
		login          string
	)
	flag.StringVar(&privateKey, "privatekey", "private_key.pem", "private key of GitHub Apps")
	flag.Int64Var(&appID, "appid", 0, "App ID of GitHub Apps")
	flag.Int64Var(&installationID, "installationid", 0, "Installation ID of organization or user on GitHub Apps")
	flag.StringVar(&login, "login", "", "login name of organization or user. if not set -installationid, search from Installation ID by use value.")

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

	tr := http.DefaultTransport
	atr, err := ghinstallation.NewAppsTransportKeyFromFile(tr, appID, privateKey)
	if err != nil {
		fmt.Printf("[ERROR] cannot create apps transport: %s\n", err)
		os.Exit(1)
	}
	client := github.NewClient(&http.Client{Transport: atr})

	if installationID == 0 {
		lo := &github.ListOptions{
			PerPage: 100,
			Page:    1,
		}
	OUTER:
		for {
			ins, resp, err := client.Apps.ListInstallations(ctx, lo)
			if err != nil {
				fmt.Printf("[ERROR] fail to fetch installations: %s\n", err)
				os.Exit(1)
			}
			for _, in := range ins {
				if login == in.GetAccount().GetLogin() {
					installationID = in.GetID()
					break OUTER
				}
			}
			if resp.LastPage == 0 || lo.Page == resp.LastPage {
				fmt.Printf("[ERROR] %s is not found in installations\n", login)
				os.Exit(1)
			}
			lo.Page++
		}
	}

	token, _, err := client.Apps.CreateInstallationToken(ctx, installationID, nil)
	if err != nil {
		fmt.Printf("[ERROR] fail to create installation token: %s", err)
		os.Exit(1)
	}

	fmt.Println("protocol=https")
	fmt.Println("host=github.com")
	fmt.Println("username=x-access-token")
	fmt.Printf("password=%s\n", token.GetToken())

	os.Exit(0)
}
