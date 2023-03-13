package asset

import (
	"io"
	"mime/multipart"
	"os"
)

func UploadFile(file *multipart.FileHeader) (string, error) {
	// upload path
	path := "./"

	f, err := os.Create(path + file.Filename)
	if err != nil {
		return "", err
	}
	defer f.Close()
	fileContent, _ := file.Open()
	_, err = io.Copy(f, fileContent)
	if err != nil {
		return "", err
	}
	return path + file.Filename, nil
}

func DeleteFile(filename string) error {
	// upload path
	path := "./"
	err := os.Remove(path + filename)
	if err != nil {
		return err
	}
	return nil
}
