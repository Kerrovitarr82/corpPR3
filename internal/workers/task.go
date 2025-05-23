package workers

import (
	"corpPR3/internal/models"
	"mime/multipart"
)

type Task struct {
	FileHeader   *multipart.FileHeader
	OriginalName string
	ID           int64
	ResponseChan chan models.AnalysisResult
	ErrorChan    chan error
}
