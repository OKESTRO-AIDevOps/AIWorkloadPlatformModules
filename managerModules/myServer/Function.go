package main

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	ys "main/ystruct"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
)

var BASE_URL = "http://" + os.Getenv("KWARE_IP") + ":" + os.Getenv("KWARE_PORT") + os.Getenv("KWARE_PATH")

// GET 요청 처리 함수
func handleGetWorkloadinfoRequest(c *gin.Context) {
	var results []ys.WorkloadInfo // 구조체를 사용하여 결과 슬라이스 정의

	query := "SELECT workload_name, yaml, metadata, created_timestamp FROM workload_info"
	rows, err := db.Query(query)
	if err != nil {
		c.JSON(500, gin.H{"error": "Database query failed"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var result ys.WorkloadInfo // 공통 구조체를 사용
		if err := rows.Scan(&result.WorkloadName, &result.YAML, &result.Metadata, &result.CreatedTimestamp); err != nil {
			c.JSON(500, gin.H{"error": "Row scan failed"})
			return
		}
		results = append(results, result) // 동일한 구조체 타입을 사용
	}

	// 응답 구조 수정
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
		results = append(results, result) // 동일한 구조체 타입을 사용
	}

	// 응답 구조 수정
	response := gin.H{
		"respond": results,
	}

	c.JSON(200, response)
}

// POST 요청 처리 함수
func handleSubmitRequest(c *gin.Context) {
	var requestData ys.RequestData

	// 요청 데이터 바인딩
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 클라이언트가 timestamp를 제공하지 않는 경우 현재 시간으로 설정
	if requestData.Timestamp == "" {
		loc, err := time.LoadLocation("Asia/Seoul")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load location"})
			return
		}
		requestData.Timestamp = time.Now().In(loc).Format("2006-01-02 15:04:05")
	}

	// metadata를 JSON 문자열로 변환
	metadataJSON, err := json.Marshal(requestData.Metadata)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize metadata"})
		return
	}

	// Resource Allocation 정보 요청
	var checkpointVal bool
	rawValue := requestData.Metadata["checkpoint"]

	switch v := rawValue.(type) {
	case bool:
		checkpointVal = v
	case string:
		parsedVal, err := strconv.ParseBool(v)
		if err != nil {
			log.Printf("checkpoint string value is not a valid bool: %v", rawValue)
		} else {
			checkpointVal = parsedVal
		}
	default:
		log.Printf("checkpoint is not a valid bool or string: %v", rawValue)
	}

	workloadTypeVal, ok := requestData.Metadata["workloadType"].(string)
	if !ok {
    	log.Printf("workloadType is not a valid string: %v", requestData.Metadata["workloadType"])
	}
	
	var workloadFeatureVal int = 0
	rawVal := requestData.Metadata["workloadFeature"]
	if num, ok := rawVal.(float64); ok {
    	workloadFeatureVal = int(num)
	} else {
   		log.Printf("workloadFeature is not a valid value: %v", rawVal)
	}

	allocInfo := ys.ResourceAllocInfo{
		Addr:            BASE_URL,
		EncodedYaml:     requestData.Yaml,
		Checkpoint:      checkpointVal,
		WorkloadTypeVal:    workloadTypeVal,
		WorkloadFeatureVal: workloadFeatureVal,
	}
	ackBody := ReqResourceAllocInfo(allocInfo)

	// 최종 워크로드 YAML 생성
	finalYaml, clusterValue := MadeFinalWorkloadYAML(ackBody, requestData.Yaml)
	// finalYaml을 YAML 문자열로 변환
	finalYamlYAML, err := yaml.Marshal(finalYaml)
	if err != nil {
		log.Fatalf("Error marshaling final YAML: %v", err)
	}
	// finalYaml을 Base64로 인코딩
	finalYamlBase64 := base64.StdEncoding.EncodeToString(finalYamlYAML)
	// cluster 값과 인코딩된 값을 Yaml Wrapper로 전송
	// POST 요청 보내기
	err = sendPostRequest(clusterValue, finalYamlBase64, false)
	if err != nil {
		log.Printf("%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send POST request"})
		return
	}
	// POST 요청이 성공한 경우에만 요청 데이터를 DB에 저장
	_, err = db.Exec("INSERT INTO workload_info (workload_name, yaml, metadata, created_timestamp) VALUES (?, ?, ?, ?)",
		requestData.Metadata["name"], finalYamlBase64, string(metadataJSON), requestData.Timestamp)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// POST 요청 처리 함수
func handleSubmitResourceRequest(c *gin.Context) {
	var requestResourceData ys.RequestResourceData

	// 요청 데이터 바인딩
	if err := c.ShouldBindJSON(&requestResourceData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 1. name 값으로 가장 최신 created_timestamp 값을 갖는 yaml 찾기
	var base64Yaml string
	err := db.QueryRow(`
        SELECT yaml 
        FROM workload_info 
        WHERE workload_name = ? 
        ORDER BY created_timestamp DESC 
        LIMIT 1`, requestResourceData.Name).Scan(&base64Yaml)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Workload not found"})
		return
	}
	// 3. 기존 metadata 조회
	var existingMetadata string
	query := "SELECT metadata FROM workload_info WHERE workload_name = ? ORDER BY created_timestamp DESC LIMIT 1"
	err = db.QueryRow(query, requestResourceData.Name).Scan(&existingMetadata)
	if err != nil {
		if err == sql.ErrNoRows {
			existingMetadata = "{}" // metadata가 없을 경우 빈 JSON 객체로 설정
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	log.Print(base64Yaml) ///////////////////////
	// 2. YAML 데이터 디코딩 (Base64 -> YAML 문자열)
	yamlData, err := base64.StdEncoding.DecodeString(base64Yaml)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode YAML"})
		return
	}

	// 3. YAML 데이터를 Map 형식으로 변환하여 편집이 용이하게 함
	var yamlMap map[string]interface{}
	err = yaml.Unmarshal(yamlData, &yamlMap)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse YAML"})
		return
	}

	// // 3.1. 디버깅을 위한 로그 출력 (yamlMap 구조 확인)
	// log.Printf("Parsed YAML Map: %+v", yamlMap)

	// 4. JSON의 containers 값으로 YAML의 자원 값 갱신
	// yamlMap["spec"]을 map[interface{}]interface{}로 처리한 뒤, 이를 map[string]interface{}로 변환
	spec, ok := yamlMap["spec"].(map[interface{}]interface{})
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "'spec' not found or has incorrect type in YAML"})
		return
	}

	// "templates"를 map[interface{}]interface{}에서 []interface{}로 변환
	templatesInterface, ok := spec["templates"].([]interface{})
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "'templates' not found or has incorrect type in YAML"})
		return
	}

	for i, container := range requestResourceData.Containers {
		// 각 template의 container 자원 값 갱신
		template, ok := templatesInterface[i].(map[interface{}]interface{})
		if !ok {
			continue // 순서가 맞지 않으면 넘어감
		}

		containerMap, ok := template["container"].(map[interface{}]interface{})
		if !ok {
			continue // container 정보가 없으면 넘어감
		}

		resources, ok := containerMap["resources"].(map[interface{}]interface{})
		if !ok {
			continue // resources가 없으면 넘어감
		}

		// 요청된 자원 값 갱신
		resources["requests"] = map[string]string{
			"cpu":            container.Resources.Requests.CPU,
			"memory":         container.Resources.Requests.Memory,
			"gpu":            container.Resources.Requests.GPU,
			"nvidia.com/gpu": container.Resources.Requests.NvidiaGPU,
		}
		resources["limits"] = map[string]string{
			"cpu":            container.Resources.Limits.CPU,
			"memory":         container.Resources.Limits.Memory,
			"gpu":            container.Resources.Limits.GPU,
			"nvidia.com/gpu": container.Resources.Limits.NvidiaGPU,
		}

		// 요청 값이 비어 있으면 해당 항목 삭제
		for key, value := range resources["requests"].(map[string]string) {
			if value == "" {
				delete(resources["requests"].(map[string]string), key)
			}
		}

		// 제한 값이 비어 있으면 해당 항목 삭제
		for key, value := range resources["limits"].(map[string]string) {
			if value == "" {
				delete(resources["limits"].(map[string]string), key)
			}
		}
	}
	// 5. 수정된 YAML을 다시 인코딩하여 Base64로 변환
	modifiedYaml, err := yaml.Marshal(yamlMap)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal modified YAML"})
		return
	}
	finalYamlBase64 := base64.StdEncoding.EncodeToString(modifiedYaml)
	// 6. 수정된 YAML을 기반으로 POST 요청
	clusterValue := "1" // 필요시 적절히 할당
	err = sendPostRequest(clusterValue, finalYamlBase64, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send POST request"})
		return
	}
	// 9. 요청 데이터 DB에 저장
	loc, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load location"})
		return
	}
	_, err = db.Exec("INSERT INTO workload_info (workload_name, yaml, metadata, created_timestamp) VALUES (?, ?, ?, ?)",
		requestResourceData.Name, finalYamlBase64, existingMetadata, time.Now().In(loc).Format("2006-01-02 15:04:05"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

///////////////////////////////////////////////////////////////

func sendPostRequest(clusterValue string, finalYamlBase64 string, retryValue bool) error {
	wrapperIp := os.Getenv("WRAPPER_IP")
	wrapperPort := os.Getenv("WRAPPER_PORT")
	wrapperPath := os.Getenv("WRAPPER_PATH")
	address := "http://" + wrapperIp + ":" + wrapperPort + wrapperPath
	
	// 요청 데이터 생성
	wrapperData := ys.WrapperData{
		Cluster: clusterValue,
		Yaml:    finalYamlBase64,
		Retry:   retryValue,
	}

	postJSON, err := json.Marshal(wrapperData)
	if err != nil {
		return fmt.Errorf("failed to create JSON for POST request: %v", err)
	}

	// 요청 로그 출력
	log.Printf("Sending POST request to: %s", address)
	log.Printf("Request Body: %s", string(postJSON))

	// HTTP POST 요청 생성
	req, err := http.NewRequest("POST", address, bytes.NewBuffer(postJSON))
	if err != nil {
		return fmt.Errorf("failed to create POST request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// 요청 헤더 출력
	for key, value := range req.Header {
		log.Printf("Request Header: %s=%s", key, value)
	}

	// HTTP 클라이언트로 요청 전송
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send POST request: %v", err)
	}
	defer resp.Body.Close()

	// 응답 상태 코드 출력
	log.Printf("Response Status Code: %d", resp.StatusCode)

	// 응답 헤더 출력
	for key, value := range resp.Header {
		log.Printf("Response Header: %s=%s", key, value)
	}

	// 응답 바디 읽기
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}
	log.Printf("Response Body: %s", string(bodyBytes))

	// 상태 코드 확인
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}
	return nil
}

func ReqResourceAllocInfo(allocInfo ys.ResourceAllocInfo) ys.RespResource {
	var err error
	allocInfo.EncodedYaml = strings.TrimSpace(allocInfo.EncodedYaml)
	// Base64 디코딩
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

	// 요청할 리소스 JSON 객체 생성
	reqJson := ys.ReqResource{}
	uuid := "dmkim"
	currentTime := time.Now()
	nowTime := currentTime.Format("2006-01-02 15:04:05")

	// 기본 정보 설정
	reqJson.Version = "0.12"
	reqJson.Request.Name = workflow.Metadata.GenerateName
	reqJson.Request.ID = uuid
	reqJson.Request.Date = nowTime

	// task의 의존성 그래프 구축
	taskOrders := make(map[string]int)
	inDegree := make(map[string]int)
	dependencyGraph := make(map[string][]string)

	// 초기화: 의존성, in-degree, 그래프 초기화
	for _, template := range workflow.Spec.Templates {
		if template.DAG != nil {
			for _, task := range template.DAG.Tasks {
				taskOrders[task.Name] = 0
				inDegree[task.Name] = 0
				dependencyGraph[task.Name] = []string{}
			}
		}
	}

	// 의존성 그래프 및 in-degree 계산
	for _, template := range workflow.Spec.Templates {
		if template.DAG != nil {
			for _, task := range template.DAG.Tasks {
				for _, dep := range task.Dependencies {
					dependencyGraph[dep] = append(dependencyGraph[dep], task.Name)
					inDegree[task.Name]++
				}
			}
		}
	}

	// 위상 정렬(Topological Sorting) 수행
	queue := []string{}
	for task, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, task)
			taskOrders[task] = 1 // 의존성이 없는 task는 order 1
		}
	}

	// 위상 정렬을 이용해 의존성 처리
	for len(queue) > 0 {
		currentTask := queue[0]
		queue = queue[1:]

		// 현재 task를 의존하는 task들의 in-degree 감소 및 처리
		for _, dependentTask := range dependencyGraph[currentTask] {
			inDegree[dependentTask]--
			if inDegree[dependentTask] == 0 {
				// 부모 task의 order 값 + 1
				taskOrders[dependentTask] = taskOrders[currentTask] + 1
				queue = append(queue, dependentTask)
			}
		}
	}

	// 컨테이너 정보에 order 값 추가
	for _, value := range workflow.Spec.Templates {
		if value.Container == nil {
			continue
		} else {
			tmpContainer := ys.Container{
				Name: value.Name,
				Resources: ys.Resources{
					Requests: ys.ResourceDetails{
						CPU:              value.Container.Resources.Requests.CPU,
						GPU:              value.Container.Resources.Requests.GPU,
						Memory:           value.Container.Resources.Requests.Memory,
						EphemeralStorage: value.Container.Resources.Requests.EphemeralStorage,
					},
					Limits: ys.ResourceDetails{
						GPU:              value.Container.Resources.Limits.NvidiaGPU,
						CPU:              value.Container.Resources.Limits.CPU,
						Memory:           value.Container.Resources.Limits.Memory,
						EphemeralStorage: value.Container.Resources.Limits.EphemeralStorage,
					},
				},
			}
			tmpContainer.Attribute.MaxReplicas = 500
			tmpContainer.Attribute.TotalSize = 500
			tmpContainer.Attribute.PredictedExecutionTime = 600
			if order, exists := taskOrders[value.Name]; exists {
				tmpContainer.Attribute.Order = order
			}

			reqJson.Request.Containers = append(reqJson.Request.Containers, tmpContainer)
		}
	}

	// 기타 요청 속성 설정
	reqJson.Request.Attribute.WorkloadType = allocInfo.WorkloadTypeVal
	reqJson.Request.Attribute.IsCronJob = true
	reqJson.Request.Attribute.DevOpsType = "DEV"
	reqJson.Request.Attribute.GPUDriverVersion = 12.34
	reqJson.Request.Attribute.CudaVersion = 342.12
	if allocInfo.WorkloadFeatureVal != 0 {
		reqJson.Request.Attribute.WorkloadFeature = allocInfo.WorkloadFeatureVal
	}
	reqJson.Request.Attribute.UserID = uuid
	reqJson.Request.Attribute.Checkpoint = allocInfo.Checkpoint
	reqJson.Request.Attribute.Yaml = base64.StdEncoding.EncodeToString(data)

	// 리소스 요청 전송
	var ackBody ys.RespResource
	ack, body := SEND_REST_DATA(allocInfo.Addr, reqJson)
	log.Printf("[ReqResource] Status code: %d, Response body: %s", ack.StatusCode, body)

	if ack.StatusCode == http.StatusOK {
		var jsonResponse map[string]interface{}
		err = json.Unmarshal([]byte(body), &jsonResponse)
		if err != nil {
			log.Printf("Failed to unmarshal ack body to JSON: %s", err)
			return ackBody
		}

		yamlData, err := yaml.Marshal(jsonResponse)
		if err != nil {
			log.Printf("Failed to marshal JSON to YAML: %s", err)
			return ackBody
		}

		log.Printf("Converted YAML:\n%s", string(yamlData))
		err = yaml.Unmarshal(yamlData, &ackBody)
		if err != nil {
			log.Printf("Failed to unmarshal YAML data to ackBody: %s", err)
		}
	} else {
		fmt.Printf("[ReqResource] Request failed with status: %s\n", ack.Status)
	}
	return ackBody
}

func MadeFinalWorkloadYAML(argBody ys.RespResource, inputYaml string) (map[string]interface{}, string) {
	// Base64로 인코딩된 YAML을 디코딩
	clusterValue := argBody.Response.Cluster
	yamlFile, err := base64.StdEncoding.DecodeString(inputYaml)
	if err != nil {
		log.Fatalf("Error decoding Base64 YAML data: %v", err)
	}
	// YAML 데이터를 저장할 변수
	var data map[string]interface{}
	// YAML 데이터 언마샬링
	err = yaml.Unmarshal(yamlFile, &data)
	if err != nil {
		log.Fatalf("Error unmarshalling YAML data: %v", err)
	}
	// templates 섹션에서 모든 container의 image 값을 출력하고 조건에 따라 새로운 키를 추가
	spec, ok := data["spec"].(map[interface{}]interface{})
	if ok {
		templates, ok := spec["templates"].([]interface{})
		if ok {

			for _, template := range templates {
				templateMap, ok := template.(map[interface{}]interface{})
				if ok {

					for _, val := range argBody.Response.Containers {
						if templateMap["name"] == val.Name {
							templateMap["nodeSelector"] = ys.NodeSelect{
								Node: val.Node,
							}
						}
					}

				}
			}
		}
	}
	return data, clusterValue
}

