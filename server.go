package main

import (
	esi "BUPT-lib/data_import"
	journal "BUPT-lib/data_import"
	news "BUPT-lib/data_import"
	"BUPT-lib/data_retrieve"
	"BUPT-lib/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"strconv"
)

func main() {
	mongoClient := ConnectToMongodb("mongodb://mongodb_container:27017")
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "BUPT Library API backend")
	})
	e.Static("/static", "static")
	e.GET("/journals/top30", func(c echo.Context) error {
		collection := mongoClient.Database("test").Collection("Journals")
		JournalJson, err := data_retrieve.GetAllDocumentsAsJson(collection)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to retrieve")

		} else {
			return c.JSONBlob(http.StatusOK, JournalJson)
		}
	})
	e.PUT("/journals/top30", func(c echo.Context) error {
		journalXlsFile, err := c.FormFile("xlsx")
		if err != nil {
			return c.String(http.StatusBadRequest, "Error getting file from form: "+err.Error())
		}
		if journal.PutT30JournalList(mongoClient, journalXlsFile) != nil {
			return c.String(http.StatusInternalServerError, "Upload Failed")
		} else {
			return c.String(http.StatusOK, "Upload OK")
		}
	})
	e.POST("/journals/top30/:id/image", func(c echo.Context) error {
		id := c.Param("id")
		image, err := c.FormFile("img")
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		err = journal.SetImage(mongoClient, id, image)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		} else {
			return c.String(http.StatusOK, "set img")
		}
	})
	e.GET("/journals/top30/:id/image", func(c echo.Context) error {
		id := c.Param("id")
		image, err := journal.GetImage(mongoClient, id)
		if err != nil {
			return c.String(http.StatusNotFound, err.Error())
		} else {
			return c.String(http.StatusOK, image)
		}
	})
	e.GET("/journals/ccf", func(c echo.Context) error {
		ccfJson, err := journal.GetCcfList(mongoClient)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to retrieve")

		} else {
			return c.JSONBlob(http.StatusOK, ccfJson)
		}
	})
	e.PUT("/journals/ccf", func(c echo.Context) error {
		file, err := c.FormFile("csv")
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		err = journal.PutCcfList(mongoClient, file)

		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to retrieve")

		} else {
			return c.String(http.StatusOK, "Set CCF")
		}
	})
	e.POST("/esi/highlycited", func(c echo.Context) error {
		title := c.FormValue("title")
		file, err := c.FormFile("xlsx")

		// 任意文件上传!
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		err = esi.AddESI(mongoClient, title, file, "ESI_Cited")
		if err != nil {
			return c.String(http.StatusNotFound, err.Error())
		} else {
			return c.String(http.StatusOK, "New Esi Added")
		}
	})
	e.GET("/esi/highlycited", func(c echo.Context) error {
		esiJson, err := esi.GetAllEsi(mongoClient, "ESI_Cited")
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to retrieve")
		} else {
			return c.JSONBlob(http.StatusOK, esiJson)
		}
	})
	e.GET("/esi/highlycited/:title", func(c echo.Context) error {
		title := c.Param("title")
		esiJson, err := esi.GetEsi(mongoClient, "ESI_Cited", title)

		if err == mongo.ErrNoDocuments {
			return c.String(http.StatusNotFound, "Not Found")
		} else if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to retrieve")
		} else {
			return c.JSONBlob(http.StatusOK, esiJson)
		}
	})
	e.DELETE("/esi/highlycited/:id", func(c echo.Context) error {
		id := c.Param("id")
		err := esi.DeleteEsi(mongoClient, id, "ESI_Cited")
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to delete")
		} else {
			return c.String(http.StatusOK, "Success")
		}
	})
	e.POST("/esi/hot", func(c echo.Context) error {
		title := c.FormValue("title")
		file, err := c.FormFile("xlsx")

		// 任意文件上传!
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		err = esi.AddESI(mongoClient, title, file, "ESI_Hot")
		if err != nil {
			return c.String(http.StatusNotFound, err.Error())
		} else {
			return c.String(http.StatusOK, "New Esi Added")
		}
	})
	e.GET("/esi/hot", func(c echo.Context) error {
		esiJson, err := esi.GetAllEsi(mongoClient, "ESI_Hot")

		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to retrieve")
		} else {
			return c.JSONBlob(http.StatusOK, esiJson)
		}
	})
	e.GET("/esi/hot/:title", func(c echo.Context) error {
		title := c.Param("title")
		esiJson, err := esi.GetEsi(mongoClient, "ESI_Hot", title)

		if err == mongo.ErrNoDocuments {
			return c.String(http.StatusNotFound, "Not Found")
		} else if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to retrieve")
		} else {
			return c.JSONBlob(http.StatusOK, esiJson)
		}
	})
	e.DELETE("/esi/hot/:id", func(c echo.Context) error {
		id := c.Param("id")
		err := esi.DeleteEsi(mongoClient, id, "ESI_Hot")
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to delete")
		} else {
			return c.String(http.StatusOK, "Success")
		}
	})
	e.POST("/incites", func(c echo.Context) error {
		title := c.FormValue("title")
		file, err := c.FormFile("csv")
		// 任意文件上传!
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		err = esi.AddIncites(mongoClient, title, file)
		if err != nil {
			return c.String(http.StatusNotFound, err.Error())
		} else {
			return c.String(http.StatusOK, "New Incites Added")
		}
	})
	e.GET("/incites", func(c echo.Context) error {
		esiJson, err := esi.GetIncites(mongoClient)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to retrieve")
		} else {
			return c.JSONBlob(http.StatusOK, esiJson)
		}
	})
	e.POST("/search/:parm", func(c echo.Context) error {
		keyword := c.Param("parm")

		// where and how to search
		data_retrieve.SearchAll(keyword)

		return nil
	})
	e.POST("/news", func(c echo.Context) error {
		newsJson := c.FormValue("parm")

		err := news.UpdateNews(mongoClient, newsJson)

		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to update")
		} else {
			return c.String(http.StatusOK, "Successfully Update News")
		}
	})
	e.GET("/news/announcement", func(c echo.Context) error {
		num, err := strconv.Atoi(c.FormValue("num"))
		if err != nil || num < 0 {
			num = 10
		}
		start, err := strconv.Atoi(c.FormValue("start"))
		if err != nil || start < 0 {
			start = 0
		}
		newsJson, err := news.GetAnnouncement(mongoClient, uint(num), uint(start))

		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to retrieve")
		} else {
			return c.JSONBlob(http.StatusOK, newsJson)
		}
	})
	e.GET("/news/announcement/total", func(c echo.Context) error {

		totalNum, err := news.CountNews(mongoClient, 3)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to retrieve")
		} else {
			return c.String(http.StatusOK, strconv.FormatInt(totalNum, 10))
		}
	})
	e.POST("/news/announcement", func(c echo.Context) error {
		file, err := c.FormFile("csv")

		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		err = news.AddAnnouncement(mongoClient, file)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		} else {
			return c.String(http.StatusOK, "News Added")
		}
	})
	e.GET("/news/news", func(c echo.Context) error {
		num, err := strconv.Atoi(c.FormValue("num"))
		if err != nil || num < 0 {
			num = 10
		}
		start, err := strconv.Atoi(c.FormValue("start"))
		if err != nil || start < 0 {
			start = 0
		}
		newsJson, err := news.GetNews(mongoClient, uint(num), uint(start))

		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to retrieve")
		} else {
			return c.JSONBlob(http.StatusOK, newsJson)
		}
	})
	e.GET("/news/news/total", func(c echo.Context) error {

		totalNum, err := news.CountNews(mongoClient, 0)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to retrieve")
		} else {
			return c.String(http.StatusOK, strconv.FormatInt(totalNum, 10))
		}
	})
	e.POST("/news/news", func(c echo.Context) error {
		file, err := c.FormFile("csv")

		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		err = esi.AddNews(mongoClient, file)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		} else {
			return c.String(http.StatusOK, "News Added")
		}
	})
	e.GET("/news/conf", func(c echo.Context) error {
		newsJson, err := news.GetOngoingConferences(mongoClient)

		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to retrieve")
		} else {
			return c.JSONBlob(http.StatusOK, newsJson)
		}
	})
	e.POST("/news/conf", func(c echo.Context) error {
		file, err := c.FormFile("csv")

		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		err = news.AddOngoingConferences(mongoClient, file)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		} else {
			return c.String(http.StatusOK, "News Added")
		}
	})
	e.GET("/news/news_lib_res", func(c echo.Context) error {
		num, err := strconv.Atoi(c.FormValue("num"))
		if err != nil || num < 0 {
			num = 10
		}
		start, err := strconv.Atoi(c.FormValue("start"))
		if err != nil || start < 0 {
			start = 0
		}
		newsJson, err := news.GetLibNews(mongoClient, uint(num), uint(start), 1)

		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to retrieve")
		} else {
			return c.JSONBlob(http.StatusOK, newsJson)
		}
	})
	e.GET("/news/news_lib_res/total", func(c echo.Context) error {

		totalNum, err := news.CountNews(mongoClient, 1)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to retrieve")
		} else {
			return c.String(http.StatusOK, strconv.FormatInt(totalNum, 10))
		}
	})
	e.POST("/news/news_lib_res", func(c echo.Context) error {
		file, err := c.FormFile("csv")

		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		err = news.AddLibNews(mongoClient, file, 1)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		} else {
			return c.String(http.StatusOK, "News Added")
		}
	})
	e.GET("/news/news_lib_ann", func(c echo.Context) error {
		num := utils.CastInt(c.FormValue("num"))
		num, err := strconv.Atoi(c.FormValue("num"))
		if err != nil || num < 0 {
			num = 10
		}
		start, err := strconv.Atoi(c.FormValue("start"))
		if err != nil || start < 0 {
			start = 0
		}
		newsJson, err := news.GetLibNews(mongoClient, uint(num), uint(start), 2)

		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to retrieve")
		} else {
			return c.JSONBlob(http.StatusOK, newsJson)
		}
	})
	e.GET("/news/news_lib_ann/total", func(c echo.Context) error {

		totalNum, err := news.CountNews(mongoClient, 2)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to retrieve")
		} else {
			return c.String(http.StatusOK, strconv.FormatInt(totalNum, 10))
		}
	})
	e.POST("/news/news_lib_ann", func(c echo.Context) error {
		file, err := c.FormFile("csv")

		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		err = news.AddLibNews(mongoClient, file, 2)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		} else {
			return c.String(http.StatusOK, "News Added")
		}
	})
	e.Logger.Fatal(e.Start(":1323"))
}

func ConnectToMongodb(uri string) *mongo.Client {
	var c options.ClientOptions
	c.ApplyURI(uri)
	client, err := mongo.NewClient(&c)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(nil)
	if err != nil {
		log.Fatal(err)
	}
	return client
}
