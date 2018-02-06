package function

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/line/line-bot-sdk-go/linebot"
)

var (
	brokerUrl = os.Getenv("BROKER_URL")
)

// Handle a serverless request
func Handle(req []byte, wg *sync.WaitGroup) string {
	lineSignature := os.Getenv("Http_X_Line_Signature")
	secret, _ := ioutil.ReadFile("/run/secrets/" + os.Getenv("CHANNEL_SECRET_KEY_NAME"))
	accessToken, _ := ioutil.ReadFile("/run/secrets/" + os.Getenv("CHANNEL_ACCESS_TOKEN_KEY_NAME"))

	wg.Add(1)
	go func() {

		findDaddy(req, lineSignature, string(secret), string(accessToken))
		defer wg.Done()

	}()
	return fmt.Sprintln("done")
}

type DFBrokerRequest struct {
	Text string `json:"text"`
}

type DFBrokerResponse struct {
	Metadata struct {
		IntentName string `json:"intentName"`
	}
	FulFillment struct {
		Speech string `json:"speech"`
	} `json:"fulfillment"`
}

func findDaddy(body []byte, sign string, secret string, accessToken string) {

	req, _ := http.NewRequest("POST", "", bytes.NewReader(body))
	req.Header.Set("X-Line-Signature", sign)

	bot, err := linebot.New(secret, accessToken)

	events, err := bot.ParseRequest(req)
	if err != nil {
		// Do something when something bad happened.
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				bReq := &DFBrokerRequest{
					Text: message.Text,
				}

				jsonValue, _ := json.Marshal(bReq)

				bResp, _ := http.Post(brokerUrl, "application/json", bytes.NewBuffer(jsonValue))

				var dfBrokerResp DFBrokerResponse
				_ = json.NewDecoder(bResp.Body).Decode(&dfBrokerResp)

				var reply = dfBrokerResp.FulFillment.Speech

				// TODO: Don't hard code this.
				if dfBrokerResp.Metadata.IntentName == "Send Message Request Intent" {
					reply = "Please send the message by yourself. You can do it directly with LINE."
				}

				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(reply)).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}
}
