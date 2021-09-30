package clientclaims

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// Upload File Handler
func UploadFile(rw http.ResponseWriter, r *http.Request) {

	// Create a MultiPart form with size of 128 KB main memory
	err := r.ParseMultipartForm(128 * 1024)
	// if this allocation is not possible or there is error in parsing
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Request Data Faulty"))
		fmt.Println(err)
		return
	}
	//we use formfile. The request should have the file key with the file as value
	// handler we get here has details about the file
	file, handler, err := r.FormFile("file")
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Request file of nonexistant/incorrect format"))
		fmt.Println(err)
		return
	}

	//create the file in the directory and copy the file to the folder
	f, err := os.OpenFile("./clientclaims/saveimgdir/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Unable to save the file to the servers disk"))
		return
	}
	io.Copy(f, file)

	// this means file upload successful
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("File Uploaded Successfully"))
}
