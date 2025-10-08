package main

import (
	"deo/handlers"
	"deo/storage"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	r.LoadHTMLGlob("templates/*")

	r.GET("/studio", func(c *gin.Context) {
		dbs, err := storage.ListDatabases()
		if err != nil {
			dbs = []string{}
		}
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"Databases": dbs,
		})
	})

	studioApi := r.Group("/studio/api")
	{
		studioApi.GET("/dbs/:db_name/collections", handlers.ListCollections)
		studioApi.GET("/dbs/:db_name/collections/:collection_name/documents", handlers.ListDocuments)
	}

	api := r.Group("/api")
	{
		api.POST("/dbs", handlers.CreateDatabase)
	}

	dbs := r.Group("/api/dbs/:db_name")
	{
		dbs.POST("/collections", handlers.CreateCollection)
		collections := dbs.Group("/collections/:collection_name")
		{
			collections.POST("/documents", handlers.CreateDocument)
			collections.GET("/documents/:document_id", handlers.ReadDocument)
			collections.PUT("/documents/:document_id", handlers.UpdateDocument)
			collections.DELETE("/documents/:document_id", handlers.DeleteDocument)
		}
	}

	fmt.Println("\033[1;34m" + `
     __
 ___/ /__ ___
/ _  / -_) _ \
\_,_/\__/\___/

` + "\033[0m")

	fmt.Println("\033[1;32mdeo (/'dioh/), is running on localhost:6741\033[0m")
	fmt.Println("\033[1;36mlocal → http://localhost:6741/studio\033[0m")
	r.Run(":6741")
}
