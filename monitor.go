package main

import (
	"encoding/json"
	"fmt"
	"github.com/win30221/HonestBeeHomeTest/server"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	var status server.Status
	limiter := time.Tick(time.Second)
	for {
		<-limiter
		req, err := http.NewRequest("GET", "http://127.0.0.1:8080/status", nil)
		if err != nil {
			fmt.Println(err.Error())
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println("監控服務不可用")
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}
		if err := json.Unmarshal([]byte(body), &status); err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		log.Println("=============================")
		log.Println("目前使用者:", status.CurrentConnectionCount)
		log.Println("目前請求率:", status.CurrentRequestRate * 100, "%")
		log.Println("已請求命令:", status.ProcessedRequestCount)
		log.Println("尚未完成命令:", status.RemaingJobs)
		log.Println("請求成功率:", status.RequestSuccessRate * 100, "%")

	}

}