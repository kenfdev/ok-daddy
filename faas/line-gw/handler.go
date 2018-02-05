package function

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"

	"github.com/line/line-bot-sdk-go/linebot"
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

func findDaddy(body []byte, sign string, secret string, accessToken string) {

	req, _ := http.NewRequest("POST", "", bytes.NewReader(body))
	req.Header.Set("X-Line-Signature", sign)

	bot, err := linebot.New(secret, accessToken)

	events, err := bot.ParseRequest(req)
	if err != nil {
		// Do something when something bad happened.
	}

	messages := []linebot.Message{linebot.NewTextMessage("Hello, world")}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			// Do Something...
			_, err := bot.ReplyMessage(event.ReplyToken, messages...).Do()
			if err != nil {
				// Do something when some bad happened
			}
		}
	}
}
