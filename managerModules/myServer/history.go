package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// /workload_history 관련 라우터 등록
func RegisterHistoryRoutes(r *gin.Engine) {
	// HTML 페이지 제공
	r.StaticFile("/workload_history", "./workload_history.html")

	// 데이터 API 엔드포인트
	r.GET("/workload_history/data", handleGetWorkloadHistoryData)
}
