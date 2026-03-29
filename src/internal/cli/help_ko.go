package cli

import (
	"fmt"
	"io"
)

func printHelpKo(w io.Writer) {
	fmt.Fprintln(w, "GWS - Go + AWS CLI v2")
	fmt.Fprintln(w, "=====================")
	fmt.Fprintln(w, "MFA가 활성화된 IAM 사용자를 위한 AWS CLI v2 인증 도구입니다.")
	fmt.Fprintln(w, "macOS, Linux, Windows를 지원합니다.")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "사용법:")
	fmt.Fprintln(w, "  gws <mfa-토큰>                   # 기본 프로필로 MFA 인증")
	fmt.Fprintln(w, "  gws <프로필> <mfa-토큰>          # 지정 프로필로 MFA 인증")
	fmt.Fprintln(w, "  gws help                         # 도움말 표시")
	fmt.Fprintln(w, "  gws --version                    # 버전 정보 표시")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "예시:")
	fmt.Fprintln(w, "  1. 기본 프로필 사용:")
	fmt.Fprintln(w, "     $ gws 123456")
	fmt.Fprintln(w, "     기본 AWS 프로필로 MFA 토큰 123456을 사용하여 인증합니다")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "  2. 특정 프로필 사용:")
	fmt.Fprintln(w, "     $ gws production 654321")
	fmt.Fprintln(w, "     'production' 프로필로 MFA 토큰 654321을 사용하여 인증합니다")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "인수:")
	fmt.Fprintln(w, "  <프로필>       ~/.aws/config에 설정된 AWS CLI 프로필 이름")
	fmt.Fprintln(w, "  <mfa-토큰>    MFA 디바이스의 6자리 코드 (예: Google Authenticator)")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "참고:")
	fmt.Fprintln(w, "  - MFA 시리얼 번호는 OS 보안 저장소에 프로필별로 안전하게 저장됩니다")
	fmt.Fprintln(w, "  - 프로필 최초 사용 시 MFA 시리얼 번호 입력을 요청합니다")
	fmt.Fprintln(w, "  - 인증 성공 시 임시 AWS 자격 증명이 설정된 새 셸 세션이 실행됩니다")
	fmt.Fprintln(w, "  - 자격 증명은 기본 12시간 유효합니다 (AWS STS 제한)")
}

func printUsageKo(w io.Writer) {
	fmt.Fprintln(w, "사용법: gws [프로필] <mfa-토큰>")
	fmt.Fprintln(w, "       gws help")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "자세한 사용법은 'gws help'를 실행하세요.")
}
