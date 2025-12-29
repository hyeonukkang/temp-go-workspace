package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

// InsertFileInfo: 파일명과 업로드 시간 정보를 DB에 저장
func InsertFileInfo(ctx context.Context, fileName string) string {
       conn, err := pgx.Connect(ctx, "postgres://postgres:postgres@localhost:5432/decg_link")
       if err != nil {
	       log.Println("DB 연결 실패:", err)
	       return "DB 연결 실패: " + err.Error()
       }
       defer conn.Close(ctx)

       // 테이블 생성(없으면)
       _, err = conn.Exec(ctx, `CREATE TABLE IF NOT EXISTS file_info (
	       id SERIAL PRIMARY KEY,
	       file_name TEXT NOT NULL,
	       uploaded_at TIMESTAMP NOT NULL
       )`)
       if err != nil {
	       log.Println("테이블 생성 실패:", err)
	       return "테이블 생성 실패: " + err.Error()
       }

       // 파일 정보 저장
       now := time.Now()
       _, err = conn.Exec(ctx, `INSERT INTO file_info (file_name, uploaded_at) VALUES ($1, $2)`, fileName, now)
       if err != nil {
	       log.Println("데이터 저장 실패:", err)
	       return "데이터 저장 실패: " + err.Error()
       }
       return fmt.Sprintf("저장 완료: %s (%s)", fileName, now.Format(time.RFC3339))
}

// GetFileInfoTable: file_info 테이블 전체 조회
func GetFileInfoTable(ctx context.Context) []map[string]interface{} {
       conn, err := pgx.Connect(ctx, "postgres://postgres:postgres@localhost:5432/decg_link")
       if err != nil {
	       log.Println("DB 연결 실패:", err)
	       return nil
       }
       defer conn.Close(ctx)

       rows, err := conn.Query(ctx, `SELECT id, file_name, uploaded_at FROM file_info ORDER BY id DESC`)
       if err != nil {
	       log.Println("테이블 조회 실패:", err)
	       return nil
       }
       defer rows.Close()

       var result []map[string]interface{}
       for rows.Next() {
	       var id int
	       var fileName string
	       var uploadedAt time.Time
	       err := rows.Scan(&id, &fileName, &uploadedAt)
	       if err != nil {
		       continue
	       }
	       result = append(result, map[string]interface{}{
		       "id": id,
		       "file_name": fileName,
		       "uploaded_at": uploadedAt.Format(time.RFC3339),
	       })
       }
       return result
}
