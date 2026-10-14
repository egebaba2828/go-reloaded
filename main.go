package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Kullanım: go run main.go <girdi_dosyasi> <cikti_dosyasi>")
		return
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	content, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("Dosya okuma hatası: %v\n", err)
		return
	}

	processedText := processText(string(content))

	err = os.WriteFile(outputFile, []byte(processedText), 0644)
	if err != nil {
		fmt.Printf("Dosya yazma hatası: %v\n", err)
		return
	}

	fmt.Println("İşlem tamam! Düzenlenen metin yazıldı:", outputFile)
}

func processText(text string) string {
	words := strings.Fields(text)
	var result []string

	for i := 0; i < len(words); i++ {
		w := words[i]

		// 1. Hexadecimal dönüşümü: (hex)
		if w == "(hex)" && len(result) > 0 {
			val, err := strconv.ParseInt(result[len(result)-1], 16, 64)
			if err == nil {
				result[len(result)-1] = fmt.Sprintf("%d", val)
			}
			continue
		}

		// 2. Binary dönüşümü: (bin)
		if w == "(bin)" && len(result) > 0 {
			val, err := strconv.ParseInt(result[len(result)-1], 2, 64)
			if err == nil {
				result[len(result)-1] = fmt.Sprintf("%d", val)
			}
			continue
		}

		// 3. Tekli harf düzenlemeleri: (up), (low), (cap)
		if (w == "(up)" || w == "(low)" || w == "(cap)") && len(result) > 0 {
			if w == "(up)" {
				result[len(result)-1] = strings.ToUpper(result[len(result)-1])
			} else if w == "(low)" {
				result[len(result)-1] = strings.ToLower(result[len(result)-1])
			} else if w == "(cap)" {
				result[len(result)-1] = capitalize(result[len(result)-1])
			}
			continue
		}

		// 4. Parametreli harf düzenlemeleri: (up, N), (low, N), (cap, N)
		if strings.HasPrefix(w, "(") && (strings.Contains(w, "up") || strings.Contains(w, "low") || strings.Contains(w, "cap")) {
			// Örn: "(up," veya "(up, 2)" durumunu ayıkla
			op, n := parseCommand(w, words, &i)
			applyCommand(result, op, n)
			continue
		}

		// 5. a -> an dönüşümü
		if strings.ToLower(w) == "a" && i+1 < len(words) {
			nextWord := strings.ToLower(words[i+1])
			if len(nextWord) > 0 {
				firstChar := nextWord[0]
				if firstChar == 'a' || firstChar == 'e' || firstChar == 'i' || firstChar == 'o' || firstChar == 'u' || firstChar == 'h' {
					if w == "A" {
						result = append(result, "An")
					} else {
						result = append(result, "an")
					}
					continue
				}
			}
		}

		result = append(result, w)
	}

	// Kelimeleri birleştirip noktalama kurallarını uygulayalım
	joined := strings.Join(result, " ")
	return fixPunctuation(joined)
}

// Komut ve sayı parametresini okuyan yardımcı fonksiyon
func parseCommand(w string, words []string, i *int) (string, int) {
	// Eğer etiket tek kelime içinde bittiyse örn: (up,2)
	cleaned := strings.Trim(w, "()")
	parts := strings.Split(cleaned, ",")
	if len(parts) == 2 {
		op := strings.TrimSpace(parts[0])
		n, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err == nil {
			return op, n
		}
	}
	// Eğer arada boşluk kaldıysa örn: (up, 2) -> sonraki kelimeye bak
	if *i+1 < len(words) && strings.HasSuffix(words[*i+1], ")") {
		numStr := strings.TrimSuffix(words[*i+1], ")")
		n, err := strconv.Atoi(numStr)
		if err == nil {
			*i++ // sayı kelimesini de atla
			return strings.Trim(parts[0], ","), n
		}
	}
	return "up", 1
}

// Son N kelimeye ilgili komutu uygulayan fonksiyon
func applyCommand(result []string, op string, n int) {
	start := len(result) - n
	if start < 0 {
		start = 0
	}
	for j := start; j < len(result); j++ {
		if op == "up" {
			result[j] = strings.ToUpper(result[j])
		} else if op == "low" {
			result[j] = strings.ToLower(result[j])
		} else if op == "cap" {
			result[j] = capitalize(result[j])
		}
	}
}

func capitalize(s string) string {
	if s == "" {
		return ""
	}
	runes := []rune(s)
	if runes[0] >= 'a' && runes[0] <= 'z' {
		runes[0] -= 32
	}
	return string(runes)
}

// Noktalama işaretleri ve tırnak boşluklarını düzelten fonksiyon
func fixPunctuation(text string) string {
	// Basit noktalama yapıştırma kuralları
	for _, p := range []string{".", ",", "!", "?", ":", ";"} {
		text = strings.ReplaceAll(text, " "+p, p)
	}
	// Tekli tırnak düzeltmeleri
	text = strings.ReplaceAll(text, " ' ", " '")
	text = strings.ReplaceAll(text, "' ", "'")
	return text
}
