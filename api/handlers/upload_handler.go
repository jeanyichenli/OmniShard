package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jeanyichenli/FileUploadSystem/internal/upload"
)

type UploadResponse struct {
	ID     string `json:"uploadid"`
	Status string `json:"status"`
}

// UploadSession is a local alias for the upload domain session type.
type UploadSession = upload.Session

var uploadService upload.Service = upload.NewService()

func Upload(c *gin.Context) {
	ses := UploadSession{}

	if err := c.ShouldBindJSON(&ses); err != nil {
		c.JSON(http.StatusInternalServerError, errParsingJson)
		return
	}

	ses = uploadService.InitiateSession(ses)

	// TODO: other processes

	// Create response
	res := UploadResponse{
		ID:     ses.ID,
		Status: ses.Status,
	}

	c.JSON(http.StatusOK, res)
}
