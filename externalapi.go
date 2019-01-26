package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const (
	EXTERNAL_API_PORT	= "8081"
)

type Response struct {
	Delay	float64	`json:"delay"`
	Message	string	`json:"message"`
}

func main() {
	rand.Seed(time.Now().UnixNano())
	http.HandleFunc("/delay", DelayResponse)
	addr := "0.0.0.0:" + EXTERNAL_API_PORT
	log.Printf("開始監聽HTTP Server[%s]端口", EXTERNAL_API_PORT)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}

func DelayResponse(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query();
	message := vars["message"][0]
	log.Println("新請求訊息", message)

	t := rand.Float64() * 2 + 1
	time.Sleep(time.Second * time.Duration(t))
	res := &Response {
		Delay: t,
		Message: message,
	}

	b, err := json.Marshal(res)
	if err != nil {
		fmt.Println("Json處理異常:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Json處理異常"))
		return
	}

	fmt.Fprintf(w, "%s", b)
}