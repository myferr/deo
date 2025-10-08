package handlers

import (
	"net/http"

	"deo/storage"

	"github.com/gin-gonic/gin"
)

func DeleteCollection(c *gin.Context) {
	dbName := c.Param("db_name")
	collectionName := c.Param("collection_name")

	if err := storage.DeleteCollection(dbName, collectionName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to delete collection"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Collection deleted successfully",
	})
}
