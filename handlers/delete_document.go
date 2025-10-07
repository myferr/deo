package handlers

import (
	"net/http"
	"os"

	"deo/storage"

	"github.com/gin-gonic/gin"
)

func DeleteDocument(c *gin.Context) {
	dbName := c.Param("db_name")
	collectionName := c.Param("collection_name")
	docID := c.Param("document_id")

	docPath, err := storage.GetDocPath(dbName, collectionName, docID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Error constructing document path"})
		return
	}

	if _, err := os.Stat(docPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Document not found"})
		return
	}

	if err := os.Remove(docPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to delete document"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Document deleted successfully",
	})
}
