package gopublica

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
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
	ErrServiceUnavailable  = errors.New("503 Service Unavailable - We’re temporarially offline for maintenance. Please try again later.")
	ErrUnkown              = errors.New("Unkown error")
)

var (
	apiKey   = getApiKey()
	baseUrl  = "https://api.propublica.org/congress/v1/"
	fileType = ".json"
)

type Result struct {
	Status    string `json:"status"`
	Copyright string `json:"copyright"`
}

func SetAPIKey(key string) {
	apiKey = key
}

// get will create an http request with the GET method
func get(parts ...string) (*http.Response, error) {
	return do(http.NewRequest("GET", makeUrl(parts), nil))
}

// do will set the api key header required by propublica
// api and execute the http request. it will check for
// erroneous status codes, however the propublica api
// is pretty inconsistent with status.
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

// makeUrl will take a slice of strings, join them
// with a slash ("/") and build an api request url
func makeUrl(parts []string) string {
	url := baseUrl + strings.Join(parts, "/") + fileType
	return url
}

// statusOk will return true only if the
// status code passed is in the range [200, 299]
func statusOk(code int) bool {
	return code >= 200 && code < 300
}

// getError will return the appropriate error
// object for the code
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

// getApiKey will first attempt to look for the api key
// in the env variable PROPUBLICA_API_KEY. if it is not
// found then it will attempt to read the file
// .propublica-api-key in the home directory of the
// current user
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
