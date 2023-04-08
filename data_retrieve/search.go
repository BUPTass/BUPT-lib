package data_retrieve

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func SearchAll(keyword string) {}

func SearchNewsCollection(client *mongo.Client, keywords string, c string) ([]byte, error) {
	collection := client.Database("test").Collection(c)
	filter := bson.D{{"$text", bson.D{{"$search", keywords}}}}
	sort := bson.D{{"score", bson.D{{"$meta", "textScore"}}}}
	projection := bson.D{{"title", 1}, {"url", 1}, {"score", bson.D{{"$meta", "textScore"}}}, {"_id", 0}}
	opts := options.Find().SetSort(sort).SetProjection(projection)
	cur, err := collection.Find(context.Background(), filter, opts)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer cur.Close(context.Background())
	/*
		newsList := make([]struct {
			Title string `json:"title"`
			Time  int64  `json:"time"`
			Url   string `json:"url"`
		}, 0, num)
		var itemCount uint
		itemCount = 0
		// iterate through the cursor and decode each document into a news entry struct
		for cur.Next(context.Background()) {
			var news News

			err := cur.Decode(&news)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			itemCount++
			if itemCount < start+1 {
				continue
			}
			if itemCount > num+start {
				break
			}

			reduced := struct {
				Title string `json:"title"`
				Time  int64  `json:"time"`
				Url   string `json:"url"`
			}{
				Title: news.Title,
				Time:  news.Time,
				Url:   news.Url,
			}
			newsList = append(newsList, reduced)
		}
		if err := cur.Err(); err != nil {
			log.Println(err)
			return nil, err
		}
		jsonBytes, err := json.Marshal(newsList)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return jsonBytes, nil

	*/
	return nil, nil
}

func SearchESICollection(client *mongo.Client, keywords string) ([]bson.M, error) {

	collection := client.Database("test").Collection("ESI_paper")

	filter := bson.D{{"$text", bson.D{{"$search", keywords}}}}
	sort := bson.D{{"score", bson.D{{"$meta", "textScore"}}}}
	projection := bson.M{
		"_id": 0,
	}
	opts := options.Find().SetSort(sort).SetProjection(projection)
	cur, err := collection.Find(context.Background(), filter, opts)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer cur.Close(context.Background())

	var papers []bson.M
	var doiS []string
	for cur.Next(context.Background()) {
		var paper bson.M
		if err := cur.Decode(&paper); err != nil {
			return nil, err
		}
		doi := paper["doi"].(string)
		if !contains(doiS, doi) {
			papers = append(papers, paper)
		}
	}
	if err := cur.Err(); err != nil {
		log.Println(err)
		return nil, err
	}
	return papers, nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
