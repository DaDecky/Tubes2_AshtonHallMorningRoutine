package handler

import (
	"backend/utils"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// ✅ This is the only entrypoint for Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	// Create a new Gin engine per request (stateless)
	router := gin.New()

	// Middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Load recipes (can be optimized with sync.Once)
	utils.LoadRecipes("recipes.json")

	// Define routes
	router.GET("/search", func(c *gin.Context) {
		start := time.Now()
		response := utils.JSONResponse{Errors: []string{}}

		target := c.Query("target")
		algorithm_mode := c.Query("algo")
		search_mode := c.Query("shortest")
		max := c.Query("max")

		if algorithm_mode == "" || (search_mode == "" && max == "") || target == "" {
			response.Errors = append(response.Errors, "Missing Query Parameters")
			c.JSON(http.StatusBadRequest, response)
			return
		}

		useBFS := (algorithm_mode == "BFS")
		findShortest := (search_mode == "true")
		maxRecipes := 1

		if max != "" {
			if val, err := strconv.Atoi(max); err == nil {
				maxRecipes = val
			} else {
				response.Errors = append(response.Errors, "Max recipes must be a number")
				c.JSON(http.StatusBadRequest, response)
				return
			}
		}

		data, nodeCount, recipeFound, err := utils.Search(target, findShortest, useBFS, maxRecipes)
		if err != nil {
			response.Errors = append(response.Errors, err.Error())
			c.JSON(http.StatusBadRequest, response)
			return
		}

		response.Data = utils.ConvertToJSONFormat(data)
		response.NodeCount = nodeCount
		response.RecipeFound = recipeFound
		response.Time = time.Since(start).Milliseconds()

		c.JSON(http.StatusOK, response)
	})

	router.GET("/elements", func(c *gin.Context) {
		jsonData, err := os.ReadFile("elements.json") // Adjust if needed
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read JSON file"})
			return
		}

		var jsonObj interface{}
		if err := json.Unmarshal(jsonData, &jsonObj); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
			return
		}

		c.JSON(http.StatusOK, jsonObj)
	})

	// ⚠ WebSockets are NOT supported on Vercel — remove or disable this route
	// You can simulate with polling if needed

	// Handle the request
	router.ServeHTTP(w, r)
}
