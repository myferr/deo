package handlers

import (
	"net/http"

	"deo/storage"

	"github.com/gin-gonic/gin"
)

type CreateDatabaseRequest struct {
	DBName string `json:"db_name" binding:"required"`
}

func CreateDatabase(c *gin.Context) {
	var req CreateDatabaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request body"})
		return
	}

	if err := storage.CreateDatabase(req.DBName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to create database"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Database created successfully",
	})
}
