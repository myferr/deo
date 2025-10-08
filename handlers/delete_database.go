package handlers

import (
	"net/http"

	"deo/storage"

	"github.com/gin-gonic/gin"
)

func DeleteDatabase(c *gin.Context) {
	dbName := c.Param("db_name")

	if err := storage.DeleteDatabase(dbName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to delete database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Database deleted successfully",
	})
}
