package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"github.com/ChimeraCoder/anaconda"
	"io/ioutil"
	"net/http"
	"net/url"
)

// FoxImage : json data from randomfox
type FoxImage struct {
	Image string `json:"image"`
	Link  string `json:"link"`
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
	api.PostTweet("らんだむこゃーんいめーじ (bot)\n"+imgLink, v)
	return
}
