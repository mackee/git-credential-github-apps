package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"syscall"

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
	if len(args) == 0 || (len(args) == 1 && args[0] == "--") {
		fmt.Println("[ERROR] not provides command from args. eg. ghawrap -- yourcli options...")
		os.Exit(1)
	}
	name := args[0]
	cmdArgs := args[1:]
	if name == "--" {
		name = args[1]
		cmdArgs = args[2:]
	}

	ctx := context.Background()
	token, err := runner.Run(ctx)
	if err != nil {
		fmt.Printf("[ERROR] %s\n", err)
		os.Exit(1)
	}

	err = runCommand(name, cmdArgs, []string{"GITHUB_TOKEN=" + token})
	if err != nil {
		fmt.Printf("[ERROR] %s\n", err)
		os.Exit(1)
	}
}

func runCommand(name string, args, envVars []string) error {
	bin, err := exec.LookPath(name)
	if err != nil {
		return err
	}

	return syscall.Exec(bin, args, envVars)
}
