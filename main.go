package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Kullanım: go run main.go <girdi_dosyasi> <cikti_dosyasi>")
		return
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	content, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("Dosya okunamadı: %v\n", err)
		return
	}

	processedText := TransformText(string(content))

	err = os.WriteFile(outputFile, []byte(processedText), 0644)
	if err != nil {
		fmt.Printf("Dosya yazılamadı: %v\n", err)
		return
	}

	fmt.Println("İşlem başarıyla tamamlandı!")
}
