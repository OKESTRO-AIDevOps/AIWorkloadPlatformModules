package main

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
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
	r.POST("/submit-resource", handleSubmitResourceRequest)
	r.Run()
}

func ReqResourceAllocInfo(allocInfo ys.ResourceAllocInfo) ys.RespResource {
	data, err := base64.StdEncoding.DecodeString(allocInfo.EncodedYaml)
	if err != nil {
		log.Printf("Failed to decode base64 data: %s", err)
		return ys.RespResource{}
	}

	var workflow ys.Workflow
	err = yaml.Unmarshal(data, &workflow)
	if err != nil {
		log.Printf("Failed to unmarshal YAML data: %s", err)
		return ys.RespResource{}
	}

	reqJson := ys.ReqResource{}
	uuid := "dmkim"
	currentTime := time.Now()
	nowTime := currentTime.Format("2006-01-02 15:04:05")

	reqJson.Version = "0.12"
	reqJson.Request.Name = workflow.Metadata.GenerateName
	reqJson.Request.ID = uuid
	reqJson.Request.Date = nowTime

	for _, value := range workflow.Spec.Templates {
		if value.Container != nil {
			tmpContainer := ys.Container{
				Name: value.Name,
				Resources: ys.Resources{
					Requests: ys.ResourceDetails{
						CPU:    value.Container.Resources.Requests.CPU,
						GPU:    value.Container.Resources.Requests.GPU,
						Memory: value.Container.Resources.Requests.Memory,
					},
					Limits: ys.ResourceDetails{
						GPU:    value.Container.Resources.Limits.NvidiaGPU,
						CPU:    value.Container.Resources.Limits.CPU,
						Memory: value.Container.Resources.Limits.Memory,
					},
				},
			}
			reqJson.Request.Containers = append(reqJson.Request.Containers, tmpContainer)
		}
	}

	var ackBody ys.RespResource
	ack, body := SEND_REST_DATA(allocInfo.Addr, reqJson)
	if ack.StatusCode == http.StatusOK {
		err = json.Unmarshal([]byte(body), &ackBody)
		if err != nil {
			log.Printf("Failed to unmarshal ack body: %s", err)
		}
	}
	return ackBody
}

func MadeFinalWorkloadYAML(argBody ys.RespResource, inputYaml string) (map[string]interface{}, string) {
	clusterValue := argBody.Response.Cluster
	yamlFile, err := base64.StdEncoding.DecodeString(inputYaml)
	if err != nil {
		log.Fatalf("Error decoding Base64 YAML data: %v", err)
	}

	var data map[string]interface{}
	err = yaml.Unmarshal(yamlFile, &data)
	if err != nil {
		log.Fatalf("Error unmarshalling YAML data: %v", err)
	}

	spec, ok := data["spec"].(map[interface{}]interface{})
	if ok {
		templates, ok := spec["templates"].([]interface{})
		if ok {
			for _, template := range templates {
				templateMap, ok := template.(map[interface{}]interface{})
				if ok {
					for _, val := range argBody.Response.Containers {
						if templateMap["name"] == val.Name {
							templateMap["nodeSelector"] = ys.NodeSelect{Node: val.Node}
						}
					}
				}
			}
		}
	}
	return data, clusterValue
}

// Stub functions
func SEND_REST_DATA(addr string, reqJson ys.ReqResource) (*http.Response, string) {
	return nil, "" // 임시 반환
}

func handleSubmitResourceRequest(c *gin.Context) {
	var requestResourceData ys.RequestResourceData

	if err := c.ShouldBindJSON(&requestResourceData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var base64Yaml string
	err := db.QueryRow("SELECT yaml FROM workload_info WHERE workload_name = ? ORDER BY created_timestamp DESC LIMIT 1", requestResourceData.Name).Scan(&base64Yaml)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Workload not found"})
		return
	}

	var existingMetadata string
	err = db.QueryRow("SELECT metadata FROM workload_info WHERE workload_name = ? ORDER BY created_timestamp DESC LIMIT 1", requestResourceData.Name).Scan(&existingMetadata)
	if err != nil {
		if err == sql.ErrNoRows {
			existingMetadata = "{}"
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	yamlData, err := base64.StdEncoding.DecodeString(base64Yaml)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode YAML"})
		return
	}

	var yamlMap map[string]interface{}
	err = yaml.Unmarshal(yamlData, &yamlMap)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse YAML"})
		return
	}

	spec, ok := yamlMap["spec"].(map[interface{}]interface{})
	if ok {
		templates, ok := spec["templates"].([]interface{})
		if ok {
			for i, container := range requestResourceData.Containers {
				template, _ := templates[i].(map[interface{}]interface{})
				containerMap, _ := template["container"].(map[interface{}]interface{})
				resources, _ := containerMap["resources"].(map[interface{}]interface{})
				resources["requests"] = map[string]string{"cpu": container.Resources.Requests.CPU, "memory": container.Resources.Requests.Memory}
				resources["limits"] = map[string]string{"cpu": container.Resources.Limits.CPU, "memory": container.Resources.Limits.Memory}
			}
		}
	}

	modifiedYaml, err := yaml.Marshal(yamlMap)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal modified YAML"})
		return
	}
	finalYamlBase64 := base64.StdEncoding.EncodeToString(modifiedYaml)

	clusterValue := "1"
	err = sendPostRequest(clusterValue, finalYamlBase64, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send POST request"})
		return
	}

	loc, _ := time.LoadLocation("Asia/Seoul")
	_, err = db.Exec("INSERT INTO workload_info (workload_name, yaml, metadata, created_timestamp) VALUES (?, ?, ?, ?)",
		requestResourceData.Name, finalYamlBase64, existingMetadata, time.Now().In(loc).Format("2006-01-02 15:04:05"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func sendPostRequest(clusterValue string, finalYamlBase64 string, retryValue bool) error {
	wrapperIp := os.Getenv("WRAPPER_IP")
	wrapperPort := os.Getenv("WRAPPER_PORT")
	wrapperPath := os.Getenv("WRAPPER_PATH")
	address := "http://" + wrapperIp + ":" + wrapperPort + wrapperPath

	wrapperData := ys.WrapperData{
		Cluster: clusterValue,
		Yaml:    finalYamlBase64,
		Retry:   retryValue,
	}

	postJSON, err := json.Marshal(wrapperData)
	if err != nil {
		return fmt.Errorf("failed to create JSON for POST request: %v", err)
	}

	req, err := http.NewRequest("POST", address, bytes.NewBuffer(postJSON))
	if err != nil {
		return fmt.Errorf("failed to create POST request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}
	return nil
}
