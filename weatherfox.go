package main

import (
	"encoding/json"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// WeatherInfo : weather data from OpenWeather
type WeatherInfo struct {
	Current Hourly   `json:"current"`
	Offset  int      `json:"timezone_offset"`
	Hourly  []Hourly `json:"hourly"`
}

// Hourly : hourly forecast
type Hourly struct {
	Dt        uint64    `json:"dt"`
	Temp      float64   `json:"temp"`
	FeelsLike float64   `json:"feels_like"`
	Pressure  int       `json:"pressure"`
	Humidity  int       `json:"humidity"`
	Weather   []Weather `json:"weather"`
}

// Weather : weather information
type Weather struct {
	ID int `json:"id"`
}

// WeatherFox : show weather forecast on screen name
func WeatherFox(api *anaconda.TwitterApi) {
	const screenName = "Arthur_Lugh"
	//jsonData := getJSON()
	//	weatherEmojiStr := weatherEmoji(strconv.Itoa(jsonData.Current.Weather[0].ID))

	userObj, err := api.GetUsersShow(screenName, nil)

	if err != nil {
		panic(err)
	}

	fmt.Println(userObj)

	return
}

// round : round function
func round(f float64) float64 {
	return math.Floor(f + .5)
}

// WeatherForecast : post today's weather forecast
func WeatherForecast(api *anaconda.TwitterApi) {
	const tweetTextHeader = "(っ ॑꒳ ॑)っ/☀おてんき at 京都（左京区）\n"
	var weatherEmojiStr string
	var tweetText string
	var minTemperature float64
	var maxTemperature float64
	weatherEmojiStr = ""
	minTemperature = 100.0
	maxTemperature = -100.0

	// get json data
	jsonData := getJSON()

	// generate weather string and find min/max temperature
	for i := 0; i < 24; i++ {
		if i == 0 {
			weatherEmojiStr += "AM : "
		} else if i == 12 {
			weatherEmojiStr += "\nPM : "
		}
		weatherEmojiStr += weatherEmoji(strconv.Itoa(jsonData.Hourly[i].Weather[0].ID))

		minTemperature = math.Min(minTemperature, jsonData.Hourly[i].Temp)
		maxTemperature = math.Max(maxTemperature, jsonData.Hourly[i].Temp)
	}

	// get time
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	nowUTC := time.Now().UTC()
	nowJST := nowUTC.In(jst)
	fmt.Println(nowJST.Format("2006-01-02"))

	roundedMinTemp := strconv.FormatFloat(round(minTemperature), 'f', 0, 32)
	roundedMaxTemp := strconv.FormatFloat(round(maxTemperature), 'f', 0, 32)
	tempStr := "気温: 最高 " + roundedMaxTemp + "℃ / 最低 " + roundedMinTemp + "℃"

	tweetText = tweetTextHeader + "\n" + weatherEmojiStr + "\n" + tempStr
	_, err := api.PostTweet(tweetText+"\n(bot)", nil)

	if err != nil {
		panic(err)
	}

	return
}

// getJSON : get weather information from OpenWeather
func getJSON() WeatherInfo {
	// ret : JSON data
	var ret WeatherInfo
	var appID = os.Getenv("APP_ID")

	url := "https://api.openweathermap.org/data/2.5/onecall?lat=35.02&lon=135.78&units=metric&appid=" + appID

	req, _ := http.NewRequest("GET", url, nil)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err.Error())
	}

	defer res.Body.Close()

	// parse json
	byteArray, err := ioutil.ReadAll(res.Body)
	if err := json.Unmarshal(byteArray, &ret); err != nil {
		panic(err)
	}

	return ret
}

// weatherEmoji : return weather emoji
func weatherEmoji(str string) string {
	var ret string
	if strings.HasPrefix(str, "2") {
		// Thunderstorm
		ret = "🌩"
	} else if strings.HasPrefix(str, "3") {
		// Drizzle
		ret = "☂"
	} else if strings.HasPrefix(str, "5") {
		// Rain
		ret = "☔"
	} else if strings.HasPrefix(str, "6") {
		// Snow
		ret = "❄"
	} else if strings.HasPrefix(str, "7") {
		// Atmosphere : mist / fog ...
		ret = "🌫"
	} else if strings.HasPrefix(str, "8") {
		if strings.HasSuffix(str, "00") {
			// Clear
			ret = "☀"
		} else if strings.HasSuffix(str, "01") || strings.HasSuffix(str, "02") {
			// few clouds or scattered clouds
			ret = "⛅"
		} else {
			// broken clouds or overcast clouds
			ret = "☁"
		}
	} else {
		ret = "❔"
	}
	return ret
}
