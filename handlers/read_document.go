package handlers

import (
	"net/http"

	"deo/storage"

	"github.com/gin-gonic/gin"
)

func ReadDocument(c *gin.Context) {
	dbName := c.Param("db_name")
	collectionName := c.Param("collection_name")
	docID := c.Param("document_id")

	data, err := storage.LoadDocument(dbName, collectionName, docID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Document not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}
