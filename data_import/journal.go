package data

import (
	"BUPT-lib/asset"
	myJson "BUPT-lib/data_retrieve"
	"context"
	"encoding/csv"
	"errors"
	"github.com/tealeg/xlsx"            // library for reading Excel files
	"go.mongodb.org/mongo-driver/bson"  // bson library for marshalling data
	"go.mongodb.org/mongo-driver/mongo" // mongo driver
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"mime/multipart"
	"strconv"
)

type Publication struct {
	ShortName string `bson:"name"`
	FullName  string `bson:"fullName"`
	Publisher string `bson:"publisher"`
	Address   string `bson:"url"`
	Level     string `bson:"level"`
	Type      string `bson:"type"`
	Category  int8   `bson:"category"`
}

func PutT30JournalList(client *mongo.Client, journalXlsFile *multipart.FileHeader) error {
	// Connect to MongoDB
	collection := client.Database("test").Collection("Journals")

	// Open the Excel file
	tmpFile, err := journalXlsFile.Open()
	if err != nil {
		log.Println(err)
		return err
	}
	file, err := xlsx.OpenReaderAt(tmpFile, journalXlsFile.Size)
	if err != nil {
		log.Println(err)
		return err
	}

	// Iterate over the worksheets in the file
	for _, sheet := range file.Sheets {
		// Create a slice to store the data for this worksheet
		data := []interface{}{}

		// Iterate over the rows in the worksheet
		for _, row := range sheet.Rows {
			// Create a map to store the data for this row
			rowData := make(map[string]string)

			// Iterate over the cells in the row
			for idx, cell := range row.Cells {
				// Get the cell value and add it to the row data
				val := cell.String()
				rowData[sheet.Cell(0, idx).String()] = val
			}

			// Add the row data to the sheet data
			if _, err := strconv.Atoi(rowData["序号"]); err == nil {
				data = append(data, rowData)
			}
		}
		if len(data) < 1 {
			log.Println("Invalid Journals File")
			err := errors.New("invalid Journals File")
			return err
		}
		// Delete existing data
		_, err := collection.DeleteMany(context.Background(), bson.M{})

		// Insert the data into MongoDB
		_, err = collection.InsertMany(context.Background(), data)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func SetImage(client *mongo.Client, id string, image *multipart.FileHeader) error {
	// Connect to MongoDB
	collection := client.Database("test").Collection("Journals")
	/*
		imageContent, _ := image.Open()
		imageBytes, err := io.ReadAll(imageContent)
		if err != nil {
			log.Println(err)
			return err
		}
	*/

	imageUri, err := asset.UploadFile(image)
	if err != nil {
		return err
	}

	filter := bson.M{"序号": id}
	update := bson.M{"$set": bson.M{"image": imageUri}}
	// update := bson.M{"$set": bson.M{"image": imageBytes}}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func GetImage(client *mongo.Client, id string) (string, error) {
	// Connect to MongoDB
	collection := client.Database("test").Collection("Journals")

	filter := bson.M{"序号": id}
	projection := bson.M{"_id": 0, "image": 1}
	findOptions := options.FindOneOptions{Projection: projection}

	var result struct {
		Image interface{} `bson:"image"`
	}

	err := collection.FindOne(context.Background(), filter, &findOptions).Decode(&result)
	if err != nil {
		return "", err
	}
	return result.Image.(string), nil

	/*
		var result bson.M
		err := collection.FindOne(context.Background(), filter, &findOptions).Decode(&result)
		if err != nil {
			return nil, errors.New("Invalid Index")
		}
		imageBinary, ok := result["image"].(primitive.Binary)
		if !ok {
			return nil, errors.New("Image Not Found")
		}
		return imageBinary.Data, nil
	*/

}

func GetCcfList(client *mongo.Client) ([]byte, error) {
	collection := client.Database("test").Collection("CCF")
	return myJson.GetAllDocumentsAsJson(collection)
}

func PutCcfList(client *mongo.Client, csv *multipart.FileHeader) error {
	collection := client.Database("test").Collection("CCF")
	data, err := parseCcfCsv(csv)
	if err != nil {
		return err
	}
	collection.DeleteMany(context.Background(), bson.M{})
	for _, data_ := range data {
		_, err = collection.InsertOne(context.Background(), data_)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func parseCcfCsv(csvFile *multipart.FileHeader) ([]Publication, error) {
	// Open the Excel file
	tmpFile, err := csvFile.Open()
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
		return nil, err
	}
	var data []Publication
	// Loop over CSV records from Row 1
	for _, record := range records[1:] {
		i64, err := strconv.ParseInt(record[7], 10, 8)
		if err != nil {
			i64 = -1
		}
		pub := Publication{
			ShortName: record[1],
			FullName:  record[2],
			Publisher: record[3],
			Address:   record[4],
			Level:     record[5],
			Type:      record[6],
			Category:  int8(i64),
		}
		data = append(data, pub)
	}

	return data, nil
}
