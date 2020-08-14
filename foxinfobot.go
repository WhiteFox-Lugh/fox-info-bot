package main

import (
	"flag"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"os"
)

// SetAPI : setting client
func SetAPI() *anaconda.TwitterApi {
	var consumerKey = os.Getenv("CONSUMER_KEY")
	var consumerKeySecret = os.Getenv("CONSUMER_KEY_SECRET")
	var accessToken = os.Getenv("ACCESS_TOKEN")
	var accessTokenSecret = os.Getenv("ACCESS_TOKEN_SECRET")

	ret := anaconda.NewTwitterApiWithCredentials(accessToken, accessTokenSecret, consumerKey, consumerKeySecret)

	return ret
}

func main() {
	flag.Parse()
	f := flag.Arg(0)

	// authentication
	api := SetAPI()

	if f == "img" {
		// post img tweet
		PostImgTweet(api)
	} else if f == "text" {
		// post text tweet
		PostTextTweet(api)
	} else if f == "weather" {
		WeatherFox(api)
	} else if f == "weatherforecast" {
		WeatherForecast(api)
	}
	fmt.Println("end")
}
