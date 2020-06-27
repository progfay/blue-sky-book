package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"golang.org/x/text/encoding/japanese"
)

type Content struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Sha         string `json:"sha"`
	Size        int64  `json:"size"`
	Url         string `json:"url"`
	HtmlUrl     string `json:"html_url"`
	GitUrl      string `json:"git_url"`
	DownloadUrl string `json:"download_url"`
	Type        string `json:"type"`
	Links       struct {
		Self string `json:"self"`
		Git  string `json:"git"`
		Html string `json:"html"`
	} `json:"_links"`
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GetRandomAozoraBookFileUrl() string {
	path := "cards"

	for {
		url := fmt.Sprintf("https://api.github.com/repos/aozorahack/aozorabunko_text/contents/%s?ref=master", path)
		res, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}

		body, err := ioutil.ReadAll(res.Body)

		contents := make([]Content, 0)
		json.Unmarshal(body, &contents)
		res.Body.Close()

		content := contents[rand.Intn(len(contents))]
		if content.Type == "file" {
			return content.DownloadUrl
		}
		path = content.Path
	}
}

func GetAozoraBookLines(url string) []string {
	lines := make([]string, 0)
	res, err := http.Get(GetRandomAozoraBookFileUrl())
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	decoder := japanese.ShiftJIS.NewDecoder()
	reader := bufio.NewReader(decoder.Reader(res.Body))

	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		lines = append(lines, string(line))
	}

	return lines
}

func main() {
	fmt.Println(GetAozoraBookLines(GetRandomAozoraBookFileUrl()))
}
