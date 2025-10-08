package handlers

import (
	"net/http"

	"deo/storage"

	"github.com/gin-gonic/gin"
)

type CreateCollectionRequest struct {
	CollectionName string `json:"collection_name" binding:"required"`
}

func CreateCollection(c *gin.Context) {
	dbName := c.Param("db_name")

	var req CreateCollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request body"})
		return
	}

	if err := storage.CreateCollection(dbName, req.CollectionName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to create collection"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Collection created successfully",
	})
}
