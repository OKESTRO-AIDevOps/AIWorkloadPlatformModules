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

// /workload_history/data 요청 처리 핸들러
func handleGetWorkloadHistoryData(c *gin.Context) {
	name := c.DefaultQuery("name", "")
	startDate := c.DefaultQuery("start_date", "")
	endDate := c.DefaultQuery("end_date", "")

	// 페이지 번호와 한 페이지당 항목 수를 쿼리 파라미터로 받음
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "20")

	// 페이지 번호와 한 페이지당 항목 수를 정수로 변환
	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		pageInt = 1
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 1 {
		limitInt = 20
	}

	// OFFSET 계산
	offset := (pageInt - 1) * limitInt

	// 기본 쿼리 시작
	query := "SELECT workload_name, yaml, metadata, created_timestamp FROM workload_info WHERE 1=1"
	countQuery := "SELECT COUNT(*) FROM workload_info WHERE 1=1"
