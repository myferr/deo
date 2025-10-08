package main

import (
	"deo/handlers"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

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

	fmt.Println("\033[1;34m" + `
     __
 ___/ /__ ___
/ _  / -_) _ \
\_,_/\__/\___/

` + "\033[0m")

	fmt.Println("\033[1;32mdeo (/'dioh/), is running on localhost:6741\033[0m")
	fmt.Println("\033[1;36mlocal → http://localhost:6741\033[0m")
	r.Run(":6741")
}
