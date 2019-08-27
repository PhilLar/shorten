package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

var api string = "https://cleanuri.com/api/v1/shorten"

type ApiAnswer struct {
	Result_url string
	Error      string
}

func ShortenUrl(args []string) (string, error) {
	if len(args) == 2 {
		url_link := args[1]
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
		} else {
			if answer.Error != "" {
				return "", errors.New(answer.Error)
			} else {
				return answer.Result_url, nil
			}
		}
		return string(body), nil
	}
	return "", errors.New("Threre're must be 1 arg - a url!")
}

func main() {
	result, err := ShortenUrl(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
}
