package data

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type News struct {
	Id            string `json:"_id"`
	Title         string `json:"title"`
	isFromOutside bool   `json:"outside"`
	Url           string `json:"url""`
	Content       string `json:"content"`
	CreateTime    int64  `json:"create_time"`
	UpdateTime    int64  `json:"update_time"`
	Type          string `json:"type"`
	isVaild       bool   `json:"valid"`
}

func UpdateNews(client *mongo.Client, jsonText string) error {
	collection := client.Database("test").Collection("News")

	var updatedNews News
	err := json.Unmarshal([]byte(jsonText), &updatedNews)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": updatedNews.Id}
	// update := bson.M{"$set": bson.M{"image": imageUri}}
	update := bson.M{"$set": updatedNews}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}
