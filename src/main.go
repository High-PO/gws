package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"syscall"
	"time"
)

type Config map[string]string

func main() {
	args := os.Args[1:]
	
	switch len(args) {
	case 0:
		showDetailedUsage()
		os.Exit(0)
	case 1:
		if args[0] == "help" {
			showDetailedUsage()
			os.Exit(0)
		} else if isSixDigitNumber(args[0]) {
			performMFAAuth("default", args[0])
		} else if args[0] == "--version" {
			fmt.Println("GWS VERSION 1.0")
			os.Exit(0)
		} else {
			// User provided a profile name without MFA token
			fmt.Printf("Profile '%s' selected. Now you need to provide a 6-digit MFA token.\n\n", args[0])
			fmt.Printf("Usage: gws %s <mfa-token>\n", args[0])
			fmt.Printf("Example: gws %s 123456\n\n", args[0])
			fmt.Println("The MFA token is the 6-digit code from your authenticator app.")
			os.Exit(1)
		}
	case 2:
		if args[0] == "help" {
			showCommandHelp(args[1])
			os.Exit(0)
		} else if isSixDigitNumber(args[1]) {
			performMFAAuth(args[0], args[1])
		} else {
			fmt.Printf("Error: Invalid MFA code '%s'. Must be exactly 6 digits.\n", args[1])
			fmt.Printf("Example: gws %s 123456\n", args[0])
			os.Exit(1)
		}
	default:
		fmt.Println("Error: Too many arguments provided.")
		fmt.Println()
		showBasicUsage()
		os.Exit(1)
	}
}

func isSixDigitNumber(s string) bool {
	match, _ := regexp.MatchString("^[0-9]{6}$", s)
	return match
}

func showDetailedUsage() {
	fmt.Println("GWS - Go + AWS CLI v2")
	fmt.Println("=====================")
	fmt.Println("A tool to easily use AWS CLI v2 as an MFA-enabled IAM user on macOS.")
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("  gws <mfa-token>                  # Use default AWS profile with MFA token")
	fmt.Println("  gws <profile> <mfa-token>        # Use specific AWS profile with MFA token")
	fmt.Println("  gws help                         # Show this help message")
	fmt.Println("  gws help <command>               # Show detailed help for a specific command")
	fmt.Println("  gws --version                    # Show version information")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("  1. Using default profile:")
	fmt.Println("     $ gws 123456")
	fmt.Println("     This will authenticate using your default AWS profile with MFA token 123456")
	fmt.Println()
	fmt.Println("  2. Using a specific profile:")
	fmt.Println("     $ gws production 654321")
	fmt.Println("     This will authenticate using the 'production' AWS profile with MFA token 654321")
	fmt.Println()
	fmt.Println("ARGUMENTS:")
	fmt.Println("  <profile>      AWS profile name configured in ~/.aws/config")
	fmt.Println("  <mfa-token>    6-digit code from your MFA device (e.g., Google Authenticator)")
	fmt.Println()
	fmt.Println("NOTES:")
	fmt.Println("  - MFA serial numbers are stored in ~/gws/config.json per profile")
	fmt.Println("  - On first use with a profile, you'll be prompted for the MFA serial number")
	fmt.Println("  - The tool creates a new shell session with temporary AWS credentials")
	fmt.Println("  - Credentials are valid for 12 hours by default (AWS STS limitation)")
	fmt.Println()
	fmt.Println("For more information on a specific command, use: gws help <command>")
}

func showBasicUsage() {
	fmt.Println("Usage: gws [profile] <mfa-token>")
	fmt.Println("       gws help")
	fmt.Println()
	fmt.Println("Run 'gws help' for detailed usage information.")
}

func showCommandHelp(command string) {
	// Check if the command is a profile name or a number
	if isSixDigitNumber(command) {
		fmt.Println("MFA Token Help")
		fmt.Println("==============")
		fmt.Println()
		fmt.Println("The MFA token is a 6-digit code from your authenticator app.")
		fmt.Println()
		fmt.Println("Format: Must be exactly 6 digits (000000-999999)")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  gws 123456              # Use with default profile")
		fmt.Println("  gws prod 987654         # Use with 'prod' profile")
		fmt.Println()
		fmt.Println("Where to find it:")
		fmt.Println("  - Google Authenticator")
		fmt.Println("  - Microsoft Authenticator") 
		fmt.Println("  - AWS Virtual MFA")
		fmt.Println("  - Hardware MFA device")
	} else {
		// Assume it's a profile name
		fmt.Printf("Profile '%s' Help\n", command)
		fmt.Println("==================")
		fmt.Println()
		fmt.Printf("To use the '%s' profile, you need to provide a 6-digit MFA token.\n", command)
		fmt.Println()
		fmt.Printf("Usage: gws %s <mfa-token>\n", command)
		fmt.Println()
		fmt.Printf("Example:\n")
		fmt.Printf("  gws %s 123456\n", command)
		fmt.Println()
		fmt.Println("Requirements:")
		fmt.Printf("  - Profile '%s' must be configured in ~/.aws/config\n", command)
		fmt.Println("  - You must have an MFA device associated with your IAM user")
		fmt.Println("  - The MFA serial number will be requested on first use")
		fmt.Println()
		fmt.Println("After successful authentication:")
		fmt.Println("  - A new shell session will be created")
		fmt.Println("  - AWS credentials will be set as environment variables")
		fmt.Println("  - Credentials are valid for 12 hours")
	}
}

func getConfigPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, "gws", "config.json")
}

func ensureConfigDir() {
	configPath := getConfigPath()
	dir := filepath.Dir(configPath)
	os.MkdirAll(dir, 0755)
	
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		config := make(Config)
		saveConfig(config)
	}
}

func loadConfig() Config {
	ensureConfigDir()
	
	data, err := os.ReadFile(getConfigPath())
	if err != nil {
		return make(Config)
	}
	
	var config Config
	json.Unmarshal(data, &config)
	return config
}

func saveConfig(config Config) {
	data, _ := json.MarshalIndent(config, "", "  ")
	os.WriteFile(getConfigPath(), data, 0644)
}

func getSerialNumber(profile string) string {
	config := loadConfig()
	return config[profile]
}

func saveSerialNumber(profile, serial string) {
	config := loadConfig()
	config[profile] = serial
	saveConfig(config)
}

func performMFAAuth(profile, mfaCode string) {
	serialNumber := getSerialNumber(profile)
	
	if serialNumber == "" {
		fmt.Print("Enter MFA serial number (arn:aws:iam::<account-id>:mfa/<device>): ")
		var input string
		if _, err := fmt.Scanln(&input); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading MFA serial number: %v\n", err)
			os.Exit(1)
		}
		serialNumber = input
		if serialNumber == "" {
			fmt.Fprintln(os.Stderr, "MFA serial number cannot be empty")
			os.Exit(1)
		}
		saveSerialNumber(profile, serialNumber)
	}
	
	// Build AWS command
	args := []string{"sts", "get-session-token", "--serial-number", serialNumber, "--token-code", mfaCode}
	if profile != "default" {
		args = append(args, "--profile", profile)
	}
	
	// Execute AWS command
	cmd := exec.Command("aws", args...)
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()
	
	if err != nil {
		os.Exit(1)
	}
	
	// Parse JSON output
	var result struct {
		Credentials struct {
			AccessKeyId     string
			SecretAccessKey string
			SessionToken    string
			Expiration      string
		}
	}
	
	if err := json.Unmarshal(output, &result); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing AWS response: %v\n", err)
		os.Exit(1)
	}
	
	// Create temporary script file
	tmpfile, err := os.CreateTemp("", "gws-*.sh")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating temp file: %v\n", err)
		os.Exit(1)
	}
	
	scriptContent := fmt.Sprintf(`#!/bin/bash
export AWS_ACCESS_KEY_ID='%s'
export AWS_SECRET_ACCESS_KEY='%s'
export AWS_SESSION_TOKEN='%s'
echo "AWS session credentials have been set successfully."
echo "Expires at: %s"
exec $SHELL -l
`, result.Credentials.AccessKeyId, result.Credentials.SecretAccessKey, result.Credentials.SessionToken, result.Credentials.Expiration)
	
	if _, err := tmpfile.Write([]byte(scriptContent)); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing script: %v\n", err)
		os.Exit(1)
	}
	tmpfile.Close()
	
	// Make script executable
	os.Chmod(tmpfile.Name(), 0755)
	
	// Execute the script in the current shell
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/bash"
	}
	
	// Clean up temp file after a delay (let the shell source it first)
	defer func() {
		time.Sleep(2 * time.Second)
		os.Remove(tmpfile.Name())
	}()
	
	// Replace the current process with a new shell that sources the script
	fmt.Printf("Launching new shell with AWS credentials...\n")
	syscall.Exec(shell, []string{shell, tmpfile.Name()}, os.Environ())
}