# Monorepo for Wails + Vite(React/TS) Apps

## 구조
- apps/app-file-controller: Wails + Vite(React/TS) 앱
- apps/app-decg-link: Wails + Vite(React/TS) 앱
- packages/domain, utils, ui: Go/프론트 공유 코드
- scripts: 빌드/유틸 스크립트
- .vscode: VSCode 설정
- go.work: Go workspace
- Taskfile.yml: cross-platform task runner

## 주요 명령어(Taskfile)
- `task bootstrap`: 개발 환경 준비(툴/의존성 설치)
- `task dev:app-file-controller`: app-file-controller 개발 실행
- `task dev:app-decg-link`: app-decg-link 개발 실행
- `task dev:all`: 두 앱 동시 실행
- `task build:app-file-controller:win`: app-file-controller Windows 빌드
- `task build:app-decg-link:win`: app-decg-link Windows 빌드
- `task build:all:win`: 전체 앱 Windows 빌드
- `task lint:go`: Go 린트
- `task test:go`: Go 테스트

## 신규 앱 추가
1. `apps/` 하위에 Wails CLI로 새 앱 생성
2. go.work에 use 추가
3. Taskfile.yml에 dev/build task 추가

## 부트스트랩
- `task bootstrap` 한 번으로 개발 준비 완료

## VSCode
- .vscode/tasks.json, settings.json, extensions.json 제공 예정

---
자세한 내용은 각 디렉터리/설정 파일 참고
