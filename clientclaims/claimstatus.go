package clientclaims

import (
	"compress/gzip"
	"io/ioutil"
	"net/http"
	"strings"
)

func DownloadFile(rw http.ResponseWriter, r *http.Request) {
	if _, ok := r.Header["Email"]; !ok {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Email Missing"))
		return
	}

	fileName := strings.Replace(r.Header["Email"][0], ".", "_", -1)
	files, err := ioutil.ReadDir("./clientclaims/claimstatusdir/")
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Unable to read the claim directory"))
		return
	}

	for _, claim := range files {
		if strings.Contains(claim.Name(), fileName) {
			writer := gzip.NewWriter(rw)
			writer.Name = "claim-status.jpg"
			writer.Comment = "Status of your Claim from the Accounting department"

			fileContent, err := ioutil.ReadFile(claim.Name())
			if err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				rw.Write([]byte("Unable to read the claim file"))

				err := closeGzipStream(writer)
				if err != nil {
					rw.WriteHeader(http.StatusInternalServerError)
					rw.Write([]byte("Unable to read the claim file"))
					return
				}
				return
			}
			_, err = writer.Write(fileContent)

			if err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				rw.Write([]byte("Error writing the Gzipped File"))

				err := closeGzipStream(writer)
				if err != nil {
					rw.WriteHeader(http.StatusInternalServerError)
					rw.Write([]byte("Unable to read the claim file"))
					return
				}
				return
			}
			err = closeGzipStream(writer)
			if err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				rw.Write([]byte("Unable to read the claim file"))
				return
			}
			rw.WriteHeader(http.StatusOK)
			rw.Write([]byte("File Gzipped and Downloaded"))
		}
	}
	rw.WriteHeader(http.StatusLocked)
	rw.Write([]byte("Claim not yet generated"))
}

func closeGzipStream(writer *gzip.Writer) error {
	err := writer.Flush()
	if err != nil {
		return err
	}
	err = writer.Close()
	if err != nil {
		return err
	}
	return nil
}
