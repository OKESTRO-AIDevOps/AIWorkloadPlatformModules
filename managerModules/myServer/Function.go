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

// [main, handleGetWorkloadinfoRequest, handleGetStratoRequest, handleSubmitRequest 유지]

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

func sendPostRequest(clusterValue string, finalYamlBase64 string, retryValue bool) error {
	return nil // 임시 반환
}
