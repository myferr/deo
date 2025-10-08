package handlers

import (
	"net/http"

	"deo/storage"

	"github.com/gin-gonic/gin"
)

func ListCollections(c *gin.Context) {
	dbName := c.Param("db_name")
	collections, err := storage.ListCollections(dbName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to list collections"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    collections,
	})
}
