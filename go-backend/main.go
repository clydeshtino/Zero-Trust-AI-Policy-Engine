package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type QueryRequest struct {
	Query string `json:"query" binding:"required"`
}

type QueryResponse struct {
	Response string `json:"response"`
}

type RAGResponse struct {
	Response string `json:"response"`
}

func main() {
	r := gin.Default()
	r.Use(CORSMiddleware())

	r.POST("/api/query", func(c *gin.Context) {
		var req QueryRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// Call Python RAG service
		resp, err := forwardToRAG(req.Query)
		if err != nil {
			c.JSON(500, gin.H{"error": "RAG service down"})
			return
		}

		c.JSON(200, QueryResponse{Response: resp.Response})
	})

	log.Println("Go API running on :8080")
	r.Run(":8080")
}

func forwardToRAG(query string) (*RAGResponse, error) {
	payload := map[string]string{"query": query}
	jsonData, _ := json.Marshal(payload)

	resp, err := http.Post("http://python-rag:8000/query", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var ragResp RAGResponse
	json.Unmarshal(body, &ragResp)
	return &ragResp, nil
}

// CORS
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
