package main

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"net/url"
	"strconv"
	"time"
)

// FoxStatus : get today's status
func FoxStatus(api *anaconda.TwitterApi) {
	const screenName = "Arthur_Lugh"
	const baseURL = "https://api.twitter.com/1.1/statuses/user_timeline.json"
	const dateFormat = "2006-01-02"
	const parseDateFormat = "Mon Jan 2 15:04:05 -0700 2006"

	values := url.Values{}
	values.Add("screen_name", screenName)
	response, err := api.GetUserTimeline(values)
	if err != nil {
		panic(err)
	}

	// get yesterday timestamp (-1h)
	time.Local = time.FixedZone("JST", 9*60*60)
	day := time.Now()
	day = day.Add(-time.Duration(1) * time.Hour)
	fmt.Println(day)

	count := 0
	jst, _ := time.LoadLocation("Asia/Tokyo")
	// check timestamp
	for i := 0; i < 5; i++ {
		timestamp := response[i].CreatedAt
		fmt.Println(timestamp)
		timeObj, _ := time.ParseInLocation(parseDateFormat, timestamp, jst)
		fmt.Println(timeObj)
	}
	tweetText := day.Format(dateFormat) + "\nTweet Count:" + strconv.Itoa(count)
	fmt.Println(tweetText)
	return
}
