package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	PhotoApi = "https://api.pexels.com/v1"
	VideoApi = "https://api.pexels.com/videos"
)

type Client struct {
	Token          string
	hc             http.Client
	RemainingTimes int32
}

func NewClient(token string) *Client {
	c := http.Client{}
	return &Client{Token: token, hc: c}
}

type SearchResult struct {
	Page         int32   `json:"page"`
	PerPage      int32   `json:"per_page"`
	TotalResults int32   `json:"total_Results"`
	NextPage     string  `json:"next_page"`
	Photos       []Photo `json:"photos"`
}

type Photo struct {
	Id              int32       `json:"id"`
	Width           int32       `json:"width"`
	Height          int32       `json:"height"`
	Url             string      `json:"url"`
	Photographer    string      `json:"photographer"`
	PhotographerUrl string      `json:"photographer_url"`
	Src             PhotoSource `json:"src"`
}

type PhotoSource struct {
	Original  string `json:"Original"`
	Large     string `json:"Large"`
	Large2x   string `json:"Large2x"`
	Medium    string `json:"Medium"`
	Small     string `json:"Small"`
	Potrait   string `json:"Potrait"`
	Square    string `json:"Square"`
	Landscape string `json:"Landscape"`
	Tiny      string `json:"Tiny            "`
}

type CuratedResult struct {
	Page     int32   `json:"page"`
	PerPage  int32   `json:"per_page"`
	NextPage string  `json:"next_page"`
	Photos   []Photo `json:"photos"`
}

type Video struct {
	Id            int32           `json:"id"`
	Width         int32           `json:"width"`
	Height        int32           `json:"height"`
	Url           int32           `json:"url"`
	Image         string          `json:"image"`
	FullRes       interface{}     `json:"full_res"`
	Duration      float64         `json:"duration"`
	VideoFiles    []VideoFiles    `json:"video_files"`
	VideoPictures []VideoPictures `json:"video_pictures"`
}

type VideoSearchResult struct {
	Page         int32   `json:"page"`
	PerPage      int32   `json:"per_page"`
	TotalResults int32   `json:"total_results"`
	NextPage     string  `json:"next_page"`
	Videos       []Video `json:"videos"`
}

type PopularVideos struct {
	Page         int32   `json:"page"`
	PerPage      int32   `json:"per_page"`
	TotalResults int32   `json:"total_results"`
	Url          string  `json:"url"`
	Videos       []Video `json:"videos"`
}

type VideoFiles struct {
	Id       int32  `json:"id"`
	Quality  string `json:"quality"`
	FileType string `json:"file_type"`
	Width    int32  `json:"width"`
	Height   int32  `json:"height"`
	Link     string `json:"link"`
}

type VideoPictures struct {
	Id      int32  `json:"id"`
	Picture string `json:"picture"`
	Number  int32  `json:"number"`
}

func (c *Client) SearchPhotos(query string, perPage, page int) (*SearchResult, error) {
	url := fmt.Sprintf(PhotoApi+"/search?query=%s&per_page=%d&page=%d", query, perPage, page)
	log.Println(url)
	resp, err := c.requstDoWithAuth("GET", url)
	log.Println(resp)
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result SearchResult
	err = json.Unmarshal(data, &result)
	defer resp.Body.Close()
	return &result, err

}

func (c *Client) requstDoWithAuth(method string, url string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", c.Token)
	log.Println(c.Token)
	resp, err := c.hc.Do(req)
	if err != nil {
		return resp, err
	}
	log.Println(resp)
	times, err := strconv.Atoi(resp.Header.Get("X-Ratelimit-Remaining"))
	if err != nil {
		return resp, nil
	} else {
		c.RemainingTimes = int32(times)
	}
	return resp, nil

}

func (c *Client) CuratedPhotos(perPage, page int) (*CuratedResult, error) {
	url := fmt.Sprintf(PhotoApi+"/curated?per_page=%d&page=%d", perPage, page)
	resp, err := c.requstDoWithAuth("GET", url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var CuratedResult CuratedResult
	err = json.Unmarshal(data, &CuratedResult)

	return &CuratedResult, nil
}

func (c *Client) GetPhoto(id int32) (*Photo, error) {
	url := fmt.Sprintf(PhotoApi+"/photos/%d", id)
	resp, err := c.requstDoWithAuth("GET", url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result Photo
	err = json.Unmarshal(data, result)
	return &result, nil
}

func (c *Client) getRandomPhoto() (*Photo, error) {
	rand.Seed(time.Now().Unix())
	randNum := rand.Intn(1001)
	result, err := c.CuratedPhotos(1, randNum)
	if err == nil && len(result.Photos) == 1 {
		return &result.Photos[0], nil
	}
	return nil, err
}

// func (c *Client) GetRandomVideo(*Video, error) {
// 	rand.Seed(time.Now().Unix())
// 	randNum := rand.Intn(1001)
// 	result, err := c.GetPopularVideo(1, randNum)
// 	if err == nil && len(result.Videos) == 1 {
// 		return &result.Videos[0], nil
// 	}
// 	return nil, err

// }

func (c *Client) SearchVideo(query string, perpage, page int) (*VideoSearchResult, error) {
	url := fmt.Sprintf(VideoApi+"/search?query=%s&per_page=%d&page=%d", query, perpage, page)
	log.Println(url)
	resp, err := c.requstDoWithAuth("GET", url)
	log.Println(resp)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}
	var VideoResult VideoSearchResult
	err = json.Unmarshal(data, &VideoResult)
	return &VideoResult, nil
}

func (c *Client) GetPopularVideo(perpage, page int) (*PopularVideos, error) {
	url := fmt.Sprintf(VideoApi+"/popular?per_page=%d&page=%d", perpage, page)
	resp, err := c.requstDoWithAuth("GET", url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	//log.Println(data)
	if err != nil {
		return nil, err
	}
	var result PopularVideos
	err = json.Unmarshal(data, &result)
	return &result, nil

}

func (c *Client) GetRemainingRequestForThisMonth() int32 {
	return c.RemainingTimes
}

func main() {
	os.Setenv("PexelsToken", "Token")
	var TOKEN = os.Getenv("PexelsToken")

	var c = NewClient(TOKEN)

	//result, err := c.SearchPhotos("waves", 15, 1)
	//result, err := c.SearchVideo("waves", 15, 1)
	//result, err := c.getRandomPhoto()
	result, err := c.GetPopularVideo(2, 2)
	if err != nil {
		fmt.Errorf("Error at Line 26 , %v", err)
	}

	// if result.Page == 0 {
	// 	fmt.Errorf("Search Result is wrong")
	// }

	fmt.Println(result)

}
