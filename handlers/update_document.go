package handlers

import (
	"net/http"

	"deo/storage"

	"github.com/gin-gonic/gin"
)

func UpdateDocument(c *gin.Context) {
	dbName := c.Param("db_name")
	collectionName := c.Param("collection_name")
	docID := c.Param("document_id")

	var jsonData map[string]interface{}
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid JSON"})
		return
	}

	if err := storage.SaveDocument(dbName, collectionName, docID, jsonData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to update document"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Document updated successfully",
		"data":    jsonData,
	})
}
