package main

import (
	"strconv"
	"strings"
)

// Main function to transform text
func TransformText(text string) string {
	words := strings.Fields(text)
	var result []string

	for i := 0; i < len(words); i++ {
		word := words[i]

		// Hex conversion
		if word == "(hex)" && i > 0 {
			if val, err := strconv.ParseInt(result[len(result)-1], 16, 64); err == nil {
				result[len(result)-1] = strconv.FormatInt(val, 10)
			}
			continue
		}

		// Bin conversion
		if word == "(bin)" && i > 0 {
			if val, err := strconv.ParseInt(result[len(result)-1], 2, 64); err == nil {
				result[len(result)-1] = strconv.FormatInt(val, 10)
			}
			continue
		}

		// Uppercase
		if word == "(up)" && i > 0 {
			result[len(result)-1] = strings.ToUpper(result[len(result)-1])
			continue
		}

		// Lowercase
		if word == "(low)" && i > 0 {
			result[len(result)-1] = strings.ToLower(result[len(result)-1])
			continue
		}

		// Capitalize
		if word == "(cap)" && i > 0 {
			result[len(result)-1] = capitalize(result[len(result)-1])
			continue
		}

		// Numeric parameter tags (e.g., (up, 2))
		if strings.HasPrefix(word, "(") && strings.HasSuffix(word, ")") && strings.Contains(word, ",") {
			cleaned := strings.Trim(word, "()")
			parts := strings.Split(cleaned, ",")
			if len(parts) == 2 {
				tag := strings.TrimSpace(parts[0])
				if count, err := strconv.Atoi(strings.TrimSpace(parts[1])); err == nil {
					for j := 0; j < count && len(result)-1-j >= 0; j++ {
						idx := len(result) - 1 - j
						switch tag {
						case "up":
							result[idx] = strings.ToUpper(result[idx])
						case "low":
							result[idx] = strings.ToLower(result[idx])
						case "cap":
							result[idx] = capitalize(result[idx])
						}
					}
					continue
				}
			}
		}

		// 'a' -> 'an' correction
		if strings.ToLower(word) == "a" && i < len(words)-1 {
			nextWord := words[i+1]
			if len(nextWord) > 0 {
				firstChar := strings.ToLower(string(nextWord[0]))
				if firstChar == "a" || firstChar == "e" || firstChar == "i" || firstChar == "o" || firstChar == "u" || firstChar == "h" {
					if word == "A" {
						word = "An"
					} else {
						word = "an"
					}
				}
			}
		}

		result = append(result, word)
	}

	finalText := strings.Join(result, " ")
	return fixPunctuation(finalText)
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	runes := []rune(s)
	if runes[0] >= 'a' && runes[0] <= 'z' {
		runes[0] -= 32
	}
	return string(runes)
}

func fixPunctuation(text string) string {
	for _, p := range []string{".", ",", "!", "?", ":", ";"} {
		text = strings.ReplaceAll(text, " "+p, p)
	}
	text = strings.ReplaceAll(text, "‘ ", "‘")
	text = strings.ReplaceAll(text, " ’", "’")
	text = strings.ReplaceAll(text, "' ", "'")
	text = strings.ReplaceAll(text, " '", "'")
	return text
}
