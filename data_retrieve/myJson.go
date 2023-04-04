package data_retrieve

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sort"
	"strconv"
)

func GetAllDocumentsAsJson(collection *mongo.Collection) ([]byte, error) {
	// Define options to limit the number of returned documents.
	// To retrieve all documents, set the limit to -1.
	findOptions := options.Find()
	//findOptions.SetLimit(-1)

	// Retrieve all documents from the collection.
	cursor, err := collection.Find(context.Background(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer cursor.Close(context.Background())

	// Iterate through the cursor and decode each document into a map[string]interface{}.
	var documents []map[string]interface{}
	for cursor.Next(context.Background()) {
		var document map[string]interface{}
		if err := cursor.Decode(&document); err != nil {
			log.Fatal(err)
			return nil, err
		}
		documents = append(documents, document)
	}
	if documents[0]["序号"] != nil {
		sortMapsByIndex(documents, "序号")
	}
	// Marshal the documents into JSON.
	jsonBytes, err := json.Marshal(documents)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return jsonBytes, nil
}

func sortMapsByIndex(maps []map[string]interface{}, index string) {
	sort.Slice(maps, func(i, j int) bool {
		a, _ := strconv.Atoi(maps[i][index].(string))
		b, _ := strconv.Atoi(maps[j][index].(string))
		return a < b
	})
}
