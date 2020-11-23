package main

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"net/url"
)

// FoxStatus : get today's status
func FoxStatus(api *anaconda.TwitterApi) {
	const screenName = "Arthur_Lugh"
	const baseURL = "https://api.twitter.com/1.1/statuses/user_timeline.json"

	values := url.Values{}
	values.Add("screen_name", screenName)
	response, err := api.GetUserTimeline(values)
	if err != nil {
		panic(err)
	}
	fmt.Println(response)
	return
}
