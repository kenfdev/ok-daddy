package function

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

var (
	whereIsDaddyUrl = os.Getenv("FIND_DADDY_URL")
	sendMessageUrl  = os.Getenv("SEND_MESSAGE_URL")
)

type Request struct {
	QueryResult struct {
		Parameters map[string]string `json:"parameters"`
		Intent     struct {
			DisplayName string `json:"displayName"`
		} `json:"intent"`
	} `json:"queryResult"`
}

type Response struct {
	FulfillmentText string `json:"fulfillmentText"`
	Source          string `json:"source"`
}

type OKDaddyRequest struct {
	IntentName string            `json:"intentName"`
	Parameters map[string]string `json:"parameters"`
}

type OKDaddyResponse struct {
	Result string `json:"result"`
}

// Handle a serverless request
func Handle(r []byte) string {
	var request Request
	json.Unmarshal(r, &request)

	req := &OKDaddyRequest{
		IntentName: request.QueryResult.Intent.DisplayName,
		Parameters: request.QueryResult.Parameters,
	}

	jsonValue, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}

	var url string
	switch req.IntentName {
	case "Where Intent":
		url = whereIsDaddyUrl
	case "Send Message Intent":
		url = sendMessageUrl
	}

	resp, _ := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))

	var okDaddyResponse OKDaddyResponse
	err = json.NewDecoder(resp.Body).Decode(&okDaddyResponse)
	if err != nil {
		panic(err)
	}

	response := &Response{
		FulfillmentText: okDaddyResponse.Result,
		Source:          "ok-daddy",
	}

	j, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s", string(j))
}
