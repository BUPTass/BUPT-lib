package hot

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"math"
	"time"
)

const (
	G = 1.8
)

func CalcNewsScore(t int64, hits int64) float64 {
	now := time.Now().Unix()
	interval := now - t
	score := (float64(hits) + 1) / math.Pow(float64(interval)+1, G)
	return score
}

func UpdateNewsScore(client *mongo.Client) error {
	collection := client.Database("test").Collection("News")
	cur, err := collection.Find(context.Background(), bson.M{"valid": true})
	if err != nil {
		log.Println(err)
		return err
	}
	defer cur.Close(context.Background())

	// Iterate through the cursor and update the score of each news document
	for cur.Next(context.Background()) {
		var news struct {
			Id   primitive.ObjectID `bson:"_id"`
			Hits int64              `bson:"hits"`
			Time int64              `bson:"time"`
		}
		err := cur.Decode(&news)
		if err != nil {
			log.Println(err)
			return err
		}
		// Calculate the score based on Time and Type
		score := CalcNewsScore(news.Time, news.Hits)
		// Update the score field in the document
		update := bson.M{"$set": bson.M{"score": score}}
		_, err = collection.UpdateOne(context.Background(), bson.M{"_id": news.Id}, update)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	if err := cur.Err(); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func IncrementNewsHits(client *mongo.Client, id primitive.ObjectID) {
	// Get the news type count from the database
	collection := client.Database("test").Collection("News")
	filter := bson.M{"_id": id}
	update := bson.M{"$inc": bson.M{"hits": 1}}
	opts := options.FindOneAndUpdate().SetUpsert(true)
	collection.FindOneAndUpdate(context.Background(), filter, update, opts)
}
