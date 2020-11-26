package main

import (
	"encoding/json"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "2006-01-02"

// EventData : github の Event のデータをもつ
type EventData struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	CreatedAt string `json:"created_at"`
}

// countTweetNum : 一日のツイート数を取得
func countTweetNum(api *anaconda.TwitterApi) (count int, day time.Time) {
	const screenName = "Arthur_Lugh"
	const myID = "3145003784"
	const baseURL = "https://api.twitter.com/1.1/statuses/user_timeline.json"
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
	day = time.Now()
	timeStr := day.Format(dateFormat)
	currentDayObj, _ := time.ParseInLocation(dateFormat, timeStr, jst)

	count = 0

	// check timestamp
	for i := 0; i < len(response); i++ {
		timestamp := response[i].CreatedAt
		idStr := response[i].User.IdStr
		source := response[i].Source

		if strings.Contains(source, "狐日和") || idStr != myID {
			continue
		}

		// parse して jst に変換
		timeObj, _ := time.Parse(parseDateFormat, timestamp)
		timeObj = timeObj.In(jst)

		// current time - tweet time
		duration := currentDayObj.Sub(timeObj)
		durationSec := int(duration.Seconds())

		if 0 <= durationSec && durationSec < durationDay {
			count++
		}
	}
	return
}

// getGithubCommitcount : 1日の GitHub event 数を取得
func getGithubCommitCount(api *anaconda.TwitterApi) (commitCount int) {
	const endpoint = "https://api.github.com/users/WhiteFox-Lugh/events"
	const githubDateFormat = "2006-01-02T15:04:05Z"
	const maxPage = 1
	const durationDay = 24 * 60 * 60
	var jsonData []EventData

	jst, _ := time.LoadLocation("Asia/Tokyo")

	// get timestamp (日付の境界)
	time.Local = time.FixedZone("JST", 9*60*60)
	day := time.Now()
	timeStr := day.Format(dateFormat)
	currentDayObj, _ := time.ParseInLocation(dateFormat, timeStr, jst)
	fmt.Println(currentDayObj)

	for i := 1; i <= maxPage; i++ {
		// endpoint 生成
		values := url.Values{}
		values.Add("page", strconv.Itoa(i))

		//リクエストの送信
		request, _ := http.NewRequest("GET", endpoint+"?"+values.Encode(), nil)
		response, err := http.DefaultClient.Do(request)
		if err != nil {
			panic(err.Error())
		}

		defer response.Body.Close()

		// parse json
		byteArray, err := ioutil.ReadAll(response.Body)
		if err := json.Unmarshal(byteArray, &jsonData); err != nil {
			panic(err)
		}

		for j := 0; j < len(jsonData); j++ {
			// created_at の文字列を時刻としてパースし、JST に変換
			createdAt := jsonData[j].CreatedAt
			timeObj, _ := time.Parse(githubDateFormat, createdAt)
			timeObj = timeObj.In(jst)
			fmt.Println(jsonData[j].CreatedAt)
			fmt.Println(timeObj)

			// current time - tweet time
			duration := currentDayObj.Sub(timeObj)
			durationSec := int(duration.Seconds())
			if 0 <= durationSec && durationSec < durationDay {
				commitCount++
			} else if durationDay <= durationSec {
				goto Checked
			}
		}
	}

Checked:
	return
}

// FoxStatus : get today's status
func FoxStatus(api *anaconda.TwitterApi) {

	// 一日のツイート数を取得
	count, day := countTweetNum(api)
	eventNum := getGithubCommitCount(api)
	tweetHeader := "(っ ॑꒳ ॑c).+(" + day.Format(dateFormat) + " 0:00 Report)"
	tweetCountText := "前日ツイート数（bot 除外）:" + strconv.Itoa(count)
	tweetCommitText := "GitHub Event 数 : " + strconv.Itoa(eventNum)
	tweetText := tweetHeader + "\n" + tweetCountText + "\n" + tweetCommitText
	fmt.Println(tweetText)
	//_, postErr := api.PostTweet(tweetText+"\n(bot)", nil)

	//if postErr != nil {
	//	panic(postErr)
	//}
	return
}
