package data

import (
	"context"
	"encoding/csv"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"mime/multipart"
)

type teacher struct {
	Name     string `json:"name"`
	Homepage string `json:"url"`
}

func AddTeacher(client *mongo.Client, csv *multipart.FileHeader) error {
	collection := client.Database("test").Collection("SCS_Teacher")
	data, err := ParseTeacher(csv)
	if err != nil {
		log.Println(err)
		return err
	}

	for _, data_ := range data {
		_, err = collection.InsertOne(context.Background(), data_)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func ParseTeacher(newFile *multipart.FileHeader) ([]teacher, error) {
	// Open the Excel file
	tmpFile, err := newFile.Open()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// Parse CSV file
	reader := csv.NewReader(tmpFile)
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1
	records, err := reader.ReadAll()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var data []teacher
	// Loop over CSV records from Row 1
	for _, record := range records[1:] {
		piece := teacher{
			Name:     record[0],
			Homepage: record[1],
		}
		data = append(data, piece)
	}

	return data, nil
}
