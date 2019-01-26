package server

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
)

type Status struct {
	CurrentConnectionCount	int		`json:"current_connection_count"`
	CurrentRequestRate		float64	`json:"current_request_rate"`
	ProcessedRequestCount	int		`json:"processed_request_count"`
	RemaingJobs				int		`json:"remaing_jobs"`
	RequestSuccessRate		float64	`json:"request_success_rate"`
}

func StartHttpServer() {
	http.HandleFunc("/status", GetStatus)
	addr := "0.0.0.0:" + HTTP_PORT
	log.Printf("開始監聽HTTP Server[%s] Port", HTTP_PORT)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}

func GetStatus(w http.ResponseWriter, r *http.Request) {

	ServerInstance.mu.RLock()
	defer ServerInstance.mu.RUnlock()

	var currentConnectionCount	= len(ServerInstance.clients)
	var currentRequestRate		= changeToPercent(float64(ServerInstance.currentRequestInSecond) / float64(REQUEST_PER_SECOND))
	var processedRequestCount	= ServerInstance.processedRequestCount
	var remaingJobs				= len(ServerInstance.insniffer)
	var requestSuccessRate		= changeToPercent(float64(ServerInstance.currentSuccessRequest) / float64(ServerInstance.processedRequestCount))

	if ServerInstance.processedRequestCount == 0 {
		requestSuccessRate = 0
	}

	status := &Status {
		CurrentConnectionCount: currentConnectionCount,
		CurrentRequestRate:		currentRequestRate,
		ProcessedRequestCount:  processedRequestCount,
		RemaingJobs:            remaingJobs,
		RequestSuccessRate:     requestSuccessRate,
	}

	b, err := json.Marshal(status)
	if err != nil {
		fmt.Println("Json處理異常:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Json處理異常"))
		return
	}

	fmt.Fprintf(w, "%s", b)
}

func changeToPercent(v float64) float64 {
	return math.Round(v * 10000) / 10000
}