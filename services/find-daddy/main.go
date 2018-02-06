package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var (
	fetchCheckInUrl = os.Getenv("CHECKIN_URL")
)

// our main function
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/find-daddy", FindDaddy).Methods("POST")
	fmt.Println("Listening at :8000...")
	log.Fatal(http.ListenAndServe(":8000", router))
}

type FindDaddyRequest struct {
	IntentName string `json:"intentName"`
	Parameters struct {
		Someone string `json:"Someone"`
	} `json:"parameters"`
}

type FindDaddyResponse struct {
	Result string `json:"result"`
}

type CheckInResponse struct {
	Name string `json:"name"`
}

func FindDaddy(w http.ResponseWriter, r *http.Request) {
	var req FindDaddyRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	if req.IntentName != "Where Intent" || req.Parameters.Someone != "Dad" {
		respondWithJSON(w, 200, &FindDaddyResponse{Result: "I'm afraid I don't have an answer for that."})
		return
	}

	resp, err := http.Post(fetchCheckInUrl, "application/json", bytes.NewBuffer([]byte("")))
	if err != nil {
		respondWithJSON(w, 200, &FindDaddyResponse{Result: "Sorry, I don't know where I am either. Maybe I can answer later."})
		return
	}

	var checkin CheckInResponse
	_ = json.NewDecoder(resp.Body).Decode(&checkin)

	respondWithJSON(w, 200, &FindDaddyResponse{Result: "I'm probably at " + checkin.Name})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
