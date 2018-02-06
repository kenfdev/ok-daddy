package function

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	expectedFulfillment := "Expected Result"
	expectedSomeone := "Dad"
	expectedIntent := "Where Intent"

	r := `{
		"queryResult": {
			"parameters": {
				"Someone": "` + expectedSomeone + `"
			},
			"intent": {
				"displayName": "` + expectedIntent + `"
			}
		}
	}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		bs, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		var fdReq OKDaddyRequest
		json.Unmarshal(bs, &fdReq)

		if expectedIntent != fdReq.IntentName {
			t.Errorf("expected %s but got %s", expectedIntent, fdReq.IntentName)
			t.Fail()
		}

		if expectedSomeone != fdReq.Parameters["Someone"] {
			t.Errorf("expected %s but got %s", expectedSomeone, fdReq.Parameters["Someone"])
			t.Fail()
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"result":"`+expectedFulfillment+`"}`)
	}))
	defer ts.Close()

	whereIsDaddyUrl = ts.URL

	// Act
	result := Handle([]byte(r))

	var actualResult Response
	json.Unmarshal([]byte(result), &actualResult)

	// Assert
	if expectedFulfillment != actualResult.FulfillmentText {
		t.Errorf("expected %s but got %s", expectedFulfillment, actualResult.FulfillmentText)
		t.Fail()
	}

	expectedSource := "ok-daddy"
	if expectedSource != actualResult.Source {
		t.Errorf("expected %s but got %s", expectedSource, actualResult.Source)
		t.Fail()
	}

}
