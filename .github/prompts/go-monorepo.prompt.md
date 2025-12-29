---
agent: agent
---

# 역할/전제
- 당신은 Go 기반 데스크톱 앱 개발을 주도하는 시니어 개발자다.
- Windows 지원이 필수이며, 프레임워크는 Wails를 사용한다.
- 우리 팀은 “여러 개의 앱”을 하나의 저장소에서 함께 개발한다.
- 상황에 따라 여러 앱을 동시에 실행하거나, 특정 앱 하나만 실행할 수 있어야 한다.
- 목표는 Monorepo 스타일(공유 코드/공유 설정/공유 빌드 파이프라인)로 개발 생산성과 일관성을 확보하는 것이다.

# 최종 목표(Outcome)
Wails 기반 멀티 앱 Monorepo 개발 환경을 설계/스캐폴딩하고, VSCode에서 바로 개발/실행/빌드/테스트가 가능한 상태로 만든다.

# 구체 요구사항(Functional Requirements)
1) 멀티 앱 구조
- apps/ 하위에 앱들이 독립적으로 존재 (예: apps/app-a, apps/app-b …)
- 각 앱은 Wails 프로젝트로서 “backend(Go)” + “frontend(TS: vite framework)”를 가진다.
- 공통 Go 코드(도메인/유틸/utils/클라이언트)는 packages/ 또는 packages/ 공유 영역으로 분리한다.
- 공통 프론트 코드(UI 컴포넌트/유틸/테마)도 공유 패키지로 분리할 수 있어야 한다(선택이 아니라 “가능”해야 함).

2) 단일 앱 실행 / 다중 앱 동시 실행
- 특정 앱 하나만 실행: 단일 커맨드로 가능해야 함
- 여러 앱 동시에 실행: 단일 커맨드로 가능해야 함(프로세스 동시 기동)
- 실행 타겟은 개발 모드(wails dev) 기준이며, 필요 시 포트 충돌 방지 가이드를 포함

3) 빌드/패키징(Windows 우선)
- 각 앱을 Windows 실행 파일(exe)로 빌드하는 표준 커맨드 제공
- 산출물 경로를 통일하고, 앱별 결과물이 서로 덮어쓰지 않도록 분리
- 아이콘/메타데이터/버전(가능하면) 앱별로 관리 가능해야 함

4) 공유 코드 관리(Go)
- Go 모듈/워크스페이스(go.work) 기반으로 앱 간 공유 패키지를 매끄럽게 참조
- 공유 패키지 변경 시 앱에서 즉시 반영되는 로컬 개발 흐름 제공

5) 개발 도구/품질
- gofmt/goimports, golangci-lint, go test 표준화
- 프론트 lint/test(사용 프레임워크에 맞춰) 표준화
- VSCode에서 실행 가능한 tasks.json, launch.json(가능하면), 추천 확장(.vscode/extensions.json) 제공

6) 부트스트랩/온보딩
- 신규 개발자가 “한 번의 명령”으로 개발 준비를 끝낼 수 있어야 함
- 예: bootstrap(툴 설치/의존성 설치/환경 점검) 커맨드 제공

# 제약조건(Constraints)
- Wails 사용은 필수(Windows 지원 포함)
- Monorepo로 운영하되, 앱별 독립성(의존성/빌드/실행)이 훼손되면 안 됨
- 전역 설치 의존을 최소화(가능하면 go install / package manager로 버전 고정)
- OS 타겟은 Windows 우선, 향후 macOS/Linux 확장 가능한 설계를 선호
- 팀 개발 환경은 VSCode 중심(Tasks/Debug 흐름 제공)

# 작업 범위(Scope)
아래 산출물을 “실제로 생성”하거나, 현재 저장소가 있다면 “기존 구조를 존중하면서” 최소 변경으로 반영한다.

1) 권장 디렉터리 구조 제안(예시)
- apps/
  - app-a/
  - app-b/
- packages/ (또는 packages/)
  - go/...
  - ui/...(선택)
- scripts/ (빌드/릴리즈/유틸)
- .vscode/ (tasks, settings, extensions)
- go.work / go.work.sum
- 루트 task runner 설정(Makefile 또는 Taskfile 또는 Mage 중 1개 선택)

2) 표준 커맨드 정의(예시; 실제로 동작해야 함)
- bootstrap: 개발 준비
- dev:<app>: 특정 앱 개발 실행
- dev:all: 여러 앱 동시 실행
- test: Go/프론트 테스트
- lint: Go/프론트 린트
- build:<app>:win: Windows 빌드
- build:all:win: 전체 앱 Windows 빌드

3) VSCode 통합
- .vscode/tasks.json에 위 커맨드 매핑
- .vscode/settings.json에 Go/프론트 개발 경험 최적화(포맷, 린트, 터미널, 파일 제외 등)
- .vscode/extensions.json에 필수 확장 추천

4) CI(선택이지만 가능하면 포함)
- GitHub Actions 기준
- PR 시 lint/test, main merge 시 Windows 빌드
- 캐시 적용(Go build cache, node/pnpm 캐시 등)

# 성공 기준(Success Criteria)
- 저장소 루트에서 단일 명령으로 특정 앱이 실행된다(dev:<app>)
- 저장소 루트에서 단일 명령으로 2개 이상의 앱이 동시에 실행된다(dev:all)
- 저장소 루트에서 단일 명령으로 Windows 빌드 산출물이 생성된다(build:<app>:win)
- go.work 기반 공유 패키지 수정이 앱에서 즉시 반영된다(로컬 개발 확인)
- VSCode에서 Task 실행만으로 개발/빌드/테스트 흐름이 재현된다
- 신규 앱 추가 절차가 문서화되어 10분 내로 스캐폴딩 가능하다

# 진행 방식(Implementation Steps)
1) Monorepo 구조/툴 체인(Go workspaces + 프론트 워크스페이스) 선택 및 이유 명시
2) 디렉터리/모듈 설계 후 실제 스캐폴딩(또는 기존 코드 반영)
3) 루트 실행 커맨드/스크립트 구현(단일/다중 실행 포함)
4) Windows 빌드 파이프라인 구현
5) VSCode tasks/settings/extensions 구성
6) README에 “실행/빌드/새 앱 추가”를 절차 중심으로 문서화
7) 간단 검증 체크리스트 제공(명령어 결과 포함)

# 출력/산출물 형식
- 최종 디렉터리 트리
- 생성/수정된 주요 파일 목록과 핵심 내용(필요 시 코드 포함)
- 실행 가능한 명령어 목록과 사용 예시
- (가능하면) CI 워크플로 파일
