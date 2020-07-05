package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	DeliminatorRegexp = regexp.MustCompile(`^-+$`)
	AnnotationRegexp  = regexp.MustCompile(`［.+?］`)
	FuriganaRegexp    = regexp.MustCompile(`《[^》]+》`)
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GetLinesFromBook(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := make([]string, 0)

	breakLineCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if DeliminatorRegexp.MatchString(line) {
			break
		}
	}

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if DeliminatorRegexp.MatchString(line) {
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

		line = AnnotationRegexp.ReplaceAllString(line, "")
		line = FuriganaRegexp.ReplaceAllString(line, "")
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

func parseLine(line string) []string {
	lines := make([]string, 0)
	chars := make([]rune, 0)
	inAngledBracket := false
	inRoundedBracket := false
	inSentence := false
	for _, char := range []rune(line) {
		chars = append(chars, char)

		switch char {
		case '「':
			inAngledBracket = true

		case '」':
			inAngledBracket = false
			if !inSentence {
				lines = append(lines, strings.TrimSpace(string(chars)))
				chars = make([]rune, 0)
			}

		case '(', '（':
			inRoundedBracket = true

		case ')', '）':
			inRoundedBracket = false
			if !inSentence {
				lines = append(lines, strings.TrimSpace(string(chars)))
				chars = make([]rune, 0)
			}

		case '!', '！', '?', '？', '。':
			if inSentence && !inAngledBracket && !inRoundedBracket {
				lines = append(lines, strings.TrimSpace(string(chars)))
				chars = make([]rune, 0)
			}

		default:
			if !inAngledBracket && !inRoundedBracket {
				inSentence = true
			}
		}
	}

	if len(chars) > 0 {
		lines = append(lines, strings.TrimSpace(string(chars)))
	}

	return lines
}

func main() {
	matches, err := filepath.Glob("texts/*.txt")
	if err != nil {
		log.Fatal(err)
	}

	for _, path := range matches {

		lines, err := GetLinesFromBook(path)
		if err != nil {
			log.Fatal(err)
		}

		for _, line := range lines {
			for _, sentence := range parseLine(line) {
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
