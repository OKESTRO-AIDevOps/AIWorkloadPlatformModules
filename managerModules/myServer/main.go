package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// 데이터베이스 핸들러
var db *sql.DB

func main() {
	// Gin 라우터 설정
	r := gin.Default()
	// 데이터베이스 초기화
	initDatabase()
	// CORS 설정
	r.Use(setupCORS())
	// 라우트 등록
	registerRoutes(r)
	// 서버 실행
	r.Run("0.0.0.0:8080")
}