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
	"bytes"
)

const CleanUrl_api string = "https://cleanuri.com/api/v1/shorten"
const Relink_api string = "https://rel.ink/api/links/"

var flagRelink *bool = flag.Bool("relink", false, "a bool")

type CleanApiAnswer struct {
	Result_url string
	Error      string
}



func ShortenUrl() (string, error) {
	flag.Parse()
	if !*flagRelink {
		return  CleanUrl(flag.Arg(0))
	} else {
		return Relink(flag.Arg(0))
	}
}

func CleanUrl(url_link string) (string, error) {
	resp, err := http.PostForm(CleanUrl_api, url.Values{
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
	answer := &CleanApiAnswer{}
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

func Relink(url_link string) (string, error) {
	jsn, err := json.Marshal(map[string]string{
		"url":url_link,
	})
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.Post("https://rel.ink/api/links/", "application/json", bytes.NewBuffer(jsn))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	answer := make(map[string]interface{})
	err = json.Unmarshal(body, &answer)
	if err != nil {
		log.Fatal(err)
	}
	clean_url := "https://rel.ink/" + answer["hashid"].(string)
	return  clean_url, err
}

func main() {
	result, err := ShortenUrl()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
}
