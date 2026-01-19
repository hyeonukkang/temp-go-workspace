import (
	"github.com/wailsapp/wails/v2/pkg/runtime"
)
// Wails 이벤트 송출(인앱/Windows 토스트)
func (wm *WatchManager) emitEvent(ctx context.Context, watcherID int64, eventType, message string, success bool) {
	// 인앱 이벤트
	runtime.EventsEmit(ctx, "watcher:event", map[string]interface{}{
		"watcherId": watcherID,
		"type": eventType,
		"message": message,
		"success": success,
	})
	// Windows 토스트(성공/실패 모두)
	// TODO: NotificationService 연동(Windows만)
}

// StartAll/StopAll(선택)
func (wm *WatchManager) StartAll() {
	// TODO: 전체 워치 Start
}
func (wm *WatchManager) StopAll() {
	// TODO: 전체 워치 Stop
}

// 워치/작업/로그 상태 조회(프론트엔드 연동용)
func (wm *WatchManager) GetWatcherStatus(id int64) (string, error) {
	// TODO: running/stopped 등 반환
	return "", nil
}
func (wm *WatchManager) GetJobStatus(jobID int64) (string, error) {
	// TODO: job 상태 반환
	return "", nil
}
import (
	"io"
	"path/filepath"
	"fmt"
)
// 파일 안정성 체크: 크기/mtime N회(3회) 연속 동일, 최대 5초 대기
func waitFileStable(path string) bool {
	var prevSize int64 = -1
	var prevMod int64 = -1
	for i := 0; i < 5; i++ {
		fi, err := os.Stat(path)
		if err != nil || fi.IsDir() {
			return false
		}
		size := fi.Size()
		mod := fi.ModTime().UnixNano()
		if size == prevSize && mod == prevMod {
			return true
		}
		prevSize = size
		prevMod = mod
		time.Sleep(1 * time.Second)
	}
	return false
}
// seenSet: skip-existing 모드에서 이미 본 파일 관리
type seenSet struct {
	mu sync.Mutex
	files map[string]struct{}
}

func newSeenSet() *seenSet {
	return &seenSet{files: make(map[string]struct{})}
}
func (s *seenSet) Add(name string) { s.mu.Lock(); s.files[name] = struct{}{}; s.mu.Unlock() }
func (s *seenSet) Remove(name string) { s.mu.Lock(); delete(s.files, name); s.mu.Unlock() }
func (s *seenSet) Has(name string) bool { s.mu.Lock(); defer s.mu.Unlock(); _, ok := s.files[name]; return ok }
func (s *seenSet) Reset(names []string) { s.mu.Lock(); s.files = make(map[string]struct{}); for _, n := range names { s.files[n] = struct{}{} }; s.mu.Unlock() }
// 기본 ignore 패턴
var defaultIgnorePatterns = []string{"*.tmp", "*.swp", "*~", ".DS_Store", "Thumbs.db", "~$*"}

// 파일명 패턴 매칭(간단 glob)
func matchAnyPattern(name string, patterns []string) bool {
	for _, pat := range patterns {
		// 단순 접미사/포함 매칭(고급 glob 필요시 확장)
		if len(pat) > 1 && pat[0] == '*' && len(name) >= len(pat)-1 && name[len(name)-len(pat)+1:] == pat[1:] {
			return true
		}
		if pat == name {
			return true
		}
	}
	return false
}
import (
	"context"
	"sync"
	"github.com/fsnotify/fsnotify"
)
// 워치 실행/중지 및 런타임 상태 관리
type watcherRuntime struct {
	cancel context.CancelFunc
	running bool
}

// WatchManager에 런타임 상태 맵 추가
type WatchManager struct {
	db *sql.DB
	mu sync.Mutex
	runtimes map[int64]*watcherRuntime // watcherID -> runtime
}

// StartWatcher: 워치 항목 실행(폴더 감시, 리스캔, 워커 시작)
func (wm *WatchManager) StartWatcher(id int64) error {
	wm.mu.Lock()
	if wm.runtimes == nil {
		wm.runtimes = make(map[int64]*watcherRuntime)
	}
	if rt, ok := wm.runtimes[id]; ok && rt.running {
		wm.mu.Unlock()
		return nil // 이미 실행 중
	}
	ctx, cancel := context.WithCancel(context.Background())
	wm.runtimes[id] = &watcherRuntime{cancel: cancel, running: true}
	wm.mu.Unlock()

	go wm.runWatcher(ctx, id)
	return nil
}

// StopWatcher: 워치 항목 중지
func (wm *WatchManager) StopWatcher(id int64) error {
	wm.mu.Lock()
	if rt, ok := wm.runtimes[id]; ok && rt.running {
		rt.cancel()
		rt.running = false
	}
	wm.mu.Unlock()
	return nil
}

// runWatcher: 폴더 감시, 리스캔, 워커 루프 실행(비동기)
func (wm *WatchManager) runWatcher(ctx context.Context, watcherID int64) {
	// TODO: fsnotify 감시, 10초 리스캔, 작업 큐 polling, worker 등 구현
	// 예시: go wm.watchFsEvents(ctx, watcherID)
	// 예시: go wm.rescanLoop(ctx, watcherID)
	// 예시: go wm.workerLoop(ctx, watcherID)
}

// watchFsEvents: fsnotify로 폴더 이벤트 감지
func (wm *WatchManager) watchFsEvents(ctx context.Context, watcherID int64) {
	entry, err := wm.getWatcherByID(watcherID)
	if err != nil || entry == nil {
		return
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return
	}
	defer watcher.Close()
	err = watcher.Add(entry.SourcePath)
	if err != nil {
		return
	}
	ignore := append(defaultIgnorePatterns, entry.IgnorePatterns...)
	for {
		select {
		case <-ctx.Done():
			return
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			// 1단계 폴더만, 디렉터리 무시
			fi, err := os.Stat(event.Name)
			if err != nil || fi.IsDir() {
				continue
			}
			base := fi.Name()
			if matchAnyPattern(base, ignore) {
				continue
			}
			// Create/Rename 유입만 처리
			if event.Op&fsnotify.Create == fsnotify.Create || event.Op&fsnotify.Rename == fsnotify.Rename {
				// Rename은 150ms 후 파일 존재 시 Create로 간주
				if event.Op&fsnotify.Rename == fsnotify.Rename {
					time.Sleep(150 * time.Millisecond)
					if _, err := os.Stat(event.Name); err != nil {
						continue
					}
				}
				// job enqueue
				_ = wm.EnqueueJob(&Job{
					WatcherID: watcherID,
					FilePath:  event.Name,
					Status:    "queued",
				})
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			// 로그 기록
			_ = wm.AddLog(&LogEntry{WatcherID: watcherID, Level: "error", Message: err.Error()})
		}
	}
}
// getWatcherByID: 단일 워치 항목 조회(내부용)
func (wm *WatchManager) getWatcherByID(id int64) (*WatchEntry, error) {
	rows, err := wm.db.Query(`SELECT id, source_path, output_path, action, on_start_behavior, ignore_patterns, created_at, updated_at FROM watchers WHERE id=?`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		var e WatchEntry
		var ignore string
		if err := rows.Scan(&e.ID, &e.SourcePath, &e.OutputPath, &e.Action, &e.OnStartBehavior, &ignore, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		if ignore != "" {
			_ = json.Unmarshal([]byte(ignore), &e.IgnorePatterns)
		}
		return &e, nil
	}
	return nil, nil
}

// rescanLoop: 10초마다 폴더 1단계 리스캔, 신규 파일 job enqueue
func (wm *WatchManager) rescanLoop(ctx context.Context, watcherID int64) {
	entry, err := wm.getWatcherByID(watcherID)
	if err != nil || entry == nil {
		return
	}
	ignore := append(defaultIgnorePatterns, entry.IgnorePatterns...)
	seen := newSeenSet()
	// skip-existing이면 최초 스냅샷 생성
	if entry.OnStartBehavior == "skip-existing" {
		files, _ := os.ReadDir(entry.SourcePath)
		var names []string
		for _, f := range files {
			if f.IsDir() || matchAnyPattern(f.Name(), ignore) {
				continue
			}
			names = append(names, f.Name())
		}
		seen.Reset(names)
	}
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			files, _ := os.ReadDir(entry.SourcePath)
			cur := make(map[string]struct{})
			for _, f := range files {
				if f.IsDir() || matchAnyPattern(f.Name(), ignore) {
					continue
				}
				cur[f.Name()] = struct{}{}
				// 신규 파일만 enqueue
				if entry.OnStartBehavior == "skip-existing" && seen.Has(f.Name()) {
					continue
				}
				// 이미 job 큐에 있는지(queued/running/done/failed 중 queued/running만 제외)
				jobs, _ := wm.ListJobs(watcherID, "queued")
				found := false
				for _, j := range jobs {
					if j.FilePath == f.Name() {
						found = true
						break
					}
				}
				if !found {
					_ = wm.EnqueueJob(&Job{
						WatcherID: watcherID,
						FilePath:  f.Name(),
						Status:    "queued",
					})
				}
			}
			// skip-existing: 폴더에 없는 파일은 seen에서 제거
			if entry.OnStartBehavior == "skip-existing" {
				for name := range seen.files {
					if _, ok := cur[name]; !ok {
						seen.Remove(name)
					}
				}
				for name := range cur {
					seen.Add(name)
				}
			}
		}
	}
}

// workerLoop: job 큐 polling 및 파일 처리(move/copy 등)
func (wm *WatchManager) workerLoop(ctx context.Context, watcherID int64) {
	entry, err := wm.getWatcherByID(watcherID)
	if err != nil || entry == nil {
		return
	}
	for {
		select {
		case <-ctx.Done():
			return
		default:
			jobs, _ := wm.ListJobs(watcherID, "queued")
			for _, job := range jobs {
				// 파일 안정성 체크
				if !waitFileStable(job.FilePath) {
					wm.UpdateJobStatus(job.ID, "failed", "file unstable or missing")
					wm.AddLog(&LogEntry{WatcherID: watcherID, Level: "error", Message: fmt.Sprintf("file unstable: %s", job.FilePath)})
					continue
				}
				// move/copy/none
				outPath := entry.OutputPath
				action := entry.Action
				if outPath == "" || action == "none" {
					wm.UpdateJobStatus(job.ID, "done", "")
					wm.AddLog(&LogEntry{WatcherID: watcherID, Level: "info", Message: fmt.Sprintf("no action for %s", job.FilePath)})
					continue
				}
				// outputPath가 sourcePath 내부/동일 금지
				absSrc, _ := filepath.Abs(entry.SourcePath)
				absOut, _ := filepath.Abs(outPath)
				if absSrc == absOut || len(absOut) > len(absSrc) && absOut[:len(absSrc)] == absSrc {
					wm.UpdateJobStatus(job.ID, "failed", "outputPath inside sourcePath")
					wm.AddLog(&LogEntry{WatcherID: watcherID, Level: "error", Message: "outputPath inside sourcePath"})
					continue
				}
				// 파일명 충돌 회피(넘버링)
				base := filepath.Base(job.FilePath)
				dst := filepath.Join(outPath, base)
				i := 1
				for {
					if _, err := os.Stat(dst); os.IsNotExist(err) {
						break
					}
					dst = filepath.Join(outPath, fmt.Sprintf("%s(%d)%s", base[:len(base)-len(filepath.Ext(base))], i, filepath.Ext(base)))
					i++
				}
				// move/copy
				var opErr error
				if action == "move" {
					opErr = os.Rename(job.FilePath, dst)
					if opErr != nil {
						// 실패 시 copy+delete
						opErr = copyFile(job.FilePath, dst)
						if opErr == nil {
							_ = os.Remove(job.FilePath)
						}
					}
				} else if action == "copy" {
					opErr = copyFile(job.FilePath, dst)
				}
				if opErr != nil {
					wm.UpdateJobStatus(job.ID, "failed", opErr.Error())
					wm.AddLog(&LogEntry{WatcherID: watcherID, Level: "error", Message: fmt.Sprintf("fail %s: %v", job.FilePath, opErr)})
				} else {
					wm.UpdateJobStatus(job.ID, "done", "")
					wm.AddLog(&LogEntry{WatcherID: watcherID, Level: "info", Message: fmt.Sprintf("%s %s -> %s", action, job.FilePath, dst)})
				}
			}
			time.Sleep(2 * time.Second)
		}
	}
}

// 파일 복사 유틸
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}
import (
	"encoding/json"
)
// Job 관리
func (wm *WatchManager) EnqueueJob(job *Job) error {
	now := time.Now()
	job.CreatedAt = now
	job.UpdatedAt = now
	res, err := wm.db.Exec(`INSERT INTO jobs (watcher_id, file_path, status, created_at, updated_at, error_msg) VALUES (?, ?, ?, ?, ?, ?)`,
		job.WatcherID, job.FilePath, job.Status, job.CreatedAt, job.UpdatedAt, job.ErrorMsg)
	if err != nil {
		return err
	}
	job.ID, _ = res.LastInsertId()
	return nil
}

func (wm *WatchManager) UpdateJobStatus(id int64, status, errorMsg string) error {
	_, err := wm.db.Exec(`UPDATE jobs SET status=?, error_msg=?, updated_at=? WHERE id=?`, status, errorMsg, time.Now(), id)
	return err
}

func (wm *WatchManager) ListJobs(watcherID int64, status string) ([]*Job, error) {
	query := `SELECT id, watcher_id, file_path, status, created_at, updated_at, error_msg FROM jobs WHERE watcher_id=?`
	args := []interface{}{watcherID}
	if status != "" {
		query += " AND status=?"
		args = append(args, status)
	}
	rows, err := wm.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []*Job
	for rows.Next() {
		var j Job
		if err := rows.Scan(&j.ID, &j.WatcherID, &j.FilePath, &j.Status, &j.CreatedAt, &j.UpdatedAt, &j.ErrorMsg); err != nil {
			return nil, err
		}
		result = append(result, &j)
	}
	return result, nil
}

// Log 관리
func (wm *WatchManager) AddLog(log *LogEntry) error {
	log.CreatedAt = time.Now()
	_, err := wm.db.Exec(`INSERT INTO logs (watcher_id, level, message, created_at) VALUES (?, ?, ?, ?)`,
		log.WatcherID, log.Level, log.Message, log.CreatedAt)
	return err
}

func (wm *WatchManager) ListLogs(watcherID int64, level string) ([]*LogEntry, error) {
	query := `SELECT id, watcher_id, level, message, created_at FROM logs WHERE watcher_id=?`
	args := []interface{}{watcherID}
	if level != "" {
		query += " AND level=?"
		args = append(args, level)
	}
	rows, err := wm.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []*LogEntry
	for rows.Next() {
		var l LogEntry
		if err := rows.Scan(&l.ID, &l.WatcherID, &l.Level, &l.Message, &l.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, &l)
	}
	return result, nil
}
// Package main - Folder Watcher Wails 앱 진입점
package main

import (
	"database/sql"
	"encoding/json"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// WatchEntry: 폴더 워치 설정(영구 저장)
// WatchManager: 워치/작업/로그 관리
type WatchManager struct {
	db *sql.DB
}

// NewWatchManager: SQLite DB 초기화 및 테이블 생성
func NewWatchManager(dbPath string) (*WatchManager, error) {
	if err := os.MkdirAll("data", 0755); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	// 테이블 생성
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS watchers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			source_path TEXT NOT NULL,
			output_path TEXT,
			action TEXT NOT NULL,
			on_start_behavior TEXT NOT NULL,
			ignore_patterns TEXT,
			created_at DATETIME,
			updated_at DATETIME
		);`,
		`CREATE TABLE IF NOT EXISTS jobs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			watcher_id INTEGER,
			file_path TEXT NOT NULL,
			status TEXT NOT NULL,
			created_at DATETIME,
			updated_at DATETIME,
			error_msg TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			watcher_id INTEGER,
			level TEXT NOT NULL,
			message TEXT NOT NULL,
			created_at DATETIME
		);`,
	}
	for _, stmt := range stmts {
		if _, err := db.Exec(stmt); err != nil {
			db.Close()
			return nil, err
		}
	}
	return &WatchManager{db: db}, nil
}

// AddWatcher: 워치 항목 추가
func (wm *WatchManager) AddWatcher(entry *WatchEntry) error {
	now := time.Now()
	entry.CreatedAt = now
	entry.UpdatedAt = now
	ignore := ""
	if len(entry.IgnorePatterns) > 0 {
		// JSON 직렬화
		b, _ := json.Marshal(entry.IgnorePatterns)
		ignore = string(b)
	}
	res, err := wm.db.Exec(`INSERT INTO watchers (source_path, output_path, action, on_start_behavior, ignore_patterns, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		entry.SourcePath, entry.OutputPath, entry.Action, entry.OnStartBehavior, ignore, entry.CreatedAt, entry.UpdatedAt)
	if err != nil {
		return err
	}
	entry.ID, _ = res.LastInsertId()
	return nil
}

// UpdateWatcher: 워치 항목 수정
func (wm *WatchManager) UpdateWatcher(entry *WatchEntry) error {
	entry.UpdatedAt = time.Now()
	ignore := ""
	if len(entry.IgnorePatterns) > 0 {
		b, _ := json.Marshal(entry.IgnorePatterns)
		ignore = string(b)
	}
	_, err := wm.db.Exec(`UPDATE watchers SET source_path=?, output_path=?, action=?, on_start_behavior=?, ignore_patterns=?, updated_at=? WHERE id=?`,
		entry.SourcePath, entry.OutputPath, entry.Action, entry.OnStartBehavior, ignore, entry.UpdatedAt, entry.ID)
	return err
}

// DeleteWatcher: 워치 항목 삭제
func (wm *WatchManager) DeleteWatcher(id int64) error {
	_, err := wm.db.Exec(`DELETE FROM watchers WHERE id=?`, id)
	return err
}

// ListWatchers: 워치 목록 조회
func (wm *WatchManager) ListWatchers() ([]*WatchEntry, error) {
	rows, err := wm.db.Query(`SELECT id, source_path, output_path, action, on_start_behavior, ignore_patterns, created_at, updated_at FROM watchers`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []*WatchEntry
	for rows.Next() {
		var e WatchEntry
		var ignore string
		if err := rows.Scan(&e.ID, &e.SourcePath, &e.OutputPath, &e.Action, &e.OnStartBehavior, &ignore, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		if ignore != "" {
			_ = json.Unmarshal([]byte(ignore), &e.IgnorePatterns)
		}
		result = append(result, &e)
	}
	return result, nil
}
type WatchEntry struct {
	ID             int64     `json:"id"`
	SourcePath     string    `json:"sourcePath"`
	OutputPath     string    `json:"outputPath,omitempty"`
	Action         string    `json:"action"` // move | copy | none
	OnStartBehavior string   `json:"onStartBehavior"` // skip-existing | process-existing
	IgnorePatterns []string  `json:"ignorePatterns,omitempty"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// Job: 파일 처리 작업 큐 항목(SQLite)
type Job struct {
	ID         int64     `json:"id"`
	WatcherID  int64     `json:"watcherId"`
	FilePath   string    `json:"filePath"`
	Status     string    `json:"status"` // queued | running | done | failed
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	ErrorMsg   string    `json:"errorMsg,omitempty"`
}

// LogEntry: 이벤트/처리 로그(SQLite)
type LogEntry struct {
	ID         int64     `json:"id"`
	WatcherID  int64     `json:"watcherId"`
	Level      string    `json:"level"` // info | warn | error
	Message    string    `json:"message"`
	CreatedAt  time.Time `json:"createdAt"`
}
