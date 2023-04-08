package data_retrieve

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func SearchAll(client *mongo.Client, keyword string) ([]byte, error) {
	var searchResult struct {
		ESI     []bson.M
		Found   []bson.M
		LibRes  []bson.M
		LibAnn  []bson.M
		Journal []bson.M
	}
	var err error

	searchResult.ESI, err = searchESICollection(client, keyword)
	if err != nil {
		searchResult.ESI = []bson.M{}
	}
	searchResult.Found, err = searchCollection(client, keyword, "Announcement")
	if err != nil {
		searchResult.Found = []bson.M{}
	}
	searchResult.LibRes, err = searchNewsCollection(client, keyword, "News", 1)
	if err != nil {
		searchResult.LibRes = []bson.M{}
	}
	searchResult.LibAnn, err = searchNewsCollection(client, keyword, "News", 2)
	if err != nil {
		searchResult.LibAnn = []bson.M{}
	}
	searchResult.Journal, err = searchCollection(client, keyword, "CCF")
	if err != nil {
		searchResult.Journal = []bson.M{}
	}

	jsonBytes, err := json.Marshal(searchResult)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return jsonBytes, nil
}

func SearchArticle(client *mongo.Client, keyword string) ([]byte, error) {
	var searchResult struct {
		ESI   []bson.M
		Found []bson.M
		News  []bson.M
	}
	var err error

	searchResult.ESI, err = searchESICollection(client, keyword)
	if err != nil {
		searchResult.ESI = []bson.M{}
	}
	searchResult.Found, err = searchCollection(client, keyword, "Announcement")
	if err != nil {
		searchResult.Found = []bson.M{}
	}
	searchResult.News, err = searchNewsCollection(client, keyword, "News", 0)
	if err != nil {
		searchResult.News = []bson.M{}
	}

	jsonBytes, err := json.Marshal(searchResult)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return jsonBytes, nil
}

func searchNewsCollection(client *mongo.Client, keyword string, colle string, newsType uint8) ([]bson.M, error) {
	collection := client.Database("test").Collection(colle)
	regex := primitive.Regex{Pattern: keyword, Options: "i"}
	filter := bson.M{"title": bson.M{"$regex": regex}, "type": newsType}
	projection := bson.M{
		"_id":   0,
		"title": 1,
		"url":   1,
	}
	opts := options.Find().SetProjection(projection)
	cur, err := collection.Find(context.Background(), filter, opts)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer cur.Close(context.Background())

	var publications []bson.M
	for cur.Next(context.Background()) {
		var publication bson.M
		if err := cur.Decode(&publication); err != nil {
			return nil, err
		}

		publications = append(publications, publication)

	}
	if err := cur.Err(); err != nil {
		log.Println(err)
		return nil, err
	}
	return publications, nil

}

func searchCollection(client *mongo.Client, keyword string, colle string) ([]bson.M, error) {
	collection := client.Database("test").Collection(colle)
	regex := primitive.Regex{Pattern: keyword, Options: "i"}
	filter := bson.M{"title": bson.M{"$regex": regex}}
	projection := bson.M{
		"_id":   0,
		"title": 1,
		"url":   1,
	}
	opts := options.Find().SetProjection(projection)
	cur, err := collection.Find(context.Background(), filter, opts)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer cur.Close(context.Background())

	var publications []bson.M
	for cur.Next(context.Background()) {
		var publication bson.M
		if err := cur.Decode(&publication); err != nil {
			return nil, err
		}

		publications = append(publications, publication)

	}
	if err := cur.Err(); err != nil {
		log.Println(err)
		return nil, err
	}
	return publications, nil

}

func searchESICollection(client *mongo.Client, keywords string) ([]bson.M, error) {

	collection := client.Database("test").Collection("ESI_papers")

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
		doi := paper["DOI"].(string)
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

func SearchTeacher(client *mongo.Client, keyword string) ([]byte, error) {
	collection := client.Database("test").Collection("SCS_Teacher")
	regex := primitive.Regex{Pattern: keyword, Options: "i"}
	filter := bson.M{"name": bson.M{"$regex": regex}}
	projection := bson.M{
		"_id":      0,
		"name":     1,
		"homepage": 1,
	}
	opts := options.Find().SetProjection(projection)
	cur, err := collection.Find(context.Background(), filter, opts)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer cur.Close(context.Background())

	var teachers []bson.M
	for cur.Next(context.Background()) {
		var teacher bson.M
		if err := cur.Decode(&teacher); err != nil {
			return nil, err
		}
		teachers = append(teachers, teacher)

	}
	jsonBytes, err := json.Marshal(teachers)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return jsonBytes, nil

}
