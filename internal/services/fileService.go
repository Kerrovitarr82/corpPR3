package services

import (
	"bufio"
	"corpPR3/internal/models"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

func ProcessFile(fileHeader *multipart.FileHeader, fileID int64) (models.AnalysisResult, error) {
	result := models.AnalysisResult{
		Filename: "",
		Lines:    0,
		Words:    0,
		Chars:    0,
	}
	originalName := fileHeader.Filename
	result.Filename = originalName

	src, err := fileHeader.Open()
	if err != nil {
		return result, fmt.Errorf("не удалось открыть оригинальный файл: %w", err)
	}
	defer func() {
		if cerr := src.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("ошибка закрытия оригинального файла: %w", cerr)
		}
	}()

	// Сохранение оригинального файла
	receivedName := fmt.Sprintf("received_%d_%s", fileID, originalName)
	receivedPath := filepath.Join("uploads", receivedName)
	outFile, err := os.Create(receivedPath)
	if err != nil {
		return result, fmt.Errorf("не удалось сохранить файл upload: %w", err)
	}
	defer func() {
		if cerr := outFile.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("ошибка закрытия файла upload: %w", cerr)
		}
	}()

	_, err = io.Copy(outFile, src)
	if err != nil {
		return result, fmt.Errorf("ошибка при копировании данных из upload файла: %w", err)
	}

	_, err = outFile.Seek(0, io.SeekStart)
	if err != nil {
		return result, fmt.Errorf("ошибка при установке курсора в начало у upload файла: %w", err)
	}

	reader := bufio.NewReader(outFile)
	for {
		line, err := reader.ReadString('\n')
		result.Lines++
		result.Words += len(strings.Fields(line))
		result.Chars += len([]rune(line))
		if err == io.EOF {
			break
		}
		if err != nil {
			return result, fmt.Errorf("ошибка чтения файла при анализе: %w", err)
		}
	}

	analysisName := fmt.Sprintf("analysis_result_%d_%s", fileID, originalName)
	analysisPath := filepath.Join("analysis_results", analysisName)
	analysisFile, err := os.Create(analysisPath)
	if err != nil {
		return result, fmt.Errorf("не удалось создать файл анализа: %w", err)
	}
	defer func() {
		if cerr := analysisFile.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("ошибка закрытия файла анализа: %w", cerr)
		}
	}()

	_, err = fmt.Fprintf(analysisFile, "Имя оригинального файла: %s\nСтрок: %d, Слов: %d, Символов: %d\n", result.Filename, result.Lines, result.Words, result.Chars)
	if err != nil {
		return result, fmt.Errorf("ошибка записи данных в файл анализа: %w", err)
	}

	return result, nil
}
