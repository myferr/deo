package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/myferr/deo/storage"

	"github.com/gin-gonic/gin"
)

func ListDocuments(c *gin.Context) {
	dbName := c.Param("db_name")
	collectionName := c.Param("collection_name")

	// Parse filters
	var filters []storage.DocumentFilter
	for key, values := range c.Request.URL.Query() {
		if strings.HasPrefix(key, "filter[") && strings.HasSuffix(key, "]") {
			field := key[len("filter[") : len(key)-1]
			if len(values) > 0 {
				filters = append(filters, storage.DocumentFilter{Field: field, Value: values[0]})
			}
		}
	}

	// Parse sorting
	var sortParams *storage.DocumentSort
	sortBy := c.Query("sort_by")
	if sortBy != "" {
		sortOrder := c.DefaultQuery("order", "asc") // Default to ascending
		sortParams = &storage.DocumentSort{Field: sortBy, Order: sortOrder}
	}

	// Parse pagination
	var pagination *storage.DocumentPagination
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")

	limit := 0
	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err == nil && l > 0 {
			limit = l
		}
	}

	offset := 0
	if offsetStr != "" {
		o, err := strconv.Atoi(offsetStr)
		if err == nil && o >= 0 {
			offset = o
		}
	}

	if limit > 0 || offset > 0 {
		pagination = &storage.DocumentPagination{Limit: limit, Offset: offset}
	}

	documents, err := storage.ListDocuments(dbName, collectionName, filters, sortParams, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to list documents"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    documents,
	})
}
