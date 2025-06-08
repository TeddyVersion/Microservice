package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type MiniApp struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type APIResponse struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

func listMiniAppsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Method not allowed"})
		return
	}
	// TODO: Replace with real mini-app discovery
	miniapps := []MiniApp{{ID: "m1", Name: "Taxi", URL: "https://taxi.example.com"}}
	json.NewEncoder(w).Encode(APIResponse{Status: "success", Data: miniapps})
}

func main() {
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})
	http.HandleFunc("/miniapps", listMiniAppsHandler)
	fmt.Println("miniapp-service running on :8090")
	http.ListenAndServe(":8090", nil)
}
