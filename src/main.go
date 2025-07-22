package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

type Config map[string]string

func main() {
	args := os.Args[1:]
	
	switch len(args) {
	case 0:
		fmt.Println("Error: Please provide MFA code")
		os.Exit(1)
	case 1:
		if isSixDigitNumber(args[0]) {
			performMFAAuth("default", args[0])
		} else {
			fmt.Println("Error: Please provide MFA code")
			os.Exit(1)
		}
	case 2:
		if isSixDigitNumber(args[1]) {
			performMFAAuth(args[0], args[1])
		} else {
			fmt.Println("Error: Invalid MFA code. Must be 6 digits.")
			os.Exit(1)
		}
	default:
		fmt.Println("Usage: gws [profile] <mfa-code>")
		fmt.Println("       gws <mfa-code>        # uses default profile")
		fmt.Println("       gws <profile> <mfa-code>")
		os.Exit(1)
	}
}

func isSixDigitNumber(s string) bool {
	match, _ := regexp.MatchString("^[0-9]{6}$", s)
	return match
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
		fmt.Scanln(&serialNumber)
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
export AWS_ACCESS_KEY_ID=%s
export AWS_SECRET_ACCESS_KEY=%s
export AWS_SESSION_TOKEN=%s
echo "AWS session credentials have been set successfully."
exec $SHELL
`, result.Credentials.AccessKeyId, result.Credentials.SecretAccessKey, result.Credentials.SessionToken)
	
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
	
	cmd = exec.Command(shell, tmpfile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error executing shell: %v\n", err)
		os.Exit(1)
	}
	
	// Clean up
	os.Remove(tmpfile.Name())
}