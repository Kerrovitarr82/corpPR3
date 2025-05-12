package main

import (
	"bufio"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func analyzeText(filePath string) (int, int, int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, 0, 0, err
	}
	defer file.Close()

	var lines, words, chars int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lines++
		words += len(strings.Fields(line))
		chars += len(line)
	}

	return lines, words, chars, nil
}

func uploadHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "Файл не получен"})
		return
	}

	timestamp := time.Now().Format("20060102_150405")
	receivedPath := filepath.Join("received", fmt.Sprintf("received_%s_%s", timestamp, file.Filename))

	// Сохраняем файл
	if err := c.SaveUploadedFile(file, receivedPath); err != nil {
		c.JSON(500, gin.H{"error": "Не удалось сохранить файл"})
		return
	}

	// Анализ
	lines, words, chars, err := analyzeText(receivedPath)
	if err != nil {
		c.JSON(500, gin.H{"error": "Ошибка анализа текста"})
		return
	}

	// Сохраняем результаты
	analysisPath := filepath.Join("analysis", fmt.Sprintf("analysis_%s.txt", timestamp))
	resultFile, _ := os.Create(analysisPath)
	defer resultFile.Close()

	analysis := fmt.Sprintf("Lines: %d\nWords: %d\nCharacters: %d\n", lines, words, chars)
	_, err = resultFile.WriteString(analysis)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Отправляем JSON-ответ
	c.JSON(200, gin.H{
		"lines":      lines,
		"words":      words,
		"characters": chars,
	})
}

func main() {
	r := gin.Default()
	r.POST("/upload", uploadHandler)
	err := r.Run(":8080")
	if err != nil {
		fmt.Println(err)
		return
	} // Сервер на http://localhost:8080
}
