package transport

import (
	"corpPR3/internal/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.POST("/upload", controllers.HandleFileUpload)

	return router
}
