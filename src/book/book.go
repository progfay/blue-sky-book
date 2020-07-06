package book

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

var (
	deliminatorRegexp = regexp.MustCompile(`^-+$`)
	annotationRegexp  = regexp.MustCompile(`［.+?］`)
	furiganaRegexp    = regexp.MustCompile(`《.+?》`)
)

type Book struct {
	Path string
}

func NewBook(path string) *Book {
	return &Book{
		Path: path,
	}
}

func (b *Book) GetLinesFromBook() ([]string, error) {
	file, err := os.Open(b.Path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	lines := make([]string, 0)

	breakLineCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if deliminatorRegexp.MatchString(line) {
			break
		}
	}

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if deliminatorRegexp.MatchString(line) {
			break
		}
	}

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if line == "" {
			breakLineCount++
			if breakLineCount >= 3 {
				break
			}
			continue
		}
		breakLineCount = 0

		line = annotationRegexp.ReplaceAllString(line, "")
		line = furiganaRegexp.ReplaceAllString(line, "")
		line = strings.ReplaceAll(line, "｜", "")

		if line == "" {
			continue
		}

		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}
