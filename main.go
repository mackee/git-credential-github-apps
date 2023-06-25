package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mackee/git-credential-github-apps/githubapps"
)

func main() {
	runner, err := githubapps.ParseArgs()
	if err != nil {
		if err == githubapps.ErrShowHelp {
			os.Exit(0)
		}
		fmt.Printf("[ERROR] %s\n", err)
		os.Exit(1)
	}

	args := runner.Args()
	if len(args) != 1 || args[0] != "get" {
		fmt.Printf("[ERROR] unexpected args\n")
		os.Exit(1)
	}

	input := &credentialInput{}
	if _, err := input.ReadFrom(os.Stdin); err != nil {
		fmt.Printf("[ERROR] %s\n", err)
		os.Exit(1)
	}
	if input.host != runner.Hostname() || !strings.HasPrefix(input.protocol, "http") {
		fmt.Print(input)
		os.Exit(0)
	}

	token, err := runner.Run(context.Background())
	if err != nil {
		fmt.Printf("[ERROR] %s\n", err)
		os.Exit(1)
	}

	printCredential(token)

	os.Exit(0)
}

func printCredential(token string) {
	fmt.Println("protocol=https")
	fmt.Println("host=github.com")
	fmt.Println("username=x-access-token")
	fmt.Printf("password=%s\n", token)
}

type credentialInput struct {
	host           string
	protocol       string
	username       string
	password       string
	wwwauthHeaders []string
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
		case "wwwauth[]":
			c.wwwauthHeaders = append(c.wwwauthHeaders, kv[1])
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
