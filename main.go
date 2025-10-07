package main

import (
	"deo/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	api := r.Group("/api")
	{
		dbs := api.Group("/dbs/:db_name")
		{
			collections := dbs.Group("/collections/:collection_name")
			{
				collections.POST("/documents", handlers.CreateDocument)
				collections.GET("/documents/:document_id", handlers.ReadDocument)
				collections.PUT("/documents/:document_id", handlers.UpdateDocument)
				collections.DELETE("/documents/:document_id", handlers.DeleteDocument)
			}
		}
	}

	r.Run()
}
