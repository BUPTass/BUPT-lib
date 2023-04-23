package data

import (
	myJson "BUPT-lib/data_retrieve"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"mime/multipart"
	"strconv"
	"time"
)

type News struct {
	Id            string `json:"id"`
	Title         string `json:"title"`
	OutsideSource string `json:"source"`
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
	StartTime    string `bson:"start_time"`
	EndTime      string `bson:"end_time"`
	Link         string `bson:"link"`
	Place        string `bson:"place"`
	Type         string `bson:"type"`
	CcfRank      string `bson:"ccf"`
	Deadline     string `bson:"deadline"`
}

// UpdateNews Not used
func UpdateNews(client *mongo.Client, jsonText string) error {
	collection := client.Database("test").Collection("News")

	var updatedNews News
	err := json.Unmarshal([]byte(jsonText), &updatedNews)
	if err != nil {
		log.Println(err)
		return err
	}

	filter := bson.M{"_id": updatedNews.Id}
	// update := bson.M{"$set": bson.M{"image": imageUri}}
	update := bson.M{"$set": updatedNews}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// GetAnnouncement 基金项目
func GetAnnouncement(client *mongo.Client, num uint, start uint) ([]byte, error) {
	collection := client.Database("test").Collection("Announcement")

	filter := bson.M{}
	opts := options.Find().SetSort(bson.M{"date": -1})

	cur, err := collection.Find(context.Background(), filter, opts)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer cur.Close(context.Background())

	newsList := make([]struct {
		Title string `json:"title"`
		Date  string `json:"date"`
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
		log.Println(err)
		return nil, err
	}
	jsonBytes, err := json.Marshal(newsList)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return jsonBytes, nil
}

// AddAnnouncement 添加基金项目
func AddAnnouncement(client *mongo.Client, news *multipart.FileHeader) error {
	collection := client.Database("test").Collection("Announcement")
	data, err := ParseNewsAnnouncement(news)
	if err != nil {
		log.Println(err)
		return err
	}

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

// ParseNewsAnnouncement 解析基金项目
func ParseNewsAnnouncement(newsFile *multipart.FileHeader) ([]News, error) {
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
	var data []News
	// Loop over CSV records from Row 1
	for _, record := range records[1:] {
		piece := News{
			Id:            "",
			Title:         record[0],
			OutsideSource: "",
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

// GetNews 外部采集新闻--领域新闻
func GetNews(client *mongo.Client, num uint, start uint) ([]byte, error) {
	collection := client.Database("test").Collection("News")

	filter := bson.M{}
	opts := options.Find().SetSort(bson.M{"time": -1})

	cur, err := collection.Find(context.Background(), filter, opts)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer cur.Close(context.Background())

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
}

func AddNews(client *mongo.Client, news *multipart.FileHeader) error {
	collection := client.Database("test").Collection("News")
	data, err := ParseNews(news)
	if err != nil {
		log.Println(err)
		return err
	}
	failedNum := 0
	successfulNum := 0
	for _, data_ := range data {
		_, err = collection.InsertOne(context.Background(), data_)
		if err != nil {
			log.Println(err)
			failedNum++
		} else {
			successfulNum++
		}
	}
	if failedNum > 0 {
		return errors.New("Successfully imported " + strconv.Itoa(successfulNum) +
			" with " + strconv.Itoa(failedNum) + " failed")
	} else {
		return nil
	}
}

// ParseNews 解析外部采集新闻--领域新闻
func ParseNews(newsFile *multipart.FileHeader) ([]News, error) {
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
		return nil, err
	}
	var data []News
	// Loop over CSV records from Row 1
	for _, record := range records[1:] {

		piece := News{
			Id:            "",
			Title:         record[0],
			OutsideSource: record[4],
			Url:           record[1],
			Content:       record[2],
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
		log.Println(err)
		return err
	}

	_, err = collection.DeleteMany(context.Background(), bson.M{})
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

func ParseOngoingConferences(newsFile *multipart.FileHeader) ([]Conference, error) {

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
	var data []Conference
	// Loop over CSV records from Row 1
	for _, record := range records[1:] {

		piece := Conference{
			Abbreviation: record[0],
			Link:         record[1],
			Place:        record[2],
			Deadline:     record[3],
			StartTime:    record[4],
			EndTime:      record[5],
			Type:         record[6],
			CcfRank:      record[7],
		}
		data = append(data, piece)
	}

	return data, nil
}

// GetLibNews 图书馆新闻 & 资源更新
func GetLibNews(client *mongo.Client, num uint, start uint, newsType uint8) ([]byte, error) {
	collection := client.Database("test").Collection("News")

	filter := bson.M{"type": newsType}
	opts := options.Find().SetSort(bson.M{"date": -1})

	cur, err := collection.Find(context.Background(), filter, opts)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer cur.Close(context.Background())

	newsList := make([]struct {
		Title string `json:"title"`
		Date  string `json:"date"`
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
		log.Println(err)
		return nil, err
	}
	jsonBytes, err := json.Marshal(newsList)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return jsonBytes, nil
}

func AddLibNews(client *mongo.Client, news *multipart.FileHeader, newsType uint8) error {
	collection := client.Database("test").Collection("News")
	data, err := ParseLibNews(news, newsType)
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

	return nil
}

// ParseLibNews 图书馆新闻 & 资源更新
func ParseLibNews(newsFile *multipart.FileHeader, newsType uint8) ([]News, error) {
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
	var data []News
	// Loop over CSV records from Row 1
	for _, record := range records[1:] {
		piece := News{
			Id:            "",
			Title:         record[1],
			OutsideSource: "",
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

func CountNews(client *mongo.Client, newsType int) (int64, error) {
	var collection *mongo.Collection
	var filter bson.M
	switch newsType {
	case 1:
		fallthrough
	case 2:
		fallthrough
	case 0:
		collection = client.Database("test").Collection("News")
		filter = bson.M{"type": newsType}
	case 3:
		collection = client.Database("test").Collection("Announcement")
		filter = bson.M{}
	default:
		return 0, nil
	}
	count, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return count, nil
}
