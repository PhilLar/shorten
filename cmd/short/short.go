package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

var api string = "https://cleanuri.com/api/v1/shorten"

var flagRelink *bool = flag.Bool("relink", false, "a bool")

type ApiAnswer struct {
	Result_url string
	Error      string
}

func ShortenUrl() (string, error) {
	flag.Parse()
	if !*flagRelink {
		return  CleanUrl(flag.Arg(0))
	} else {
		return "", errors.New("underdone")
	}
}

func CleanUrl(url_link string) (string, error) {
	resp, err := http.PostForm(api, url.Values{
		"url": {url_link},
	})
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	answer := &ApiAnswer{}
	err = json.Unmarshal(body, answer)
	if err != nil {
		log.Fatalln(err)
	}
	if answer.Error != "" {
		return "", errors.New(answer.Error)
	} else {
		return answer.Result_url, nil
	}
}

//func CleanUrl(url_link string) (string, error) {

func main() {
	result, err := ShortenUrl()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
}
