package main

import (
	"fmt"
	"log"
	"github.com/PhilLar/webshorten/short"
)

func main() {
	resultMap := short.ShortenUrls()
	for short, err := range resultMap {
		if err == nil {
			fmt.Println(short)
		} else {
			log.Fatal(err)
		}
	}
}
