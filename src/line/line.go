package line

import "strings"

func ParseLine(line string) []string {
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
