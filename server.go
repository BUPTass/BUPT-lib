package main

import (
	esi "BUPT-lib/data_import"
	journal "BUPT-lib/data_import"
	news "BUPT-lib/data_import"
	data_retrieve "BUPT-lib/data_retrieve"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
)

func main() {
	mongoClient := ConnectToMongodb("mongodb://127.0.0.1:27017")
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "BUPT Library API backend")
	})
	e.GET("/journals", func(c echo.Context) error {
		collection := mongoClient.Database("test").Collection("Journals")
		JournalJson, err := data_retrieve.GetAllDocumentsAsJson(collection)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to retrieve")

		} else {
			return c.JSONBlob(http.StatusOK, JournalJson)
		}
	})
	e.POST("/journals", func(c echo.Context) error {
		journalXlsFile, err := c.FormFile("xlsx")
		if err != nil {
			return c.String(http.StatusBadRequest, "Error getting file from form: "+err.Error())
		}
		if journal.UpdateJournalList(mongoClient, journalXlsFile) != nil {
			return c.String(http.StatusInternalServerError, "Upload Failed")
		} else {
			return c.String(http.StatusOK, "Upload OK")
		}
	})
	e.POST("/journals/:id/image", func(c echo.Context) error {
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
	e.GET("/journals/:id/image", func(c echo.Context) error {
		id := c.Param("id")
		image, err := journal.GetImage(mongoClient, id)
		if err != nil {
			return c.String(http.StatusNotFound, err.Error())
		} else {
			return c.Blob(http.StatusOK, "image/jpeg", image)
		}
	})
	e.POST("/esi/paper", func(c echo.Context) error {
		title := c.FormValue("title")
		file, err := c.FormFile("xlsx")

		// 任意文件上传!
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		err = esi.AddESI(mongoClient, title, file)
		if err != nil {
			return c.String(http.StatusNotFound, err.Error())
		} else {
			return c.String(http.StatusOK, "New Esi Added")
		}
	})
	e.GET("/esi/paper", func(c echo.Context) error {
		esiJson, err := esi.GetEsi(mongoClient)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to retrieve")
		} else {
			return c.JSONBlob(http.StatusOK, esiJson)
		}
	})
	e.DELETE("/esi/paper/:id", func(c echo.Context) error {
		id := c.Param("id")
		err := esi.DeleteEsi(mongoClient, id)
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
		return nil
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
