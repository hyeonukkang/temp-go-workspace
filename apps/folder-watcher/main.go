package folderwatcher
package main


import (
	"context"
	"github.com/wailsapp/wails/v2"
)

func main() {
	wm, err := NewWatchManager("data/folder-watcher.db")
	if err != nil {
		panic(err)
	}
	app := &App{wm: wm}
	wails.Run(&wails.Options{
		Title:    "Folder Watcher",
		Width:    1024,
		Height:   768,
		Assets:   nil, // 프론트엔드 빌드 후 경로로 교체
		Bind:     []interface{}{app},
	})
}

// App: Wails export용 래퍼
type App struct {
	wm *WatchManager
}

// Watcher CRUD
func (a *App) AddWatcher(ctx context.Context, entry *WatchEntry) error         { return a.wm.AddWatcher(entry) }
func (a *App) UpdateWatcher(ctx context.Context, entry *WatchEntry) error      { return a.wm.UpdateWatcher(entry) }
func (a *App) DeleteWatcher(ctx context.Context, id int64) error              { return a.wm.DeleteWatcher(id) }
func (a *App) ListWatchers(ctx context.Context) ([]*WatchEntry, error)        { return a.wm.ListWatchers() }

// Job/Log
func (a *App) ListJobs(ctx context.Context, watcherID int64, status string) ([]*Job, error) { return a.wm.ListJobs(watcherID, status) }
func (a *App) ListLogs(ctx context.Context, watcherID int64, level string) ([]*LogEntry, error) { return a.wm.ListLogs(watcherID, level) }

// Start/Stop
func (a *App) StartWatcher(ctx context.Context, id int64) error { return a.wm.StartWatcher(id) }
func (a *App) StopWatcher(ctx context.Context, id int64) error  { return a.wm.StopWatcher(id) }
func (a *App) StartAll(ctx context.Context)                     { a.wm.StartAll() }
func (a *App) StopAll(ctx context.Context)                      { a.wm.StopAll() }

// 상태조회
func (a *App) GetWatcherStatus(ctx context.Context, id int64) (string, error) { return a.wm.GetWatcherStatus(id) }
func (a *App) GetJobStatus(ctx context.Context, jobID int64) (string, error)  { return a.wm.GetJobStatus(jobID) }
