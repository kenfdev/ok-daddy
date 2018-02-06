package function

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	fsUrl    = os.Getenv("FS_URL")
	fsApiVer = os.Getenv("FS_API_VER")
)

type Location struct {
	Address          string   `json:"address"`
	Lat              float64  `json:"lat"`
	Lng              float64  `json:"lng"`
	FormattedAddress []string `json:"formattedAddress"`
}

type CheckInItem struct {
	Venue struct {
		Name     string   `json:"name"`
		Location Location `json:"location"`
		URL      string   `json:"url"`
	} `json:"venue"`
}

type FourSquareResponse struct {
	Response struct {
		Checkins struct {
			Items []CheckInItem `json:"items"`
		} `json:"checkins"`
	} `json:"response"`
}

type Response struct {
	Name     string   `json:"name"`
	Location Location `json:"location"`
	URL      string   `json:"url"`
}

// Handle a serverless request
func Handle(r []byte) string {
	oauthToken, _ := ioutil.ReadFile("/run/secrets/" + os.Getenv("OAUTH_TOKEN_KEY_NAME"))

	fsResp, err := find(string(oauthToken))
	if err != nil {
		panic(err)
	}

	venue := fsResp.Response.Checkins.Items[0].Venue
	response := &Response{
		Name:     venue.Name,
		Location: venue.Location,
		URL:      venue.URL,
	}

	j, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%s", string(j))
}

func find(token string) (*FourSquareResponse, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", fsUrl, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("oauth_token", string(token))
	q.Add("v", fsApiVer)
	q.Add("limit", "1")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var fsResp FourSquareResponse
	err = json.NewDecoder(resp.Body).Decode(&fsResp)
	if err != nil {
		return nil, err
	}

	return &fsResp, nil
}
