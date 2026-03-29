package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

// SessionCredential represents AWS STS temporary credentials.
type SessionCredential struct {
	AccessKeyID     string `json:"AccessKeyId"`
	SecretAccessKey string `json:"SecretAccessKey"`
	SessionToken    string `json:"SessionToken"`
	Expiration      string `json:"Expiration"`
}

// stsResponse wraps the STS API JSON response.
type stsResponse struct {
	Credentials SessionCredential `json:"Credentials"`
}

// STSClient abstracts AWS STS API calls.
type STSClient interface {
	GetSessionToken(profile, serialNumber, tokenCode string) (*SessionCredential, error)
}

// AWSSTSClient implements STSClient by calling the AWS CLI.
type AWSSTSClient struct{}

// GetSessionToken calls `aws sts get-session-token` and parses the JSON response.
func (c *AWSSTSClient) GetSessionToken(profile, serialNumber, tokenCode string) (*SessionCredential, error) {
	args := []string{"sts", "get-session-token", "--serial-number", serialNumber, "--token-code", tokenCode}
	if profile != "default" {
		args = append(args, "--profile", profile)
	}

	cmd := exec.Command("aws", args...)
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("STS API 호출 실패: %w", err)
	}

	var resp stsResponse
	if err := json.Unmarshal(output, &resp); err != nil {
		return nil, fmt.Errorf("AWS 응답 파싱 실패: %w", err)
	}

	return &resp.Credentials, nil
}

// ParseSTSResponse parses a raw STS JSON response into SessionCredential.
func ParseSTSResponse(data []byte) (*SessionCredential, error) {
	var resp stsResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("AWS 응답 파싱 실패: %w", err)
	}
	return &resp.Credentials, nil
}
