package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
)

func main() {
	par := "client"
	filePath := path.Join(par, "text.txt")
	file, err := os.Open(filePath)
	if err != nil {
		panic(fmt.Errorf("не удалось открыть файл: %w", err))
	}
	defer file.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	part, err := writer.CreateFormFile("file", filePath)
	if err != nil {
		panic(fmt.Errorf("ошибка создания части файла: %w", err))
	}

	_, err = io.Copy(part, file)
	if err != nil {
		panic(fmt.Errorf("ошибка копирования файла: %w", err))
	}

	err = writer.Close()
	if err != nil {
		panic(fmt.Errorf("ошибка закрытия writer: %w", err))
	}

	url := "http://localhost:8080/upload"
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		panic(fmt.Errorf("ошибка создания запроса: %w", err))
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(fmt.Errorf("ошибка выполнения запроса: %w", err))
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Errorf("ошибка чтения ответа: %w", err))
	}

	fmt.Println("Ответ сервера:")
	fmt.Println(string(respBody))
}
