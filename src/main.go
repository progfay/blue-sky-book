package main

import (
	"fmt"
	"log"
	"flag"
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
	targetDir := flag.String("target-dir", "texts", "directory contains aozora book data files")
	min := flag.Int("min", 50, "minimum length of sentence to extract")
	max := flag.Int("max", 80, "maximum length of sentence to extract")
	flag.Parse()

	matches, err := filepath.Glob(filepath.Join(*targetDir, "*.txt"))
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
				if length < *min || *max < length {
					continue
				}
				fmt.Println(sentence)
			}
		}
	}
}
