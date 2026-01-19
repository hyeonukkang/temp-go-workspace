
[Copilot Agent Prompt] Wails 멀티앱 Monorepo에 신규 앱 "folder-watcher" 추가 및 구현

## 역할
너는 Go + Wails 기반 Windows 데스크톱 앱을 주도하는 시니어 개발자다. 현재 저장소는 “Wails 멀티 앱 Monorepo” 구조이며, 기존 규칙/툴체인/Task runner/CI/VSCode 설정을 최우선으로 존중한다. 최소 변경 원칙을 지킨다.

## 최종 목표(Outcome)
- apps/folder-watcher 하위에 Wails 앱(backend Go + frontend Vite TS)을 생성하고,
- 폴더 워칭 설정/수정/삭제/목록, 개별 Start/Stop, 변경 알림(인앱 + Windows 토스트), output(move/copy), 로그(UI + 로컬 영구 저장)을 구현한다.
- 루트 단일 커맨드로 dev/build/test/lint가 가능하도록 monorepo 통합을 완료한다.
- VSCode Tasks로 dev/build/test/lint 흐름이 재현된다.

---

## 확정 요구사항(Functional Requirements)
1) 워칭 범위: 재귀 아님(해당 폴더 1단계만)
2) watcher 항목: sourcePath는 항목당 1개, 항목은 여러 개 등록 가능
3) 목록 조회/수정/삭제 가능
4) 변경 감지 시 알림: 인앱 토스트/배너 + Windows 시스템 토스트(성공/실패 모두)
5) watcher 항목별 옵션:
   - outputPath(옵션)
   - action: move | copy (outputPath 없으면 none)
6) output 처리 트리거: “추가된 파일”만
   - Create 이벤트
   - 드롭/이동 유입도 포함: Rename 유입(폴더로 들어온 파일)을 Create로 간주
7) 충돌 정책: output에 동일 파일명 있으면 이름 변경(넘버링)하여 회피
8) Start는 항목별 수동(Start 버튼)로만 동작(자동 시작 금지)
9) 모든 이벤트/처리 기록은 로그로 남긴다:
   - UI에서 조회
   - 로컬에 영구 저장
10) 실패 정책:
   - 실패하면 원본 파일은 삭제하지 않는다
   - 실패는 로그/상태로 남겨 추후 확인(자동 재시도 폭주 방지)

---

## “누락 0(보장)”을 위한 핵심 설계(필수)
- fsnotify 이벤트만 믿지 말고, “내구성 작업 큐 + 주기 리스캔”으로 결과 누락을 방지한다.
- watcher별로 SQLite 기반 job queue를 사용한다(앱 종료/재시작에도 작업 유실 방지).
- 리스캔 주기: 10초 (watcher running 상태일 때만)
- 리스캔은 sourcePath 1단계만 스캔하며, “새로 추가된 파일”만 job으로 enqueue한다.

### 작업 큐(Job Queue) 상태
- queued | running | done | failed
- failed는 자동 enqueue 대상에서 제외(리스캔 폭주 방지)
- (선택/가능하면) UI에 “Retry(단일 항목)” 제공: failed -> queued로 전환

---

## 임시파일 무시(필수)
기본 ignore 패턴을 이벤트 수집 + 리스캔 모두에 동일 적용:
- *.tmp, *.swp, *~, .DS_Store, Thumbs.db, ~$*
추후 watcher별 커스텀은 가능하게 설계만 해두고(인터페이스/필드), UI는 지금은 필수 아님.

---

## Start 시 기존 파일 처리 옵션(필수: 옵션으로 제공)
watcher 항목에 OnStartBehavior 옵션 추가:
- skip-existing (Start 시점에 이미 있던 파일은 처리하지 않음)
- process-existing (Start 시점에 이미 있던 파일도 처리)
기본값은 안전 우선으로 skip-existing.
skip-existing일 때는 “seen set(스냅샷)”을 만들고, 리스캔에서 현재 폴더에 없는 파일은 seen에서 제거하여 재유입 시 처리되게 한다.

---

## 안정성/안전장치(필수)
- outputPath가 sourcePath와 동일하거나 sourcePath 내부(하위)면 금지(루프/자기복사 방지)
- 디렉터리는 move/copy 대상에서 제외(파일만)
- 파일 안정성 체크 후 처리(쓰기 중 파일 조기 이동 방지):
  - 크기/mtime이 연속 N회 동일할 때까지 대기(최대 타임아웃)
- Rename 유입 처리:
  - Rename 이벤트 수신 시 150ms 정도 지연 후 파일 존재하면 “유입”으로 간주하여 create job enqueue

---

## Monorepo 통합(반드시 수행)
1) 저장소 스캔
- 기존 Task runner 종류(Taskfile/Makefile/Mage 등), node 패키지매니저(pnpm/yarn/npm), 기존 apps 구조, Wails 버전, CI 유무 확인
- 기존 패턴 그대로 folder-watcher에 적용

2) 스캐폴딩
- apps/folder-watcher에 Wails 프로젝트 생성(기존 앱들과 동일 버전/템플릿/구성)
- packages/go 공통 코드 영역 재사용 가능하도록 설계(필요 시만 최소 추가)
- go.work에 apps/folder-watcher 및 packages/go 모듈 포함

3) 루트 커맨드(기존 방식에 추가)
- dev:folder-watcher
- build:folder-watcher:win
- lint / test / dev:all / build:all:win 흐름에 folder-watcher 포함(이미 존재하면 확장만)

4) VSCode
- .vscode/tasks.json에 dev/build/test/lint 매핑 추가
- .vscode/settings.json과 extensions.json은 기존 정책 유지 + 필요한 최소만 추가

5) CI(존재할 경우)
- PR: lint + test
- main merge: Windows build
- 캐시(go, node) 적용(기존 워크플로 스타일 준수)

---

## 구현 작업 분해(Implementation Steps)
Step 1. 저장소 분석 결과를 먼저 요약(발견한 Task runner/패키지매니저/Wails 버전/CI)
Step 2. apps/folder-watcher 생성 및 기본 실행 확인(단일 앱 dev)
Step 3. Backend 구현
- WatchEntry 모델(설정 저장용) + Runtime 상태 분리
- SQLite 저장:
  - watchers 설정 테이블(또는 JSON 설정 + SQLite jobs만 사용 중 1 선택; “보장”을 위해 jobs는 SQLite 필수)
  - jobs 테이블(queued/running/done/failed)
  - logs 저장(권장: SQLite logs 테이블 또는 jsonl; “조회+필터” 편의상 SQLite logs 권장)
- WatchManager:
  - Add/Update/Delete/List
  - Start/Stop(개별)
  - StartAll/StopAll(선택)
- fsnotify 수집:
  - 이벤트 정규화 + ignore 적용 + create 후보만 job enqueue(Create + Rename 유입)
- 리스캔(10초):
  - 파일 목록 스캔 -> create 후보만 job enqueue
  - skip-existing 옵션 반영(seen set)
  - failed 파일은 자동 enqueue 제외
- worker(워처별 1개):
  - job을 순차 처리
  - 파일 안정성 체크
  - 충돌 이름 변경
  - copy / move(rename 우선, 실패 시 copy+delete)
  - 실패 시 원본 삭제 금지, status=failed, 로그 기록, 실패 알림
  - 성공 시 status=done, 성공 로그, 성공 알림
- 알림:
  - 인앱: Wails EventsEmit으로 프론트에 전달
  - Windows 토스트: NotificationService(Windows 구현 + 타 OS noop)
  - 성공/실패 모두 토스트

Step 4. Frontend 구현(Vite + TS)
- Watchers 화면:
  - 리스트 + Add/Edit/Delete
  - 항목별 Start/Stop 버튼
  - OnStartBehavior 옵션 설정 UI 포함
- Logs 화면:
  - 최근 로그 조회 + 필터(워처/레벨/키워드) 최소 구현
  - failed 상태 표시(추후 Retry 확장 고려)
- 이벤트 구독:
  - 백엔드에서 emit한 이벤트로 인앱 toast 표시 + 화면 갱신
- 폴더 선택:
  - Wails directory picker 사용

Step 5. Monorepo 커맨드/VSCode/CI 연결
- 루트 커맨드로 dev/build/test/lint 동작 확인
- dev:all 실행 시 포트 충돌 가능성이 있으면 해결:
  - Vite 포트를 앱별로 분리(앱별 env 또는 vite config)
  - 문서에 가이드 추가

Step 6. 문서(README)
- bootstrap
- dev:folder-watcher / build:folder-watcher:win
- watcher 추가/Start/옵션 설명
- “누락 0”을 위해 job queue + 10초 리스캔이 동작하는 방식 설명
- 트러블슈팅(권한, 파일 잠김, output 경로 금지 등)

---

## 산출물(Deliverables)
- 최종 디렉터리 트리
- 생성/수정한 주요 파일 목록 + 핵심 내용
- 실행 가능한 루트 명령어 목록 + 예시
- README(실행/빌드/새 watcher 추가/옵션 설명)
- 간단 검증 체크리스트

---

## 수용 기준(Acceptance Tests)
- 루트에서 dev:folder-watcher 1번으로 개발 실행
- watcher 2개 등록 후 각각 Start하면 개별 워칭 시작/정지
- 소스 폴더에 파일을 “새로 추가”하면:
  - 인앱 토스트 + Windows 토스트 발생
  - 로그(UI/로컬)에 기록
  - output + action 설정 시 move/copy 수행
- 드롭/이동 유입(rename)도 create로 처리되어 동일하게 작동
- 충돌 시 넘버링된 파일명으로 저장
- 실패 시 원본 삭제되지 않고, failed로 기록되어 재확인 가능
- 10초 리스캔으로 이벤트 누락이 있어도 결과 누락 없이 처리됨
- build:folder-watcher:win으로 exe 산출물 생성(앱별 경로 분리)

---

## 진행 방식 지시
- 작업 시작 전: “저장소 스캔 결과”를 요약하고 그 기반으로 변경 범위를 확정한다.
- 이후: 단계별로 커밋 단위/파일 단위로 구현한다(가능하면 작은 단위로).
- 기존 구조/스타일을 바꾸지 말고, 필요한 부분만 추가한다.
