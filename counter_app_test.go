package main

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	com "gitlab.com/leapbit-practice/tweety-lib-communication/comms"
	tw "gitlab.com/leapbit-practice/tweety-lib-twitter/twitter"
)

func TestSendTweetsToDB(t *testing.T) {
	cdb := HttpClientDB{
		RequestClient: HttpRequestClient{http.Client{Timeout: time.Duration(15) * time.Second}},
		DbIpAndPort:   "172.17.31.148:8080",
	}

	userTweets := []tw.RespTwitterApiTweet{

		{
			Created_at: tw.TwitterTime{Time: time.Now()},
			Id:         16,
			Id_str:     "16",
			Text:       "bok ana evo ti moj tweet potrudio sam se da ga smislim",
			Url:        "http://mojtweet/blabla/12345",
			User: struct {
				Id          uint64 `json:"id"`
				Screen_name string `json:"screen_name"`
			}{
				Id:          1068831,
				Screen_name: "planky",
			},
		},
		{
			Created_at: tw.TwitterTime{Time: time.Now()},
			Id:         17,
			Id_str:     "17",
			Text:       "bok ana evo ti još jedan tweet nadam se da radi",
			Url:        "http://mojdrugitweet/blabla/67890",
			User: struct {
				Id          uint64 `json:"id"`
				Screen_name string `json:"screen_name"`
			}{
				Id:          1068831,
				Screen_name: "planky",
			},
		},
	}

	rankedWordCount := rankMostUsedWords(userTweets)

	err, errMsg := cdb.sendTweetsToDB(userTweets, rankedWordCount)
	if err != nil {
		t.Fatalf("Internal error while sending tweets to database. Error: %s", err)
	}
	if errMsg != nil {
		t.Fatalf("External error while sending tweets to database. Error: %s", errMsg)
	}
}

func TestRankMostUsedWords(t *testing.T) {
	userTweets := []tw.RespTwitterApiTweet{

		{
			Created_at: tw.TwitterTime{Time: time.Now()},
			Id:         16,
			Id_str:     "16",
			Text:       "bok ana evo ti moj tweet potrudio sam se da ga smislim",
			Url:        "http://mojtweet/blabla/12345",
			User: struct {
				Id          uint64 `json:"id"`
				Screen_name string `json:"screen_name"`
			}{
				Id:          1068831,
				Screen_name: "planky",
			},
		},
		{
			Created_at: tw.TwitterTime{Time: time.Now()},
			Id:         17,
			Id_str:     "17",
			Text:       "bok ana evo ti još jedan tweet nadam se da radi",
			Url:        "http://mojdrugitweet/blabla/67890",
			User: struct {
				Id          uint64 `json:"id"`
				Screen_name string `json:"screen_name"`
			}{
				Id:          1068831,
				Screen_name: "planky",
			},
		},
	}

	resultWordCount := rankMostUsedWords(userTweets)

	fmt.Println(resultWordCount)
}

func TestGetTweetsFromTwitter(t *testing.T) {
	ctw := HttpClientTW{
		RequestClient: HttpRequestClient{Client: http.Client{Timeout: time.Duration(15) * time.Second}},
		TweetNo:       10,
		Bearer:        bearer,
	}

	tweets, err := ctw.getTweetsFromTwitter("2625272871")
	if err != nil {
		t.Fatalf("Error while getting tweets from twitter. Error: %s", err.Error())
	}
	t.Log(tweets)
}

func TestGetImageUrlsFromTwitter(t *testing.T) {
	ctw := HttpClientTW{
		RequestClient: HttpRequestClient{Client: http.Client{Timeout: time.Duration(15) * time.Second}},
		TweetNo:       10,
		Bearer:       bearer,
	}

	urlProfileImage, urlBanner, err := ctw.getImageUrlsFromTwitter("813286")
	if err != nil {
		t.Fatalf("Error while getting image urls from twitter. Error: %s", err.Error())
	}

	fmt.Println(urlProfileImage)
	fmt.Println(urlBanner)
}

func TestDownloadFile(t *testing.T) {
	rc := HttpRequestClient{Client: http.Client{Timeout: time.Duration(15) * time.Second}}
	f, err := rc.DownloadFile("https://pbs.twimg.com/profile_images/1329647526807543809/2SGvnHYV_normal.jpg")
	if err != nil {
		t.Fatalf("Cannot download file. Error: %s", err.Error())
	}

	fmt.Printf("File(byte array): %v\n", f)
	fmt.Printf("File(string): %s\n", string(f))
}

func TestSendImagesToDB(t *testing.T) {
	userId := "2458938607"
	cdb := HttpClientDB{
		RequestClient: HttpRequestClient{http.Client{Timeout: time.Duration(15) * time.Second}},
		DbIpAndPort:   "tweety-dbsaver-tck-test.demobet.lan:8080",
	}

	ctw := HttpClientTW{
		RequestClient: HttpRequestClient{Client: http.Client{Timeout: time.Duration(15) * time.Second}},
		TweetNo:       10,
		Bearer:        bearer,
	}

	urlProfileImage, urlBanner, err := ctw.getImageUrlsFromTwitter(userId)
	com.TweetyLog(com.DEBUG, fmt.Sprintf("Obtained urls from twitter: %s %s", urlProfileImage, urlBanner))
	if err != nil {
		t.Fatalf("Error while getting image urls from twitter. Error: %s", err.Error())
	}
	imgNames := make([]string, 0)
	dataToZip := make([][]byte, 0)

	if urlProfileImage != "" {
		com.TweetyLog(com.INFO, fmt.Sprintf("Downloading profile image of user %s...", userId))
		img1, err := ctw.RequestClient.DownloadFile(urlProfileImage)
		com.TweetyLog(com.DEBUG, fmt.Sprintf("Img1: %s\n", string(img1)))
		if err != nil {
			com.TweetyLog(com.ERROR, fmt.Sprintf("Worker failed to process user id: %s. Error: %s", userId, err.Error()))
			t.FailNow()
		}
		com.TweetyLog(com.INFO, fmt.Sprintf("Downloading profile image of user %s DONE.", userId))
		imgNames = append(imgNames, "profile_image.jpg")
		dataToZip = append(dataToZip, img1)
	}

	if urlBanner != "" {
		com.TweetyLog(com.INFO, fmt.Sprintf("Downloading profile banner of user %s...", userId))
		img2, err := ctw.RequestClient.DownloadFile(urlBanner)
		com.TweetyLog(com.DEBUG, fmt.Sprintf("Img2: %x\n", img2))
		if err != nil {
			com.TweetyLog(com.ERROR, fmt.Sprintf("Worker failed to process user id: %s. Error: %s", userId, err.Error()))
			t.FailNow()
		}
		com.TweetyLog(com.INFO, fmt.Sprintf("Downloading profile banner of user %s DONE.", userId))
		imgNames = append(imgNames, "banner.png")
		dataToZip = append(dataToZip, img2)
	}

	if len(dataToZip) > 0 {
		com.TweetyLog(com.INFO, fmt.Sprintf("Zipping images of user %s...", userId))
		zippedData, err := ZipFiles(imgNames, dataToZip)
		com.TweetyLog(com.DEBUG, fmt.Sprintf("Zipped data: %x\n", zippedData))
		if err != nil {
			com.TweetyLog(com.ERROR, fmt.Sprintf("Worker failed to process user id: %s. Error: %s", userId, err.Error()))
			t.FailNow()
		}
		com.TweetyLog(com.INFO, fmt.Sprintf("Zipping images of user %s DONE.", userId))

		com.TweetyLog(com.INFO, fmt.Sprintf("Sending images of user %s to database...", userId))
		err, errMsg := cdb.sendImagesToDB(userId, zippedData)
		if err != nil || errMsg != nil {
			com.TweetyLog(com.ERROR, fmt.Sprintf("Worker failed to process user id: %s. Error: %s", userId, err.Error()))
			t.FailNow()
		}
		com.TweetyLog(com.INFO, fmt.Sprintf("Sending images of user %s to database DONE.", userId))
	}
}
