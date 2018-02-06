package function

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestFind(t *testing.T) {

	expectedVenueName := "Some Place"
	expectedToken := "token"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		values, err := url.ParseQuery(r.URL.String())
		if err != nil {
			t.Fail()
		}

		actualToken := values.Get("oauth_token")
		if expectedToken != actualToken {
			t.Errorf("expected %s but got %s", expectedToken, actualToken)
			t.Fail()
		}

		w.Header().Set("Content-Type", "application/json")
		resp := `{
			"response": {
				"checkins": {
					"items": [
						{
							"venue": {
								"name": "` + expectedVenueName + `"
							}
						}
					]
				}
			}
		}`
		fmt.Fprintln(w, resp)
	}))
	defer ts.Close()

	fsUrl = ts.URL

	j, err := find(expectedToken)
	if err != nil {
		t.Fail()
	}

	if expectedVenueName != j.Response.Checkins.Items[0].Venue.Name {
		t.Errorf("expected %s but got %s", expectedVenueName, j.Response.Checkins.Items[0].Venue.Name)
		t.Fail()
	}

}
