package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"path/filepath"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	matches, err := filepath.Glob("texts/*.txt")
	if err != nil {
		log.Fatal(err)
	}

	file := matches[rand.Intn(len(matches))]
	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(content))
}
