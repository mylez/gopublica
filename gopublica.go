package gopublica

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"strconv"
	"strings"
)

var (
	ErrBadRequest          = errors.New("400 Bad Request - Your request sucks")
	ErrUnauthorized        = errors.New("401 Unauthorized - Your API key is wrong")
	ErrForbidden           = errors.New("403 Forbidden - The kitten requested is hidden for administrators only")
	ErrNotFound            = errors.New("404 Not Found - The specified kitten could not be found")
	ErrMethodNotAllowed    = errors.New("405 Method Not Allowed - You tried to access a kitten with an invalid method")
	ErrNotAcceptable       = errors.New("406 Not Acceptable - You requested a format that isn’t json")
	ErrGone                = errors.New("410 Gone - The kitten requested has been removed from our servers")
	ErrTeapot              = errors.New("418 I’m a teapot")
	ErrTooManyRequests     = errors.New("429 Too Many Requests - You’re requesting too many kittens! Slow down!")
	ErrInternalServerError = errors.New("500 Internal Server Error - We had a problem with our server. Try again later.")
	ErrServiceUnavailable  = errors.New("503 Service Unavailable - We’re temporarially offline for maintanance. Please try again later.")
	ErrUnkown              = errors.New("Unkown error")
)

var (
	apiKey   = getApiKey()
	baseUrl  = "https://api.propublica.org/congress/v1/"
	fileType = ".json"
)

type Result struct {
	Status string `json:"status"`
	//Copyright string `json:"copyright"`
}

type CosponsorsResult struct {
	Result
	Results []struct {
		Cosponsors []struct {
			CosponsorId string `json:"cosponsor_id"`
			Name        string `json:"name"`
			Date        string `json:"date"`
		} `json:"cosponsors"`
		Congress              string `json:"congress"`
		Number                string `json:"number"`
		BillUri               string `json:"bill_uri"`
		Title                 string `json:"title"`
		SponsorId             string `json:"sponsor_id"`
		IntroducedDate        string `json:"introduced_date"`
		Committees            string `json:"committees"`
		LatestMajorActionDate string `json:"latest_major_action_date"`
		LatestMajorAction     string `json:"latest_major_action"`
	} `json:"results"`
}

// :congress/bills/:bill-id/cosponsors.js
//
func GetCosponsorsForBill(congress, billId string) (*CosponsorsResult, error) {
	rsp, err := get(congress, "bills", billId, "cosponsors")

	if err != nil {
		return nil, err
	}

	bod, _ := ioutil.ReadAll(rsp.Body)
	res := &CosponsorsResult{}
	json.Unmarshal(bod, res)

	if err != nil {
		return nil, err
	}

	code, _ := strconv.Atoi(res.Result.Status)

	if !statusOk(code) && res.Status != "OK" {
		return nil, getError(code)
	}

	return res, nil
}

func get(parts ...string) (*http.Response, error) {
	return do(http.NewRequest("GET", makeUrl(parts), nil))
}

func do(req *http.Request, err error) (*http.Response, error) {
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-API-Key", apiKey)

	rsp, err := http.DefaultClient.Do(req)

	if !statusOk(rsp.StatusCode) {
		return nil, getError(rsp.StatusCode)
	}

	return rsp, err
}

func makeUrl(parts []string) string {
	url := baseUrl + strings.Join(parts, "/") + fileType
	return url
}

func statusOk(code int) bool {
	return code >= 200 && code < 300
}

func getError(code int) error {
	switch code {
	case 400:
		return ErrBadRequest
	case 401:
		return ErrUnauthorized
	case 403:
		return ErrForbidden
	case 404:
		return ErrNotFound
	case 405:
		return ErrMethodNotAllowed
	case 406:
		return ErrNotAcceptable
	case 410:
		return ErrGone
	case 418:
		return ErrTeapot
	case 429:
		return ErrTooManyRequests
	case 500:
		return ErrInternalServerError
	case 503:
		return ErrServiceUnavailable
	}
	return ErrUnkown
}

func getApiKey() string {
	key := os.Getenv("PROPUBLICA_API_KEY")
	if key != "" {
		return key
	}
	usr, _ := user.Current()
	if usr == nil {
		return ""
	}
	bytes, _ := ioutil.ReadFile(usr.HomeDir + "/.propublica-api-key")
	return string(bytes)
}
