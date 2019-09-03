package main

import (
	"flag"
	"fmt"
	"github.com/PhilLar/webshorten/short"
	"log"
	"sync"
)

var flagRelink = flag.Bool("relink", false, "use rel.ink service to shorten URL")

// runInParallel is a Wrapper for CleanUrl() and Relink()
func runInParallel(a func(urlLink string) (string, error), args []string) map[string]error {
	shorts := make(map[string]error)
	var mutex = &sync.Mutex{}
	var wg sync.WaitGroup
	wg.Add(len(args))
	for _, arg := range args {
		go func(arg string) {
			defer wg.Done()
			short, err := a(arg)
			mutex.Lock()
			shorts[short] = err
			mutex.Unlock()
		}(arg)
	}
	wg.Wait()
	return shorts
}

func shortenUrls() map[string]error {
	flag.Parse()
	if !*flagRelink {
		return runInParallel(short.CleanURL, flag.Args())
	}
	return runInParallel(short.Relink, flag.Args())
}

func main() {
	resultMap := shortenUrls()
	for short, err := range resultMap {
		if err == nil {
			fmt.Println(short)
		} else {
			log.Fatal(err)
		}
	}
}
