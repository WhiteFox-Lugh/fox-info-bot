package main

import (
	"github.com/ChimeraCoder/anaconda"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// FoxStatus : get today's status
func FoxStatus(api *anaconda.TwitterApi) {
	const screenName = "Arthur_Lugh"
	const myID = "3145003784"
	const baseURL = "https://api.twitter.com/1.1/statuses/user_timeline.json"
	const dateFormat = "2006-01-02"
	const parseDateFormat = "Mon Jan 2 15:04:05 -0700 2006"
	const getTweetCount = "200"
	const durationDay = 24 * 60 * 60
	jst, _ := time.LoadLocation("Asia/Tokyo")

	values := url.Values{}
	values.Add("screen_name", screenName)
	values.Add("count", getTweetCount)
	response, err := api.GetUserTimeline(values)
	if err != nil {
		panic(err)
	}

	// get timestamp
	time.Local = time.FixedZone("JST", 9*60*60)
	day := time.Now()
	timeStr := day.Format(dateFormat)
	currentDayObj, _ := time.ParseInLocation(dateFormat, timeStr, jst)

	count := 0

	// check timestamp
	for i := 0; i < len(response); i++ {
		timestamp := response[i].CreatedAt
		idStr := response[i].User.IdStr
		source := response[i].Source

		if strings.Contains(source, "狐日和") || idStr != myID {
			continue
		}

		timeObj, _ := time.Parse(parseDateFormat, timestamp)
		timeObj = timeObj.In(jst)

		// current time - tweet time
		duration := currentDayObj.Sub(timeObj)
		durationSec := int(duration.Seconds())

		if 0 <= durationSec && durationSec < durationDay {
			count++
		}
	}
	tweetText := "(っ ॑꒳ ॑c).+(" + day.Format(dateFormat) + " Report)\n前日ツイート数（bot 除外）:" + strconv.Itoa(count)
	_, postErr := api.PostTweet(tweetText+"\n(bot)", nil)

	if postErr != nil {
		panic(postErr)
	}
	return
}
