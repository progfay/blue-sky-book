package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

	b "github.com/progfay/blue-sky-book/book"
	l "github.com/progfay/blue-sky-book/line"
	q "github.com/progfay/blue-sky-book/queue"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	targetDir := flag.String("target-dir", "texts", "directory contains aozora book data files")
	min := flag.Int("min", 1000, "minimum length of printing sentences")
	flag.Parse()

	matches, err := filepath.Glob(filepath.Join(*targetDir, "*.txt"))
	if err != nil {
		log.Fatal(err)
	}

	handler := func(path string) {
		book := b.NewBook(path)
		lines, err := book.GetLinesFromBook()
		if err != nil {
			log.Fatal(err)
		}

		sentences := ""
		for _, line := range lines {
			for _, sentence := range l.ParseLine(line) {
				if strings.HasPrefix(sentence, "「") || strings.HasPrefix(sentence, "（") {
					sentences = ""
					continue
				}
				if !strings.HasSuffix(sentence, "。") {
					sentences = ""
					continue
				}
				if strings.Contains(sentence, "※") {
					continue
				}
				sentences += sentence
				if utf8.RuneCountInString(sentences) >= *min {
					fmt.Println(sentences)
					sentences = ""
				}
			}
		}
	}

	queue := q.NewQueue(context.Background(), 100, handler)

	for _, path := range matches {
		queue.Add(path)
	}

	queue.Start()
}
