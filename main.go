package main

import (
	"fmt"
	"net/http"

	"github.com/myferr/deo/handlers"
	"github.com/myferr/deo/storage"

	"github.com/gin-gonic/gin"

	"embed"
	"html/template"
)

//go:embed templates/*
var templates embed.FS

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	templ := template.Must(template.New("").ParseFS(templates, "templates/*"))
	r.SetHTMLTemplate(templ)

	r.GET("/studio", func(c *gin.Context) {
		dbs, err := storage.ListDatabases()
		if err != nil {
			dbs = []string{}
		}
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"Databases": dbs,
		})
	})

	api := r.Group("/api")
	{
		api.POST("/dbs", handlers.CreateDatabase)
		api.GET("/dbs", handlers.ListDatabases)
	}

	dbsGroup := r.Group("/api/dbs/:db_name")
	{
		dbsGroup.GET("/collections", handlers.ListCollections)
		dbsGroup.POST("/collections", handlers.CreateCollection)
		dbsGroup.DELETE("/", handlers.DeleteDatabase)

		collectionsGroup := dbsGroup.Group("/collections/:collection_name")
		{
			collectionsGroup.GET("/documents", handlers.ListDocuments)
			collectionsGroup.POST("/documents", handlers.CreateDocument)
			collectionsGroup.GET("/documents/:document_id", handlers.ReadDocument)
			collectionsGroup.PUT("/documents/:document_id", handlers.UpdateDocument)
			collectionsGroup.DELETE("/documents/:document_id", handlers.DeleteDocument)
			collectionsGroup.DELETE("/", handlers.DeleteCollection)
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
