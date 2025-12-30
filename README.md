
# Monorepo for Wails + Vite(React/TS) Apps

## 구조
- apps/app-file-controller: Wails + Vite(React/TS) 앱
- apps/app-decg-link: Wails + Vite(React/TS) 앱
- packages/domain, utils, ui-antd: Go/프론트 공통 모듈/유틸/컴포넌트
- scripts: 빌드/유틸 스크립트
- .vscode: VSCode 협업용 에디터 설정(포맷, 린트, 확장 등)
- go.work: Go workspace(모듈/패키지 통합 관리)
- Taskfile.yml: cross-platform task runner(공통 명령어)

## 주요 명령어(Taskfile)
- `task bootstrap`: 개발 환경 준비(툴/의존성 설치)
- `task dev:app-file-controller`: app-file-controller 개발 실행
- `task dev:app-decg-link`: app-decg-link 개발 실행
- `task dev:all`: 두 앱 동시 실행(프론트/백엔드 모두)
- `task build:app-file-controller:win`: app-file-controller Windows 빌드
- `task build:app-decg-link:win`: app-decg-link Windows 빌드
- `task build:all:win`: 전체 앱 Windows 빌드
- `task lint:go`: Go 린트
- `task test:go`: Go 테스트
- `npm run format`: 프론트/공통 코드 prettier 포맷
- `npm run lint`: 프론트/공통 코드 eslint 린트
- `npm run go:fmt`: 전체 Go 코드 gofmt 포맷
- `npm run go:imports`: 전체 Go 코드 goimports 포맷

## 신규 앱 추가
1. `apps/` 하위에 Wails CLI로 새 앱 생성
2. go.work에 use 추가
3. Taskfile.yml에 dev/build task 추가
4. 필요시 .vscode/settings.json, package.json 등 공통 설정/포맷/린트 적용

## 부트스트랩
- `task bootstrap` 한 번으로 개발 준비 완료(Go/Node 의존성, 툴 자동 설치)


## VSCode 협업 환경
- .vscode/settings.json: 저장 시 자동 포맷/린트(Go: goimports, JS/TS: prettier, ESLint)
- .vscode/extensions.json: 추천 확장(Go, Prettier, ESLint 등)
- 모든 팀원이 동일한 개발 환경을 유지할 수 있도록 루트 설정을 사용

## 빌드 후 앱 실행 방법

1. 앱별 빌드
	 - `task build:app-file-controller:win` 또는 `task build:app-decg-link:win` 실행
	 - macOS/Linux는 각 앱 디렉터리에서 `wails build` 사용 가능

2. 빌드 산출물 실행
	 - 각 앱의 `build/bin/` 디렉터리에서 실행 파일(.app, .exe 등) 확인
	 - 예시 (macOS):
		 ```bash
		 ./apps/app-file-controller/build/bin/app-file-controller.app/Contents/MacOS/app-file-controller &
		 ./apps/app-decg-link/build/bin/app-decg-link.app/Contents/MacOS/app-decg-link &
		 ```
	 - 예시 (Windows):
		 ```cmd
		 .\apps\app-file-controller\build\bin\app-file-controller.exe
		 .\apps\app-decg-link\build\bin\app-decg-link.exe
		 ```

3. 두 앱을 동시에 실행하려면 위 명령을 각각 터미널에서 실행하거나, 백그라운드(&)로 실행

---
자세한 내용은 각 디렉터리/설정 파일 참고
