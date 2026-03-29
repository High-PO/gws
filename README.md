# gws

> go + aws cli v2 = gws

MFA가 활성화된 IAM 사용자를 위한 AWS CLI v2 인증 도구입니다. macOS, Linux, Windows를 지원합니다.

## 주요 기능

- MFA 토큰 기반 AWS STS 임시 자격 증명 발급
- 임시 자격 증명이 설정된 셸 세션 자동 실행
- GWS 셸 내부에서 프로필 전환 시 중첩 셸 없이 현재 셸 교체
- MFA 시리얼 번호를 OS 보안 저장소(macOS Keychain, Linux Secret Service, Windows Credential Manager)에 안전하게 저장
- 한국어/영어 도움말 지원 (시스템 로케일 자동 감지)

## 빌드

```bash
cd src
go build -ldflags "-X main.version=1.1.0" -o gws
```

## 설치

macOS / Linux:
```bash
sudo cp ./gws /usr/local/bin/
```

Windows:
```
gws.exe를 PATH에 포함된 디렉토리에 복사 (예: C:\Users\<username>\bin)
```

## 사용법

```bash
# 기본 프로필로 MFA 인증
gws <mfa_token>

# 특정 프로필로 MFA 인증
gws <profile> <mfa_token>

# 도움말
gws help

# 버전 정보
gws --version
```

## 셸 세션 전환

GWS가 생성한 셸 세션 내부에서 다른 프로필로 전환하면, 중첩 셸을 만들지 않고 현재 셸을 교체합니다.

```bash
$ gws dev 123456          # dev 프로필로 새 셸 세션 시작
$ gws production 654321   # 중첩 없이 production으로 전환
$ gws staging 111222      # 중첩 없이 staging으로 전환
```

`exit`을 반복할 필요 없이 프로필 간 자유롭게 이동할 수 있습니다.

## 테스트

```bash
cd src
go test ./...
```

## 프로젝트 구조

```
src/
├── main.go                    # 엔트리포인트
├── internal/
│   ├── auth/                  # STS 인증 및 셸 실행 오케스트레이션
│   ├── cli/                   # CLI 파싱 및 도움말 (한/영)
│   ├── config/                # MFA 시리얼 설정 관리
│   ├── credential/            # OS 보안 저장소 연동
│   └── shell/                 # 셸 세션 실행 및 교체
```
