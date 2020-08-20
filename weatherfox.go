package main

import (
	"encoding/json"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/dghubble/oauth1"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"net/url"
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
	Daily   []Daily  `json:"daily"`
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

// Daily : daily forecast
type Daily struct {
	Sunrise  int64   `json:"sunrise"`
	Sunset   int64   `json:"sunset"`
	Pressure int     `json:"pressure"`
	Humidity int     `json:"humidity"`
	Wind     float64 `json:"wind_speed"`
}

// Weather : weather information
type Weather struct {
	ID int `json:"id"`
}

// WeatherFox : show weather forecast on screen name
func WeatherFox(api *anaconda.TwitterApi) {
	const screenName = "Arthur_Lugh"
	const baseURL = "https://api.twitter.com/1.1/account/update_profile.json"
	var consumerKey = os.Getenv("CONSUMER_KEY")
	var consumerKeySecret = os.Getenv("CONSUMER_KEY_SECRET")
	var accessToken = os.Getenv("ACCESS_TOKEN")
	var accessTokenSecret = os.Getenv("ACCESS_TOKEN_SECRET")
	jsonData := getJSON()
	weatherEmojiStr, face := weatherEmoji(strconv.Itoa(jsonData.Current.Weather[0].ID))
	fmt.Println("weather -> " + weatherEmojiStr)

	config := oauth1.NewConfig(consumerKey, consumerKeySecret)
	token := oauth1.NewToken(accessToken, accessTokenSecret)
	client := config.Client(oauth1.NoContext, token)

	userObj, err := api.GetUsersShow(screenName, nil)

	if err != nil {
		panic(err)
	}

	currentName := userObj.Name
	idx := strings.Index(currentName, "(")
	if idx == -1 {
		idx = len(currentName)
	}
	newName := currentName[:idx] + face + weatherEmojiStr

	values := url.Values{}
	values.Add("name", newName)

	//リクエストの送信
	request, err := http.NewRequest("POST", baseURL+"?"+values.Encode(), nil)
	if err != nil {
		panic(err)
	}

	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}

	fmt.Println(response)

	return
}

// round : round function
func round(f float64) float64 {
	return math.Floor(f + .5)
}

// WeatherForecast : post today's weather forecast
func WeatherForecast(api *anaconda.TwitterApi) {
	const tweetTextHeader = "(っ ॑꒳ ॑)っ/ 天気(京都市左京区)\n"
	var weatherEmojiStr string
	var emoji string
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
		emoji, _ = weatherEmoji(strconv.Itoa(jsonData.Hourly[i].Weather[0].ID))
		weatherEmojiStr += emoji

		minTemperature = math.Min(minTemperature, jsonData.Hourly[i].Temp)
		maxTemperature = math.Max(maxTemperature, jsonData.Hourly[i].Temp)
	}

	// get temperature
	roundedMinTemp := strconv.FormatFloat(round(minTemperature), 'f', 0, 32)
	roundedMaxTemp := strconv.FormatFloat(round(maxTemperature), 'f', 0, 32)
	tempStr := "気温🌡: 最高 " + roundedMaxTemp + "℃ / 最低 " + roundedMinTemp + "℃"
	if maxTemperature >= 35 {
		tempStr += " (猛暑日)"
	} else if maxTemperature >= 30 {
		tempStr += " (真夏日)"
	} else if maxTemperature >= 25 {
		tempStr += " (夏日)"
	} else if maxTemperature < 0 {
		tempStr += " (真冬日)"
	} else if minTemperature < 0 {
		tempStr += "(冬日)"
	}

	// get windspeed
	windspeed := jsonData.Daily[0].Wind
	windStr := "風速: " + strconv.FormatFloat(jsonData.Daily[0].Wind, 'f', 0, 32) + " m/s"
	if windspeed >= 30 {
		windStr += " (猛烈な風)"
	} else if windspeed >= 20 {
		windStr += " (非常に強い風)"
	} else if windspeed >= 15 {
		windStr += " (強い風)"
	} else if windspeed >= 10 {
		windStr += " (やや強い風)"
	}

	// get humidity
	humidity := jsonData.Daily[0].Humidity
	humidStr := "湿度: " + strconv.Itoa(humidity) + " %"

	// get pressure
	pressure := jsonData.Daily[0].Pressure
	preStr := "気圧: " + strconv.Itoa(pressure) + " hPa"

	// sunset and sunrise
	sunriseJST := time.Unix(jsonData.Daily[0].Sunrise, 0)
	sunriseJST = sunriseJST.Add(9 * time.Hour)
	sunsetJST := time.Unix(jsonData.Daily[0].Sunset, 0)
	sunsetJST = sunsetJST.Add(9 * time.Hour)
	const layout = "15:04:05"
	sunTime := "日の出 " + sunriseJST.Format(layout) + " / 日の入り " + sunsetJST.Format(layout)

	tweetText = tweetTextHeader + "\n" + weatherEmojiStr + "\n" + tempStr + "\n" + windStr + "\n" + humidStr + "\n" + preStr + "\n" + sunTime
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
func weatherEmoji(str string) (weather string, face string) {
	var faceArray = [...]string{"(っ ॑꒳ ॑)っ/", "(っ˘꒳˘)っ/", "() ੭•͈ω•͈)っ/", "(*`꒳´)っ/"}
	const scared = "(っºΔº)っ/"

	if strings.HasPrefix(str, "2") {
		// Thunderstorm
		weather = "⚡"
		face = scared
		return
	} else if strings.HasPrefix(str, "3") {
		// Drizzle
		weather = "☂"
	} else if strings.HasPrefix(str, "5") {
		// Rain
		weather = "☔"
	} else if strings.HasPrefix(str, "6") {
		// Snow
		weather = "❄"
	} else if strings.HasPrefix(str, "7") {
		// Atmosphere : mist / fog ...
		weather = "🌫"
	} else if strings.HasPrefix(str, "8") {
		if strings.HasSuffix(str, "00") || strings.HasSuffix(str, "01") {
			// Clear (800) or few clouds (801 : 11%-25%)
			weather = "☀"
		} else if strings.HasSuffix(str, "02") || strings.HasSuffix(str, "03") {
			// scattered clouds (802 : 25%-50%) or broken clouds (803 : 51%-84%)
			weather = "⛅"
		} else {
			// broken clouds or overcast clouds
			weather = "☁"
		}
	} else {
		weather = "❔"
	}
	// face generate
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(faceArray))
	face = faceArray[n]
	return
}
