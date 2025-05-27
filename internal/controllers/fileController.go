package controllers

import (
	"corpPR3/internal/workers"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HandleFileUpload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Файл не был получен"})
		return
	}

	task := workers.NewTask(file)
	workers.TaskQueue <- task

	select {
	case res := <-task.ResponseChan:
		c.JSON(http.StatusOK, res)
	case err := <-task.ErrorChan:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
