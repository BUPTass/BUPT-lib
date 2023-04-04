package data

import (
	myJson "BUPT-lib/data_retrieve"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"mime/multipart"
	"time"
)

type News struct {
	Id            string `json:"id"`
	Title         string `json:"title"`
	isFromOutside bool   `json:"outside"`
	Url           string `json:"url"`
	Content       string `json:"content"`
	Date          string `json:"date"`
	Time          int64  `json:"time"`
	CreateTime    int64  `json:"create_time"`
	UpdateTime    int64  `json:"update_time"`
	Type          uint8  `json:"type"`
	isVaild       bool   `json:"valid"`
}

type Conference struct {
	Abbreviation string `bson:"abbreviation"`
	FullName     string `bson:"full_name"`
	StartTime    string `bson:"start_time"`
	EndTime      string `bson:"end_time"`
	Link         string `bson:"link"`
	Address      string `bson:"address"`
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

func GetAnnouncement(client *mongo.Client, num int) ([]byte, error) {
	collection := client.Database("test").Collection("Announcement")

	filter := bson.M{}
	opts := options.Find().SetSort(bson.M{"date": -1})

	if !(num > 0 && num < 50) {
		num = 50
	}
	// set the limit for number of news entries to retrieve
	opts.SetLimit(int64(num))

	cur, err := collection.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	newsList := make([]struct {
		Title string `json:"title"`
		Date  string `json:"date"`
		Url   string `json:"url"`
	}, 0, num)

	// iterate through the cursor and decode each document into a news entry struct
	for cur.Next(context.Background()) {
		var news News

		err := cur.Decode(&news)
		if err != nil {
			return nil, err
		}
		reduced := struct {
			Title string `json:"title"`
			Date  string `json:"date"`
			Url   string `json:"url"`
		}{
			Title: news.Title,
			Date:  news.Date,
			Url:   news.Url,
		}
		newsList = append(newsList, reduced)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	jsonBytes, err := json.Marshal(newsList)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return jsonBytes, nil
}

func AddAnnouncement(client *mongo.Client, news *multipart.FileHeader) error {
	collection := client.Database("test").Collection("Announcement")
	data, err := ParseNewsAnnouncement(news)
	if err != nil {
		return err
	}

	for _, data_ := range data {
		_, err = collection.InsertOne(context.Background(), data_)
		if err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}
	return nil
}

func ParseNewsAnnouncement(newsFile *multipart.FileHeader) ([]News, error) {
	// Open the Excel file
	tmpFile, err := newsFile.Open()
	if err != nil {
		log.Fatal(err)
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
	var data []News
	// Loop over CSV records from Row 1
	for _, record := range records[1:] {
		piece := News{
			Id:            "",
			Title:         record[0],
			isFromOutside: true,
			Url:           record[2],
			Content:       "",
			Date:          record[1],
			Time:          0,
			CreateTime:    time.Now().Unix(),
			UpdateTime:    time.Now().Unix(),
			Type:          0,
			isVaild:       true,
		}
		data = append(data, piece)
	}

	return data, nil
}

func GetNews(client *mongo.Client, num int) ([]byte, error) {
	collection := client.Database("test").Collection("News")

	filter := bson.M{}
	opts := options.Find().SetSort(bson.M{"time": -1})

	if !(num > 0 && num < 50) {
		num = 50
	}

	// set the limit for number of news entries to retrieve
	opts.SetLimit(int64(num))

	cur, err := collection.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	newsList := make([]struct {
		Title string `json:"title"`
		Time  int64  `json:"time"`
		Url   string `json:"url"`
	}, 0, num)

	// iterate through the cursor and decode each document into a news entry struct
	for cur.Next(context.Background()) {
		var news News

		err := cur.Decode(&news)
		if err != nil {
			return nil, err
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
		return nil, err
	}
	jsonBytes, err := json.Marshal(newsList)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return jsonBytes, nil
}

func AddNews(client *mongo.Client, news *multipart.FileHeader) error {
	collection := client.Database("test").Collection("News")
	data, err := ParseNews(news)
	if err != nil {
		return err
	}

	for _, data_ := range data {
		_, err = collection.InsertOne(context.Background(), data_)
		if err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}
	return nil
}

func ParseNews(newsFile *multipart.FileHeader) ([]News, error) {
	// Open the Excel file
	tmpFile, err := newsFile.Open()
	if err != nil {
		log.Fatal(err)
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
	var data []News
	// Loop over CSV records from Row 1
	for _, record := range records[1:] {

		piece := News{
			Id:            "",
			Title:         record[0],
			isFromOutside: true,
			Url:           record[1],
			Content:       "",
			Date:          "",
			Time:          castRFC1123(record[3]),
			CreateTime:    time.Now().Unix(),
			UpdateTime:    time.Now().Unix(),
			Type:          0,
			isVaild:       true,
		}
		data = append(data, piece)
	}

	return data, nil
}

func castRFC1123(dateString string) int64 {
	t, err := time.Parse(time.RFC1123, dateString)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return 0
	}

	// Convert the time.Time value to a Unix timestamp.
	return t.Unix()
}

func GetOngoingConferences(client *mongo.Client) ([]byte, error) {
	collection := client.Database("test").Collection("OngoingConf")
	return myJson.GetAllDocumentsAsJson(collection)
}

func AddOngoingConferences(client *mongo.Client, conf *multipart.FileHeader) error {
	collection := client.Database("test").Collection("OngoingConf")
	data, err := ParseOngoingConferences(conf)
	if err != nil {
		return err
	}

	_, err = collection.DeleteMany(context.Background(), bson.M{})
	for _, data_ := range data {
		_, err = collection.InsertOne(context.Background(), data_)
		if err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}
	return nil
}

func ParseOngoingConferences(newsFile *multipart.FileHeader) ([]Conference, error) {

	tmpFile, err := newsFile.Open()
	if err != nil {
		log.Fatal(err)
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
	var data []Conference
	// Loop over CSV records from Row 1
	for _, record := range records[1:] {

		piece := Conference{
			Abbreviation: record[0],
			FullName:     record[1],
			StartTime:    record[2],
			EndTime:      record[3],
			Link:         record[4],
			Address:      record[5],
		}
		data = append(data, piece)
	}

	return data, nil
}

func GetLibNews(client *mongo.Client, num int, newsType uint8) ([]byte, error) {
	collection := client.Database("test").Collection("News")

	filter := bson.M{"type": newsType}
	opts := options.Find().SetSort(bson.M{"date": -1})

	if !(num > 0 && num < 50) {
		num = 50
	}
	// set the limit for number of news entries to retrieve
	opts.SetLimit(int64(num))

	cur, err := collection.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	newsList := make([]struct {
		Title string `json:"title"`
		Date  string `json:"date"`
		Url   string `json:"url"`
	}, 0, num)

	// iterate through the cursor and decode each document into a news entry struct
	for cur.Next(context.Background()) {
		var news News

		err := cur.Decode(&news)
		if err != nil {
			return nil, err
		}
		reduced := struct {
			Title string `json:"title"`
			Date  string `json:"date"`
			Url   string `json:"url"`
		}{
			Title: news.Title,
			Date:  news.Date,
			Url:   news.Url,
		}
		newsList = append(newsList, reduced)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	jsonBytes, err := json.Marshal(newsList)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return jsonBytes, nil
}

func AddLibNews(client *mongo.Client, news *multipart.FileHeader, newsType uint8) error {
	collection := client.Database("test").Collection("News")
	data, err := ParseLibNews(news, newsType)
	if err != nil {
		return err
	}

	for _, data_ := range data {
		_, err = collection.InsertOne(context.Background(), data_)
		if err != nil {
			return err
		}
	}

	return nil
}

func ParseLibNews(newsFile *multipart.FileHeader, newsType uint8) ([]News, error) {
	// Open the Excel file
	tmpFile, err := newsFile.Open()
	if err != nil {
		log.Fatal(err)
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
	var data []News
	// Loop over CSV records from Row 1
	for _, record := range records[1:] {
		piece := News{
			Id:            "",
			Title:         record[1],
			isFromOutside: true,
			Url:           record[3],
			Content:       "",
			Date:          record[2],
			Time:          0,
			CreateTime:    time.Now().Unix(),
			UpdateTime:    time.Now().Unix(),
			Type:          uint8(newsType),
			isVaild:       true,
		}
		data = append(data, piece)
	}

	return data, nil
}
