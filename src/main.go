package main

import (
	"fmt"
	"log"
	"math/rand"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

	b "github.com/progfay/blue-sky-book/book"
	l "github.com/progfay/blue-sky-book/line"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	matches, err := filepath.Glob("texts/*.txt")
	if err != nil {
		log.Fatal(err)
	}

	for _, path := range matches {
		book := b.NewBook(path)
		lines, err := book.GetLinesFromBook()
		if err != nil {
			log.Fatal(err)
		}

		for _, line := range lines {
			for _, sentence := range l.ParseLine(line) {
				if strings.HasPrefix(sentence, "「") || strings.HasPrefix(sentence, "（") {
					continue
				}
				if !strings.HasSuffix(sentence, "。") {
					continue
				}
				if strings.Contains(sentence, "※") {
					continue
				}
				length := utf8.RuneCountInString(sentence)
				if length < 50 || 80 < length {
					continue
				}
				fmt.Println(sentence)
			}
		}
	}
}
