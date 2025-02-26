package main

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
	ys "main/ystruct"
)

var db *sql.DB
var BASE_URL = "http://" + os.Getenv("KWARE_IP") + ":" + os.Getenv("KWARE_PORT") + os.Getenv("KWARE_PATH")

func main() {
	r := gin.Default()
	r.GET("/workloadinfo", handleGetWorkloadinfoRequest)
	r.GET("/strato", handleGetStratoRequest)
	r.POST("/submit", handleSubmitRequest)
	r.Run()
}

// [기존 handleGetWorkloadinfoRequest, handleGetStratoRequest 유지]

func handleSubmitRequest(c *gin.Context) {
	var requestData ys.RequestData

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if requestData.Timestamp == "" {
		loc, err := time.LoadLocation("Asia/Seoul")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load location"})
			return
		}
		requestData.Timestamp = time.Now().In(loc).Format("2006-01-02 15:04:05")
	}

	metadataJSON, err := json.Marshal(requestData.Metadata)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize metadata"})
		return
	}

	allocInfo := ys.ResourceAllocInfo{
		Addr:        BASE_URL,
		EncodedYaml: requestData.Yaml,
	}
	ackBody := ReqResourceAllocInfo(allocInfo)

	finalYaml, clusterValue := MadeFinalWorkloadYAML(ackBody, requestData.Yaml)
	finalYamlYAML, err := yaml.Marshal(finalYaml)
	if err != nil {
		log.Fatalf("Error marshaling final YAML: %v", err)
	}
	finalYamlBase64 := base64.StdEncoding.EncodeToString(finalYamlYAML)

	err = sendPostRequest(clusterValue, finalYamlBase64, false)
	if err != nil {
		log.Printf("%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send POST request"})
		return
	}

	_, err = db.Exec("INSERT INTO workload_info (workload_name, yaml, metadata, created_timestamp) VALUES (?, ?, ?, ?)",
		requestData.Metadata["name"], finalYamlBase64, string(metadataJSON), requestData.Timestamp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// Stub functions (다음 커밋에서 구현)
func ReqResourceAllocInfo(allocInfo ys.ResourceAllocInfo) ys.RespResource {
	return ys.RespResource{} // 임시 반환
}

func MadeFinalWorkloadYAML(argBody ys.RespResource, inputYaml string) (map[string]interface{}, string) {
	return nil, "1" // 임시 반환
}

func sendPostRequest(clusterValue string, finalYamlBase64 string, retryValue bool) error {
	return nil // 임시 반환
}
