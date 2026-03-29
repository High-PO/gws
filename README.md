# gws

> go + aws cli v2 = gws

A CLI tool for AWS MFA authentication. Supports macOS, Linux, and Windows.

MFA가 활성화된 IAM 사용자를 위한 AWS CLI v2 인증 도구입니다. macOS, Linux, Windows를 지원합니다.

---

## Features / 주요 기능

- Authenticate with AWS STS using MFA tokens and launch a shell with temporary credentials
- Switch profiles inside a GWS shell session without nesting shells
- MFA serial numbers stored securely in OS keyring (macOS Keychain, Linux Secret Service, Windows Credential Manager)
- Bilingual help: auto-detects system locale (English / Korean)

---

- MFA 토큰 기반 AWS STS 임시 자격 증명 발급 및 셸 세션 자동 실행
- GWS 셸 내부에서 프로필 전환 시 중첩 셸 없이 현재 셸 교체
- MFA 시리얼 번호를 OS 보안 저장소에 안전하게 저장
- 한국어/영어 도움말 지원 (시스템 로케일 자동 감지)

---

## Prerequisites / 사전 요구사항

- Go 1.24+
- AWS CLI v2 installed and configured (`~/.aws/config`)
- At least one IAM profile with MFA enabled

---

- Go 1.24+
- AWS CLI v2 설치 및 설정 (`~/.aws/config`)
- MFA가 활성화된 IAM 프로필 1개 이상

---

## Build / 빌드

```bash
cd src
go build -ldflags "-X main.version=1.1.0" -o gws
```

## Install / 설치

macOS / Linux:
```bash
sudo cp ./gws /usr/local/bin/
```

Windows:
```
Copy gws.exe to a directory in your PATH (e.g. C:\Users\<username>\bin)
```

---

## Usage / 사용법

### Basic Authentication / 기본 인증

```bash
# Authenticate with default profile / 기본 프로필로 인증
gws <mfa_token>

# Authenticate with a specific profile / 특정 프로필로 인증
gws <profile> <mfa_token>
```

On first use with a profile, GWS will prompt for the MFA serial number (ARN format). It is then stored securely in the OS keyring and reused automatically.

프로필 최초 사용 시 MFA 시리얼 번호(ARN 형식)를 입력받습니다. 이후 OS 보안 저장소에 저장되어 자동으로 재사용됩니다.

### Examples / 예시

```bash
# First run: prompts for MFA serial, then launches a shell with temporary credentials
# 첫 실행: MFA 시리얼 입력 후 임시 자격 증명이 설정된 셸 실행
$ gws 123456
Enter MFA serial number (arn:aws:iam::<account-id>:mfa/<device>): arn:aws:iam::000000000000:mfa/mydevice
AWS session credentials set. Expires at: 2024-01-01T12:00:00Z
Launching new shell with AWS credentials...

# Subsequent runs: no prompt, just authenticate
# 이후 실행: 프롬프트 없이 바로 인증
$ gws 654321
AWS session credentials set. Expires at: 2024-01-02T00:00:00Z
Launching new shell with AWS credentials...

# Use a named profile / 특정 프로필 사용
$ gws production 123456
```

### Shell Session Switching / 셸 세션 전환

When you run `gws` inside an existing GWS shell session, it replaces the current shell instead of creating a nested one. No more stacking `exit` commands.

GWS 셸 세션 내부에서 `gws`를 실행하면 중첩 셸을 만들지 않고 현재 셸을 교체합니다. `exit`을 반복할 필요가 없습니다.

```bash
$ gws dev 123456              # Start a new GWS shell for 'dev'
                               # 'dev' 프로필로 새 GWS 셸 시작
$ gws production 654321       # Switch to 'production' (no nesting)
                               # 'production'으로 전환 (중첩 없음)
Switching to profile [production]...
$ gws staging 111222          # Switch to 'staging' (no nesting)
                               # 'staging'으로 전환 (중첩 없음)
Switching to profile [staging]...
```

### Other Commands / 기타 명령어

```bash
# Show help / 도움말 표시
gws help

# Show version / 버전 정보
gws --version
```

---

## How It Works / 동작 원리

1. GWS calls `aws sts get-session-token` with your MFA token to obtain temporary credentials.
2. A new shell process is launched with `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_SESSION_TOKEN`, and `GWS_SESSION` set as environment variables.
3. If already inside a GWS session (`GWS_SESSION` is set), the current shell is replaced via `exec` (Unix) or process spawn + exit (Windows) instead of creating a nested shell.

---

1. GWS는 MFA 토큰으로 `aws sts get-session-token`을 호출하여 임시 자격 증명을 발급받습니다.
2. `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_SESSION_TOKEN`, `GWS_SESSION` 환경 변수가 설정된 새 셸 프로세스를 실행합니다.
3. 이미 GWS 세션 내부(`GWS_SESSION` 설정됨)라면, 중첩 셸 대신 `exec`(Unix) 또는 프로세스 생성+종료(Windows) 방식으로 현재 셸을 교체합니다.

---

## Platform Support / 플랫폼 지원 현황

| Platform | Status | Target Version |
|----------|--------|----------------|
| macOS | Tested | v1.1.0 (current) |
| Windows | Planned | v1.2.0 |
| Linux | Planned | v1.3.0 |

Currently tested on macOS only. Windows support is planned for the next release, and Linux for the release after that.

현재 macOS에서만 테스트되었습니다. Windows는 다음 버전(v1.2.0), Linux는 그 다음 버전(v1.3.0)에서 지원 예정입니다.

---

## Testing / 테스트

```bash
cd src
go test ./...
```

## Project Structure / 프로젝트 구조

```
src/
├── main.go                    # Entry point / 엔트리포인트
├── internal/
│   ├── auth/                  # STS authentication & shell orchestration
│   ├── cli/                   # CLI parsing & help (en/ko)
│   ├── config/                # MFA serial config management
│   ├── credential/            # OS keyring integration
│   └── shell/                 # Shell session launch & replacement
```
