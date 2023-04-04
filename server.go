package main

import (
	esi "BUPT-lib/data_import"
	journal "BUPT-lib/data_import"
	news "BUPT-lib/data_import"
	data_retrieve "BUPT-lib/data_retrieve"
	utils "BUPT-lib/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
)

func main() {
	mongoClient := ConnectToMongodb("mongodb://mongodb_container:27017")
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
		esiJson, err := esi.GetEsi(mongoClient, "ESI_Cited")
		if err != nil {
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
		esiJson, err := esi.GetEsi(mongoClient, "ESI_Hot")
		if err != nil {
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
		num := utils.CastInt(c.FormValue("num"))

		newsJson, err := news.GetAnnouncement(mongoClient, num)

		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to retrieve")
		} else {
			return c.JSONBlob(http.StatusOK, newsJson)
		}
	})
	e.POST("/news/announcement", func(c echo.Context) error {
		file, err := c.FormFile("csv")

		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		err = esi.AddAnnouncement(mongoClient, file)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		} else {
			return c.String(http.StatusOK, "News Added")
		}
	})
	e.GET("/news/news", func(c echo.Context) error {
		num := utils.CastInt(c.FormValue("num"))

		newsJson, err := news.GetNews(mongoClient, num)

		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to retrieve")
		} else {
			return c.JSONBlob(http.StatusOK, newsJson)
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
		err = esi.AddOngoingConferences(mongoClient, file)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		} else {
			return c.String(http.StatusOK, "News Added")
		}
	})
	e.GET("/news/news_lib_res", func(c echo.Context) error {
		num := utils.CastInt(c.FormValue("num"))

		newsJson, err := news.GetLibNews(mongoClient, num, 1)

		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to retrieve")
		} else {
			return c.JSONBlob(http.StatusOK, newsJson)
		}
	})
	e.POST("/news/news_lib_res", func(c echo.Context) error {
		file, err := c.FormFile("csv")

		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		err = esi.AddLibNews(mongoClient, file, 1)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		} else {
			return c.String(http.StatusOK, "News Added")
		}
	})
	e.GET("/news/news_lib_ann", func(c echo.Context) error {
		num := utils.CastInt(c.FormValue("num"))

		newsJson, err := news.GetLibNews(mongoClient, num, 2)

		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to retrieve")
		} else {
			return c.JSONBlob(http.StatusOK, newsJson)
		}
	})
	e.POST("/news/news_lib_ann", func(c echo.Context) error {
		file, err := c.FormFile("csv")

		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		err = esi.AddLibNews(mongoClient, file, 2)
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
