package cli

import (
	"fmt"
	"regexp"
)

// CommandType represents the type of CLI command parsed from arguments.
type CommandType int

const (
	CmdHelp CommandType = iota
	CmdVersion
	CmdAuth
	CmdUsage
)

// Command represents a parsed CLI command.
type Command struct {
	Type    CommandType
	Profile string
	Token   string
}

var sixDigitRe = regexp.MustCompile(`^[0-9]{6}$`)

// profileRe allows alphanumeric, hyphens, underscores, and dots (safe for AWS CLI --profile).
var profileRe = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

// isSixDigitNumber checks if a string is exactly 6 digits.
func isSixDigitNumber(s string) bool {
	return sixDigitRe.MatchString(s)
}

// Parse parses CLI arguments and returns a Command.
// args should be os.Args[1:] (without the program name).
func Parse(args []string) (Command, error) {
	switch len(args) {
	case 0:
		return Command{Type: CmdUsage}, nil
	case 1:
		switch {
		case args[0] == "help":
			return Command{Type: CmdHelp}, nil
		case args[0] == "--version":
			return Command{Type: CmdVersion}, nil
		default:
			if err := ValidateMFAToken(args[0]); err != nil {
				return Command{}, fmt.Errorf("invalid argument: %s", args[0])
			}
			return Command{Type: CmdAuth, Profile: "default", Token: args[0]}, nil
		}
	case 2:
		if !profileRe.MatchString(args[0]) {
			return Command{}, fmt.Errorf("invalid profile name: %q (alphanumeric, hyphens, underscores, dots only)", args[0])
		}
		if err := ValidateMFAToken(args[1]); err != nil {
			return Command{}, err
		}
		return Command{Type: CmdAuth, Profile: args[0], Token: args[1]}, nil
	default:
		return Command{}, fmt.Errorf("too many arguments")
	}
}
