package workers

import (
	"corpPR3/internal/models"
	"corpPR3/internal/services"
	"mime/multipart"
	"path/filepath"
	"sync/atomic"
)

var fileID int64 = 0

var TaskQueue = make(chan Task, 100)

func StartWorkerPool(num int) {
	for i := 0; i < num; i++ {
		go worker()
	}
}

func worker() {
	for task := range TaskQueue {
		result, err := services.ProcessFile(task.FileHeader, task.ID)

		if err != nil {
			task.ErrorChan <- err
		} else {
			task.ResponseChan <- result
		}
	}
}

func NewTask(fileHeader *multipart.FileHeader) Task {
	id := atomic.AddInt64(&fileID, 1)
	return Task{
		FileHeader:   fileHeader,
		OriginalName: filepath.Base(fileHeader.Filename),
		ID:           id,
		ResponseChan: make(chan models.AnalysisResult, 1),
		ErrorChan:    make(chan error, 1),
	}
}
