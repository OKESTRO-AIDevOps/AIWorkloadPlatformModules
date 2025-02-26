package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	ys "main/ystruct"
)

var db *sql.DB // 전역 DB 변수 (실제 구현에서는 초기화 필요)

func main() {
	r := gin.Default()
	r.GET("/workloadinfo", handleGetWorkloadinfoRequest)
	r.GET("/strato", handleGetStratoRequest)
	r.Run()
}

// GET 요청 처리 함수
func handleGetWorkloadinfoRequest(c *gin.Context) {
	var results []ys.WorkloadInfo

	query := "SELECT workload_name, yaml, metadata, created_timestamp FROM workload_info"
	rows, err := db.Query(query)
	if err != nil {
		c.JSON(500, gin.H{"error": "Database query failed"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var result ys.WorkloadInfo
		if err := rows.Scan(&result.WorkloadName, &result.YAML, &result.Metadata, &result.CreatedTimestamp); err != nil {
			c.JSON(500, gin.H{"error": "Row scan failed"})
			return
		}
		results = append(results, result)
	}

	response := gin.H{
		"respond": results,
	}
	c.JSON(200, response)
}

func handleGetStratoRequest(c *gin.Context) {
	var results []ys.Strato

	query := "SELECT mlid, yaml, data FROM strato"
	rows, err := db.Query(query)
	if err != nil {
		c.JSON(500, gin.H{"error": "Database query failed"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var result ys.Strato
		if err := rows.Scan(&result.MlId, &result.YAML, &result.Data); err != nil {
			c.JSON(500, gin.H{"error": "Row scan failed"})
			return
		}
		results = append(results, result)
	}

	response := gin.H{
		"respond": results,
	}
	c.JSON(200, response)
}
