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
	//"go/build"
	//"os"
	"sync"
)

const CleanUrl_api string = "https://cleanuri.com/api/v1/shorten"
const Relink_api string = "https://rel.ink/api/links/"

var flagRelink *bool = flag.Bool("relink", false, "a bool")

type CleanUrlAnswer struct {
	ResultUrl  string `json:"result_url"`
	Error      string
}

// Wrapper for CleanUrl() and Relink()
func RunInParallel(a func(url_link string) (string, error), args []string) map[string]error {
	shorts := make(map[string]error)
	var mutex = &sync.Mutex{}
	var wg sync.WaitGroup
	var short string
	var err error
	wg.Add(len(args))
	for _, arg := range args {
		go func(arg string) {
			defer wg.Done()
			short, err = a(arg)
			mutex.Lock()
            shorts[short] = err
            mutex.Unlock()
		}(arg)
	}
	wg.Wait()
	return shorts
}


func ShortenUrls() map[string]error {
	flag.Parse()
	if !*flagRelink {
		// shorts := make(map[string]error)
		// str := make(chan string)
		// err := make(chan error)
		// for arg := range flag.Args() {
		// 	go func() {

		// 	}
		// }
		return RunInParallel(CleanUrl, flag.Args())
	} else {
		return RunInParallel(Relink, flag.Args())
	}
}

func CleanUrl(url_link string) (string, error) {
	resp, err := http.PostForm(CleanUrl_api, url.Values{
		"url": {url_link},
	})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	answer := &CleanUrlAnswer{}
	err = json.Unmarshal(body, answer)
	if err != nil {
		return "", err
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
		return "", err
	}
	resp, err := http.Post("https://rel.ink/api/links/", "application/json", bytes.NewBuffer(jsn))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	answer := make(map[string]interface{})
	err = json.Unmarshal(body, &answer)
	if err != nil {
		return "", err
	}
	clean_url := "https://rel.ink/" + answer["hashid"].(string)
	return  clean_url, err
}

func main() {
	//fmt.Println("sdf", build.Default.GOPATH)
	result_map := ShortenUrls()
	for short, err := range result_map {
		if err == nil {
			fmt.Println(short)
		} else {
			log.Fatal(err)
		}
	}
}
