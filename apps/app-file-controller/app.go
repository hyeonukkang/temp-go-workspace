package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"monorepo/packages/utils"

	"github.com/jackc/pgx/v5"
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
       // API 서버 시작 (포트 환경변수 지원)
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
	       log.Println("API 서버 시작:", ":"+port)
	       log.Fatal(http.ListenAndServe(":"+port, nil))
       }()
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
       // decg-link API로 이름 전송
       go func(name string) {
	       payload := map[string]string{"name": name}
	       body, _ := json.Marshal(payload)
	       resp, err := http.Post("http://localhost:8080/api/fileinfo", "application/json", bytes.NewBuffer(body))
	       if err != nil {
		       log.Println("decg-link API 요청 실패:", err)
		       return
	       }
	       defer resp.Body.Close()
	       var res map[string]string
	       if err := json.NewDecoder(resp.Body).Decode(&res); err == nil {
		       log.Println("decg-link 응답:", res["result"])
	       }
       }(name)
       return fmt.Sprintf("Hello %s, It's show time!", name)
}

// PostgresTestQuery: PostgreSQL 연결 및 샘플 쿼리 결과 반환
func (a *App) PostgresTestQuery() string {
       conn, err := pgx.Connect(a.ctx, "postgres://postgres:postgres@localhost:5432/postgres")
       if err != nil {
	       log.Println("DB 연결 실패:", err)
	       return "DB 연결 실패: " + err.Error()
       }
       defer conn.Close(a.ctx)

       var greeting string
       err = conn.QueryRow(a.ctx, "select 'Hello, PostgreSQL from Go!' as greeting").Scan(&greeting)
       if err != nil {
	       log.Println("쿼리 실패:", err)
	       return "쿼리 실패: " + err.Error()
       }
       return greeting
}

func InsertFileInfo(ctx context.Context, name string) string {
    // 실제 저장 로직이 필요하다면 여기에 구현
		log.Println("InsertFileInfo called with name:", name)
		sum := utils.Add(3, 5)
    log.Println("sum:", sum)
    return "ok"
}