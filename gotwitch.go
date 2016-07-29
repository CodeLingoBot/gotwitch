package gotwitch

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

type errorResponse struct {
	Error   string `json:"error"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// TwitchAPI struct
type TwitchAPI struct {
	ClientID string
}

// SuccessCallback runs on a successfull request and parse
type SuccessCallback func()

// HTTPErrorCallback runs on a errored HTTP request
type HTTPErrorCallback func(statusCode int, statusMessage, errorMessage string)

// InternalErrorCallback runs on an internal error
type InternalErrorCallback func(error)

// New instantiates a new TwitchAPI object
func New(clientID string) *TwitchAPI {
	return &TwitchAPI{
		ClientID: clientID,
	}
}

var client = &http.Client{}

func (twitchAPI *TwitchAPI) request(verb, baseURL string, parameters url.Values, requestBody interface{}, responseBody interface{},
	onSuccess SuccessCallback, onHTTPError HTTPErrorCallback, onInternalError InternalErrorCallback) {
	url := "https://api.twitch.tv/kraken" + baseURL + "?" + parameters.Encode()
	var request *http.Request
	var err error
	if requestBody != nil {
		serializedRequestBody, err := json.Marshal(requestBody)
		if err != nil {
			onInternalError(err)
			return
		}

		serializedRequestBodyReader := bytes.NewReader(serializedRequestBody)
		request, err = http.NewRequest(verb, url, serializedRequestBodyReader)
	} else {
		request, err = http.NewRequest(verb, url, nil)
	}
	if err != nil {
		onInternalError(err)
		return
	}

	twitchAPI.setHeaders(request)
	response, err := client.Do(request)
	if err != nil {
		onInternalError(err)
		return
	}

	if response.StatusCode >= 300 {
		handleHTTPError(response, onHTTPError, onInternalError)
		return
	}

	handleSuccess(response, responseBody, onSuccess, onInternalError)
}

// Get request
func (twitchAPI *TwitchAPI) Get(baseURL string, parameters url.Values, responseBody interface{}, onSuccess SuccessCallback,
	onHTTPError HTTPErrorCallback, onInternalError InternalErrorCallback) {
	twitchAPI.request("GET", baseURL, parameters, nil, responseBody, onSuccess, onHTTPError, onInternalError)
}

// Put request
func (twitchAPI *TwitchAPI) Put(baseURL string, parameters url.Values, requestBody interface{}, responseBody interface{}, onSuccess SuccessCallback,
	onHTTPError HTTPErrorCallback, onInternalError InternalErrorCallback) {
	twitchAPI.request("PUT", baseURL, parameters, requestBody, responseBody, onSuccess, onHTTPError, onInternalError)
}

// Post request
func (twitchAPI *TwitchAPI) Post(baseURL string, parameters url.Values, requestBody interface{}, responseBody interface{}, onSuccess SuccessCallback,
	onHTTPError HTTPErrorCallback, onInternalError InternalErrorCallback) {
	twitchAPI.request("POST", baseURL, parameters, requestBody, responseBody, onSuccess, onHTTPError, onInternalError)
}

func (twitchAPI *TwitchAPI) setHeaders(request *http.Request) {
	request.Header.Add("Client-ID", twitchAPI.ClientID)
	request.Header.Add("Accept", "application/vnd.twitchtv.v3+json")
}

// Delete request
func (twitchAPI *TwitchAPI) Delete(baseURL string, parameters url.Values, responseBody interface{}, onSuccess SuccessCallback,
	onHTTPError HTTPErrorCallback, onInternalError InternalErrorCallback) {
	twitchAPI.request("DELETE", baseURL, parameters, nil, responseBody, onSuccess, onHTTPError, onInternalError)
}

func handleSuccess(response *http.Response, data interface{}, onSuccess SuccessCallback, onInternalError InternalErrorCallback) {
	body, err := body(response)
	if err != nil {
		onInternalError(err)
		return
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		onInternalError(err)
		return
	}

	onSuccess()
}

func handleHTTPError(response *http.Response, onHTTPError HTTPErrorCallback, onInternalError InternalErrorCallback) {
	body, err := body(response)
	if err != nil {
		onInternalError(err)
		return
	}

	var errorResponse errorResponse
	err = json.Unmarshal(body, &errorResponse)
	if err != nil {
		onInternalError(err)
		return

	}

	onHTTPError(errorResponse.Status, errorResponse.Message, errorResponse.Error)
}

func body(response *http.Response) ([]byte, error) {
	body, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	return body, err
}
