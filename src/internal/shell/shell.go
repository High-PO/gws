package shell

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"gws/internal/auth"
)

// Launcher launches or replaces shell sessions with AWS credentials.
type Launcher struct{}

// IsInsideGWSSession checks whether the current process is running inside
// a GWS-managed shell session by looking for the GWS_SESSION environment variable.
// Returns false if the variable is not set (safe fallback per requirement 6.2).
func IsInsideGWSSession() bool {
	_, ok := os.LookupEnv("GWS_SESSION")
	return ok
}

// buildEnv constructs a new environment variable slice by:
// 1. Removing existing AWS credential vars and GWS_SESSION from baseEnv
// 2. Appending the new credentials and GWS_SESSION=<profile>
// It is a pure function that does not read os.Environ() itself, making it testable.
func buildEnv(baseEnv []string, creds *auth.SessionCredential, profile string) []string {
	filtered := make([]string, 0, len(baseEnv))
	for _, e := range baseEnv {
		if strings.HasPrefix(e, "AWS_ACCESS_KEY_ID=") ||
			strings.HasPrefix(e, "AWS_SECRET_ACCESS_KEY=") ||
			strings.HasPrefix(e, "AWS_SESSION_TOKEN=") ||
			strings.HasPrefix(e, "GWS_SESSION=") {
			continue
		}
		filtered = append(filtered, e)
	}
	return append(filtered,
		"AWS_ACCESS_KEY_ID="+creds.AccessKeyID,
		"AWS_SECRET_ACCESS_KEY="+creds.SecretAccessKey,
		"AWS_SESSION_TOKEN="+creds.SessionToken,
		"GWS_SESSION="+profile,
	)
}

// getShellPath returns the appropriate shell path for the current OS.
// - macOS/Linux: uses $SHELL environment variable (falls back to /bin/sh)
// - Windows: uses cmd.exe
// The returned path is validated to be an absolute path (Unix) or a known shell (Windows).
func getShellPath() string {
	switch runtime.GOOS {
	case "windows":
		return "cmd.exe"
	default: // darwin, linux, etc.
		shellPath := os.Getenv("SHELL")
		if shellPath == "" || !filepath.IsAbs(shellPath) {
			return "/bin/sh"
		}
		return shellPath
	}
}

// launchNewShell creates a new child shell process with the given environment.
// It connects stdin, stdout, and stderr to the current process.
func launchNewShell(shellPath string, env []string) error {
	cmd := exec.Command(shellPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = env

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("셸 실행 실패: %w", err)
	}

	return nil
}

// Launch starts or replaces a shell session with the given credentials.
// It uses IsInsideGWSSession() to branch:
// - Outside GWS session: creates a new child shell via launchNewShell()
// - Inside GWS session: replaces the current shell via replaceShell()
// Both paths use buildEnv() to construct the environment with credentials and GWS_SESSION.
func (l *Launcher) Launch(creds *auth.SessionCredential, profile string) error {
	shellPath := getShellPath()
	env := buildEnv(os.Environ(), creds, profile)

	if IsInsideGWSSession() {
		fmt.Printf("Switching to profile [%s]...\n", profile)
		return replaceShell(shellPath, env)
	}

	fmt.Println("Launching new shell with AWS credentials...")
	return launchNewShell(shellPath, env)
}

// BuildCommand creates an exec.Cmd configured with AWS credentials as environment variables.
// This is exported for testing purposes.
func (l *Launcher) BuildCommand(creds *auth.SessionCredential, profile string) *exec.Cmd {
	shellPath := getShellPath()
	cmd := exec.Command(shellPath)
	cmd.Env = buildEnv(os.Environ(), creds, profile)
	return cmd
}
