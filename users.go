package gotwitch

import (
	"net/http"
	"net/url"

	"github.com/pajlada/jsonapi"
)

// User json to struct
type User struct {
	ID              string `json:"id"`
	Login           string `json:"login"`
	DisplayName     string `json:"display_name"`
	Type            string `json:"type"`
	BroadcasterType string `json:"broadcaster_type"`
	Description     string `json:"description"`
	ProfileImageURL string `json:"profile_image_url"`
	OfflineImageURL string `json:"offline_image_url"`
	ViewCount       int    `json:"view_count"`
	Email           string `json:"email"`
}

// FollowedStreams json to struct
type FollowedStreams struct {
	Streams []Stream `json:"streams"`
	Total   int      `json:"_total"`
}

type usersListChannel struct {
	Data []User `json:"data"`
}

// GetFollowedStreams request for GET https://api.twitch.tv/kraken/streams/followed
func (twitchAPI *TwitchAPI) GetFollowedStreams(oauthToken string, parameters url.Values,
	onSuccess func(FollowedStreams), onHTTPError jsonapi.HTTPErrorCallback,
	onInternalError jsonapi.InternalErrorCallback) {
	var followedStreams FollowedStreams
	onSuccessfulRequest := func() {
		onSuccess(followedStreams)
	}
	twitchAPI.AuthenticatedGet("/streams/followed", parameters, oauthToken, &followedStreams,
		onSuccessfulRequest, onHTTPError, onInternalError)
}

// GetUsers request for GET https://api.twitch.tv/helix/users
func (twitchAPI *TwitchAPI) GetUsers(userIDs []string,
	onSuccess func([]User),
	onHTTPError jsonapi.HTTPErrorCallback,
	onInternalError jsonapi.InternalErrorCallback) (response *http.Response, err error) {
	var usersListChannel usersListChannel
	onSuccessfulRequest := func() {
		onSuccess(usersListChannel.Data)
	}
	parameters := make(url.Values)
	for _, userID := range userIDs {
		parameters.Add("id", userID)
	}
	return twitchAPI.authenticatedAPI.Get("/users", parameters, &usersListChannel,
		onSuccessfulRequest, onHTTPError, onInternalError)
}

// GetUsersSimple request for GET https://api.twitch.tv/helix/users
func (twitchAPI *TwitchAPI) GetUsersSimple(userIDs []string) (data []User, response *http.Response, err error) {
	var errorChannel = make(chan error)
	onSuccessfulRequest := func(users []User) {
		data = users
		errorChannel <- nil
	}
	go func() {
		response, err = twitchAPI.GetUsers(userIDs, onSuccessfulRequest, simpleOnHTTPError(errorChannel), simpleOnInternalError(errorChannel))
	}()

	err = <-errorChannel

	return
}

// GetUsersByLogin request for GET https://api.twitch.tv/helix/users
func (twitchAPI *TwitchAPI) GetUsersByLogin(userLogins []string,
	onSuccess func([]User),
	onHTTPError jsonapi.HTTPErrorCallback,
	onInternalError jsonapi.InternalErrorCallback) (response *http.Response, err error) {
	var usersListChannel usersListChannel
	onSuccessfulRequest := func() {
		onSuccess(usersListChannel.Data)
	}
	parameters := make(url.Values)
	for _, userLogin := range userLogins {
		parameters.Add("login", userLogin)
	}
	return twitchAPI.get("/users", parameters, &usersListChannel,
		onSuccessfulRequest, onHTTPError, onInternalError)
}

// GetUsersByLoginSimple request for GET https://api.twitch.tv/helix/users
func (twitchAPI *TwitchAPI) GetUsersByLoginSimple(userLogins []string) (data []User, response *http.Response, err error) {
	var errorChannel = make(chan error)
	onSuccessfulRequest := func(users []User) {
		data = users
		errorChannel <- nil
	}
	go func() {
		response, err = twitchAPI.GetUsersByLogin(userLogins, onSuccessfulRequest, simpleOnHTTPError(errorChannel), simpleOnInternalError(errorChannel))
	}()

	err = <-errorChannel

	return
}
