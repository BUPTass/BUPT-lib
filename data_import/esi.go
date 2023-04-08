package data

import (
	Asset "BUPT-lib/asset"
	myJson "BUPT-lib/data_retrieve"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"github.com/tealeg/xlsx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"mime/multipart"
	"path"
	"strconv"
	"strings"
)

type ESIHot struct {
	AccessionNumber string
	DOI             string
	PMID            string
	ArticleName     string
	Authors         []string
	Source          string
	ResearchField   string
	TimesCited      int
	Countries       []string
	Addresses       []string
	Institutions    []string
	PublicationDate int
}

type IncitesRow struct {
	Author              string  `bson:"名称"`
	WebOfSciencePapers  int     `bson:"Web of Science 论文数"`
	Rank                int     `bson:"排名"`
	CitationFrequency   int     `bson:"被引频次"`
	Researchers         string  `bson:"研究人员 Web of Science ResearcherID"`
	CitationImpact      float64 `bson:"学科规范化的引文影响力"`
	Top1Percent         float64 `bson:"被引次数排名前 1% 的论文百分比"`
	Top10Percent        float64 `bson:"被引次数排名前 10% 的论文百分比"`
	HighCitationPercent float64 `bson:"高被引论文百分比"`
	HighCitationPapers  int     `bson:"高被引论文"`
	PopularPercent      float64 `bson:"热门论文百分比"`
	HIndex              int     `bson:"h 指数"`
	PatentCitations     int     `bson:"专利引用"`
	PopularPapers       int     `bson:"热门论文"`
}

func AddESI(client *mongo.Client, title string, EsiFile *multipart.FileHeader, name string) error {
	// Connect to MongoDB
	collection := client.Database("test").Collection(name)

	esiCollection := client.Database("test").Collection("ESI_papers")

	fileUri, err := Asset.UploadFile(EsiFile)
	if err != nil {
		return err
	}

	data, err := parseESIHot(EsiFile)
	if err != nil {
		log.Println(err)
		return err
	}

	result, err := esiCollection.InsertMany(context.TODO(), data)
	if err != nil {
		log.Println(err)
		return err
	}

	document := bson.M{
		"title":          title,
		"filename":       EsiFile.Filename,
		"file":           fileUri,
		"parsed_content": result.InsertedIDs,
	}

	_, err = collection.InsertOne(context.Background(), document)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func GetAllEsi(client *mongo.Client, name string) ([]byte, error) {
	collection := client.Database("test").Collection(name)
	filter := bson.M{}
	projection := bson.M{"parsed_content": 0}
	cur, err := collection.Find(context.Background(), filter, options.Find().SetProjection(projection))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer cur.Close(context.Background())
	var documents []map[string]interface{}
	for cur.Next(context.Background()) {
		var document map[string]interface{}
		if err := cur.Decode(&document); err != nil {
			log.Println(err)
			return nil, err
		}
		documents = append(documents, document)
	}

	jsonBytes, err := json.Marshal(documents)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if string(jsonBytes) == "null" {
		jsonBytes = []byte("[]")
	}
	return jsonBytes, nil
}

func GetEsi(client *mongo.Client, name string, title string) ([]byte, error) {
	collection := client.Database("test").Collection(name)
	filter := bson.M{"title": title}
	projection := bson.M{"_id": 0, "parsed_content": 1}
	findOptions := options.FindOneOptions{Projection: projection}

	esiCollection := client.Database("test").Collection("ESI_papers")

	var esiResult map[string][]primitive.ObjectID
	err := collection.FindOne(context.Background(), filter, &findOptions).Decode(&esiResult)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	var results []interface{}
	for _, id := range esiResult["parsed_content"] {
		result := bson.M{}
		err = esiCollection.FindOne(context.Background(), bson.M{"_id": id},
			&options.FindOneOptions{Projection: bson.M{"_id": 0}}).Decode(result)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, result)
	}
	jsonBytes, err := json.Marshal(results)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if string(jsonBytes) == "null" {
		jsonBytes = []byte("[]")
	}
	return jsonBytes, nil
}

// TODO: Delete from ESI_paper
func DeleteEsi(client *mongo.Client, id string, name string) error {
	collection := client.Database("test").Collection(name)
	filter := bson.M{"_id": id}
	projection := bson.M{"_id": 0, "file": 1}
	findOptions := options.FindOneOptions{Projection: projection}
	var result string
	err := collection.FindOne(context.Background(), filter, &findOptions).Decode(&result)
	if err != nil {
		return errors.New("Invaild Index")
	}
	fileName := path.Base(result)
	err = Asset.DeleteFile(fileName)
	if err != nil {
		return errors.New("Unable to remove")
	}
	return nil
}
func parseESIHot(EsiFile *multipart.FileHeader) ([]interface{}, error) {
	// Open the Excel file
	tmpFile, err := EsiFile.Open()
	if err != nil {
		return nil, err
	}
	file, err := xlsx.OpenReaderAt(tmpFile, EsiFile.Size)
	if err != nil {
		return nil, err
	}

	// Iterate over the worksheets in the file
	for _, sheet := range file.Sheets {
		// Create a slice to store the data for this worksheet
		data := []interface{}{}

		// Iterate over the rows in the worksheet
	RowIter:
		for _, row := range sheet.Rows {
			// Create a map to store the data for this row
			rowData := make(map[string]interface{})

			// Iterate over the cells in the row
			for idx, cell := range row.Cells {
				// Get the cell value and add it to the row data
				val := cell.String()
				index := sheet.Cell(5, idx).String()
				switch index {
				case "Times Cited":
					number, err := strconv.Atoi(val)
					if err != nil {
						continue RowIter
					} else {
						rowData[index] = number
					}
				case "Publication Date":
					number, err := strconv.Atoi(val)
					if err != nil {
						continue RowIter
					} else {
						rowData[index] = number
					}
				case "Authors", "Countries", "Addresses", "Institutions":
					rowData[index] = strings.Split(val, ";")
				default:
					rowData[index] = val
				}

			}

			// Add the row data to the sheet data
			if len(rowData) > 10 {
				data = append(data, rowData)
			}

		}
		if len(data) < 1 {
			log.Println("Invalid ESI File")
			err := errors.New("invalid ESI File")
			return nil, err
		}
		return data, nil
	}
	return nil, nil
}

func AddIncites(client *mongo.Client, title string, IncitesFile *multipart.FileHeader) error {
	// Connect to MongoDB
	collection := client.Database("test").Collection("Incites")

	fileUri, err := Asset.UploadFile(IncitesFile)
	if err != nil {
		log.Println(err)
		return err
	}

	data, err := parseIncites(IncitesFile)
	if err != nil {
		log.Println(err)
		return err
	}

	document := bson.M{
		"title":          title,
		"filename":       IncitesFile.Filename,
		"file":           fileUri,
		"parsed_content": data,
	}

	_, err = collection.InsertOne(context.Background(), document)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func parseIncites(IncitesFile *multipart.FileHeader) ([]IncitesRow, error) {
	// Open the Excel file
	tmpFile, err := IncitesFile.Open()
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
	var data []IncitesRow
	// Loop over CSV records from Row 1
	for _, record := range records[1:] {
		if len(record) < 14 {
			continue // skip short records
		}
		incite := IncitesRow{
			Author:              record[0],
			WebOfSciencePapers:  castInt(record[1]),
			Rank:                castInt(record[2]),
			CitationFrequency:   castInt(record[3]),
			Researchers:         record[5],
			CitationImpact:      castFloat(record[6]),
			Top1Percent:         castFloat(record[7]),
			Top10Percent:        castFloat(record[8]),
			HighCitationPercent: castFloat(record[9]),
			HighCitationPapers:  castInt(record[10]),
			PopularPercent:      castFloat(record[11]),
			HIndex:              castInt(record[12]),
			PatentCitations:     castInt(record[13]),
			PopularPapers:       castInt(record[14]),
		}
		data = append(data, incite)
	}

	return data, nil
}

func castInt(in string) int {
	i, _ := strconv.Atoi(in)
	return i
}

func castFloat(in string) float64 {
	i, _ := strconv.ParseFloat(in, 64)
	return i
}

func GetIncites(client *mongo.Client) ([]byte, error) {
	collection := client.Database("test").Collection("Incites")
	return myJson.GetAllDocumentsAsJson(collection)
}
