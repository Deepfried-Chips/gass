package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"net/url"
	"path"
)

func contains(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}

	return false
}

func validateFileType(fileBytes *[]byte) (string, error) {

	// check file type, detectcontenttype only needs the first 512 bytes
	filetype := http.DetectContentType(*fileBytes)
	fileEndings, err := mime.ExtensionsByType(filetype)
	if err != nil {
		return "", &fileTypeReadError
	}

	return fileEndings[0], nil
}

func randomToken(len int) string {
	b := make([]byte, len)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", b)
}

// createResponse returns a JSON object for ShareX upload requests.
func createResponse(prefix string, file string) []byte {

	u, _ := url.Parse("http://" + config.Host)
	u.Path = path.Join(u.Path, prefix, file)

	res := response{
		Status: http.StatusOK,
		Result: result{
			Name: file,
			URL:  u.String(),
			// doesnt do anything yet.
			//DeleteKey:  randomToken(6)},
		},
	}

	j, _ := json.Marshal(&res)
	return j

}
