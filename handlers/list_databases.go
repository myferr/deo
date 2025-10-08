package handlers

import (
	"net/http"

	"deo/storage"

	"github.com/gin-gonic/gin"
)

func ListDatabases(c *gin.Context) {
	dbs, err := storage.ListDatabases()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to list databases"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    dbs,
	})
}
