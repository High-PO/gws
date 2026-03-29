package cli

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// isKorean checks if the system locale is Korean.
// 1. LC_ALL, LANG 환경 변수 확인
// 2. macOS인 경우 defaults read -g AppleLanguages fallback
func isKorean() bool {
	for _, key := range []string{"LC_ALL", "LANG"} {
		val := strings.ToLower(os.Getenv(key))
		if strings.HasPrefix(val, "ko") {
			return true
		}
	}

	switch runtime.GOOS {
	case "darwin":
		out, err := exec.Command("defaults", "read", "-g", "AppleLanguages").Output()
		if err == nil && strings.Contains(strings.ToLower(string(out)), "ko") {
			return true
		}
	case "windows":
		out, err := exec.Command("powershell", "-NoProfile", "-Command",
			"(Get-Culture).Name").Output()
		if err == nil && strings.HasPrefix(strings.ToLower(strings.TrimSpace(string(out))), "ko") {
			return true
		}
	}

	return false
}

// PrintHelp prints the full help/usage information to the given writer.
func PrintHelp(w io.Writer) {
	if isKorean() {
		printHelpKo(w)
	} else {
		printHelpEn(w)
	}
}

// PrintUsage prints a brief usage summary to the given writer.
func PrintUsage(w io.Writer) {
	if isKorean() {
		printUsageKo(w)
	} else {
		printUsageEn(w)
	}
}

// PrintVersion prints the version string to the given writer.
func PrintVersion(w io.Writer, version string) {
	fmt.Fprintf(w, "gws version %s\n", version)
}
