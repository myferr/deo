package handlers

import (
	"net/http"

	"deo/storage"

	"github.com/gin-gonic/gin"
)

func ListDocuments(c *gin.Context) {
	dbName := c.Param("db_name")
	collectionName := c.Param("collection_name")
	documents, err := storage.ListDocuments(dbName, collectionName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to list documents"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    documents,
	})
}
