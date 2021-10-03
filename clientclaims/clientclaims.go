package clientclaims

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"go.uber.org/zap"
)

// UploadController is the Upload route handler
type UploadController struct {
	logger *zap.Logger
}

// NewUploadController returns a frsh Upload controller
func NewUploadController(logger *zap.Logger) *UploadController {
	return &UploadController{
		logger: logger,
	}
}

// Upload File Handler
func (ctrl *UploadController) UploadFile(rw http.ResponseWriter, r *http.Request) {

	// Create a MultiPart form with size of 128 KB main memory
	err := r.ParseMultipartForm(128 * 1024)
	// if this allocation is not possible or there is error in parsing
	if err != nil {
		ctrl.logger.Warn("Unable to parse request body", zap.Error(err))
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Request Data Faulty"))
		return
	}
	//we use formfile. The request should have the file key with the file as value
	// handler we get here has details about the file
	file, handler, err := r.FormFile("file")
	if err != nil {
		ctrl.logger.Warn("Unable to return file for the provided form key", zap.String("key", "file"), zap.Error(err))
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Request file of nonexistant/incorrect format"))
		return
	}

	//create the file in the directory and copy the file to the folder
	f, err := os.OpenFile("./clientclaims/saveimgdir/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		ctrl.logger.Warn("Unable to create file", zap.String("file", fmt.Sprintf("./clientclaims/saveimgdir/%s", handler.Filename)), zap.Error(err))
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Unable to save the file to the servers disk"))
		return
	}
	io.Copy(f, file)

	ctrl.logger.Info("File Uploaded Successfully", zap.String("file", fmt.Sprintf("./clientclaims/saveimgdir/%s", handler.Filename)))
	// this means file upload successful
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("File Uploaded Successfully"))
}
