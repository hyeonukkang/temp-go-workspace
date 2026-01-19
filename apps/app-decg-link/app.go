package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"monorepo/packages/logger"

	"github.com/google/uuid"
	"github.com/samber/lo"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx


	u := uuid.New()
	logger.Info("uuid.New() result:", u)

	// samber/lo v1.49.1 사용 예제
	nums := []int{1, 2, 3}
	hasTwo := lo.Contains(nums, 2)
	logger.Info("lo.Contains(nums, 2):", hasTwo)

	// API 서버 시작
	go func() {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		http.HandleFunc("/api/fileinfo", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
				return
			}
			var req struct {
				Name string `json:"name"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}
			result := InsertFileInfo(a.ctx, req.Name)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"result": result})
		})
		logger.Info("API 서버 시작:", ":"+port)
		err := http.ListenAndServe(":"+port, nil)
		logger.Error("API 서버 종료:", err)
	}()
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return Greet(name)
}

// InsertFileInfo: 파일명과 업로드 시간 정보를 DB에 저장
func (a *App) InsertFileInfo(fileName string) string {
	return InsertFileInfo(a.ctx, fileName)
}

// GetFileInfoTable: file_info 테이블 전체 조회
func (a *App) GetFileInfoTable() []map[string]interface{} {
	return GetFileInfoTable(a.ctx)
}

// LogButtonClick: 버튼 클릭 로그 기록
func (a *App) LogButtonClick() {
	logger.Info("button click")
}
