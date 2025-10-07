package handlers

import (
	"net/http"

	"deo/storage"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateDocument(c *gin.Context) {
	dbName := c.Param("db_name")
	collectionName := c.Param("collection_name")

	var jsonData map[string]any
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid JSON"})
		return
	}

	docID := uuid.New().String()
	jsonData["_id"] = docID

	if err := storage.SaveDocument(dbName, collectionName, docID, jsonData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to create document"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Document created successfully",
		"data":    jsonData,
	})
}
