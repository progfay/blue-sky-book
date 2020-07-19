package job

import (
	"context"
	"fmt"
	"log"
	"strings"
	"unicode/utf8"

	b "github.com/progfay/blue-sky-book/book"
	l "github.com/progfay/blue-sky-book/line"
)

type Job struct {
	Min  int
	Max  int
	Path string
}

func (job Job) Run(ctx context.Context) {
	book := b.NewBook(job.Path)
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
			if length < job.Min || job.Max < length {
				continue
			}
			fmt.Println(sentence)
		}
	}
}
