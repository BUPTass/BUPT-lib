package data

import (
	"context"
	"errors"
	"github.com/tealeg/xlsx"           // library for reading excel files
	"go.mongodb.org/mongo-driver/bson" // bson library for marshalling data
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo" // mongo driver
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"log"
	"mime/multipart"
	"strconv"
)

func UpdateJournalList(client *mongo.Client, journalXlsFile *multipart.FileHeader) error {
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
			err := errors.New("Invalid Journals File")
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
	imageContent, _ := image.Open()
	imageBytes, err := io.ReadAll(imageContent)
	if err != nil {
		log.Println(err)
		return err
	}

	// TODO: Store images another way
	/*
		imageUri, err = Asset.UploadImage(image)
		if err != nil {
			return err
		}
	*/

	filter := bson.M{"序号": id}
	// update := bson.M{"$set": bson.M{"image": imageUri}}
	update := bson.M{"$set": bson.M{"image": imageBytes}}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func GetImage(client *mongo.Client, id string) ([]byte, error) {
	// Connect to MongoDB
	collection := client.Database("test").Collection("Journals")

	filter := bson.M{"序号": id}
	projection := bson.M{"_id": 0, "image": 1}
	findOptions := options.FindOneOptions{Projection: projection}

	// TODO: Store images another way
	/*
		var imageUri string
		err := collection.FindOne(context.Background(), filter, &findOptions).Decode(&imageUri)
		if err != nil {
			return nil, err
		}
		return imageUri, nil
	*/

	var result bson.M
	err := collection.FindOne(context.Background(), filter, &findOptions).Decode(&result)
	if err != nil {
		return nil, errors.New("Invaild Index")
	}
	imageBinary, ok := result["image"].(primitive.Binary)
	if !ok {
		return nil, errors.New("Image Not Found")
	}
	return imageBinary.Data, nil
}
