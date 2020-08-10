package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"github.com/ChimeraCoder/anaconda"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"
)

// FoxText : array of tweet text
type FoxText struct {
	Body []string `json:"text"`
}

// FoxImage : json data from randomfox
type FoxImage struct {
	Image string `json:"image"`
	Link  string `json:"link"`
}

// SetAPI : setting client
func SetAPI() *anaconda.TwitterApi {
	var consumerKey = os.Getenv("CONSUMER_KEY")
	var consumerKeySecret = os.Getenv("CONSUMER_KEY_SECRET")
	var accessToken = os.Getenv("ACCESS_TOKEN")
	var accessTokenSecret = os.Getenv("ACCESS_TOKEN_SECRET")

	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerKeySecret)
	ret := anaconda.NewTwitterApi(accessToken, accessTokenSecret)

	return ret
}

// PostTextTweet : post tweet text
func PostTextTweet(api *anaconda.TwitterApi) {
	// body : array of tweet body
	var body FoxText

	jsonData, err := ioutil.ReadFile("./tweet/tweet_body.json")
	if err != nil {
		panic(err.Error())
	}

	json.Unmarshal(jsonData, &body)

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(body.Body))
	text := body.Body[n]
	tweet, err := api.PostTweet(text, nil)

	if err != nil {
		panic(err)
	}

	print(tweet.Text)
	return
}

// PostImgTweet : post image tweet
func PostImgTweet(api *anaconda.TwitterApi) {
	const urlRandFox = "https://randomfox.ca/floof/"
	var foxImage FoxImage

	// http request
	response, reqErr := http.Get(urlRandFox)
	if reqErr != nil {
		panic(reqErr)
	}
	defer response.Body.Close()

	// parse json
	byteArray, _ := ioutil.ReadAll(response.Body)
	if err := json.Unmarshal(byteArray, &foxImage); err != nil {
		panic(err)
	}

	// img url and img link
	imgURL := foxImage.Image
	imgLink := foxImage.Link

	// http request for img data
	responseImg, reqErr := http.Get(imgURL)
	if reqErr != nil {
		panic(reqErr)
	}
	defer responseImg.Body.Close()

	// parse img
	imgByteArray, _ := ioutil.ReadAll(responseImg.Body)

	// encode
	encodedImg := b64.StdEncoding.EncodeToString(imgByteArray)

	// upload img
	img, _ := api.UploadMedia(encodedImg)

	// tweet
	v := url.Values{}
	v.Add("media_ids", img.MediaIDString)
	api.PostTweet("こゃーん\n"+imgLink, v)
	return
}

func main() {
	// authentication
	api := SetAPI()

	// post img tweeet
	PostImgTweet(api)
}
