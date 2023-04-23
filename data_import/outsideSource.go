package data

import (
	myJson "BUPT-lib/data_retrieve"
	"context"
	"encoding/csv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"mime/multipart"
)

type EResource struct {
	Name         string `json:"name"`
	Url          string `json:"url"`
	Subject      string `json:"subject"`
	ResourceType string `json:"type"`
	Intro        string `json:"intro"`
}

// GetEResource 电子资源
func GetEResource(client *mongo.Client) ([]byte, error) {
	collection := client.Database("test").Collection("EResource")

	return myJson.GetAllDocumentsAsJson(collection)
}

// AddEResource 添加电子资源
func AddEResource(client *mongo.Client, csv *multipart.FileHeader) error {
	collection := client.Database("test").Collection("EResource")
	data, err := ParseEResource(csv)
	if err != nil {
		log.Println(err)
		return err
	}

	// Delete before insert
	_, err = collection.DeleteMany(context.Background(), bson.M{})
	for _, data_ := range data {
		_, err = collection.InsertOne(context.Background(), data_)
		if err != nil {
			return err
		}
	}

	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// ParseEResource 解析电子资源
func ParseEResource(newsFile *multipart.FileHeader) ([]EResource, error) {
	// Open the Excel file
	tmpFile, err := newsFile.Open()
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
	var data []EResource
	// Loop over CSV records from Row 1
	for _, record := range records[1:] {
		piece := EResource{
			Name:         record[0],
			Url:          record[1],
			Subject:      record[2],
			ResourceType: record[3],
			Intro:        record[4],
		}
		data = append(data, piece)
	}

	return data, nil
}
