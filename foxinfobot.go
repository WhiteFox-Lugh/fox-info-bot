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

// PrintHelp : if arg has no value, show help message
func PrintHelp() {
	fmt.Println("---------------------")
	fmt.Println("foxinfobot.go Usage")
	fmt.Println("---------------------")
	fmt.Println("img : post image tweet (using randomfox api)")
	fmt.Println("text : post text tweet about fox")
	fmt.Println("weather : change twitter screen name (depends on current weather at Kyoto)")
	fmt.Println("foxstatus : tweet todays twitter usage")
	fmt.Println("---------------------")
	return
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
	} else if f == "foxstatus" {
		FoxStatus(api)
	} else {
		PrintHelp()
	}
	fmt.Println("end")
}
