package asset

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

func UploadFile(file *multipart.FileHeader) (string, error) {
	// upload path
	path := "./static/"

	ext := filepath.Ext(file.Filename)
	randomName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	f, err := os.Create(path + randomName)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer f.Close()
	fileContent, _ := file.Open()
	_, err = io.Copy(f, fileContent)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return "/static/" + randomName, nil
}

func DeleteFile(filename string) error {
	// upload path
	path := "./static/"
	err := os.Remove(path + filename)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
