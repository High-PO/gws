package main

import (
	"fmt"
	"os"

	"gws/internal/auth"
	"gws/internal/cli"
	"gws/internal/config"
	"gws/internal/credential"
	"gws/internal/shell"
)

var version = "dev"

func main() {
	cmd, err := cli.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "gws: %v\n", err)
		cli.PrintUsage(os.Stderr)
		os.Exit(1)
	}

	switch cmd.Type {
	case cli.CmdUsage:
		cli.PrintUsage(os.Stdout)
	case cli.CmdHelp:
		cli.PrintHelp(os.Stdout)
	case cli.CmdVersion:
		cli.PrintVersion(os.Stdout, version)
	case cli.CmdAuth:
		store := credential.NewStore()
		configMgr := &config.Manager{Store: store}
		authenticator := &auth.Authenticator{
			STS:    &auth.AWSSTSClient{},
			Config: configMgr,
			Shell:  &shell.Launcher{},
		}
		if err := authenticator.Authenticate(cmd.Profile, cmd.Token); err != nil {
			fmt.Fprintf(os.Stderr, "gws: %v\n", err)
			os.Exit(1)
		}
	}
}
