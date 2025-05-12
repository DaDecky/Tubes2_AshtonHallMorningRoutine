package handler

import (
	"backend/utils"
	"encoding/json"
	"fmt"
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
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var router = gin.Default()

func main() {
	// initialize recipes data
	// utils.InitializeData() <-------- scrapping. just uncomment for production
	utils.LoadRecipes("recipes.json")

	// cors
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080"}, // Frontend origin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// search biasa
	router.GET("/search", func(c *gin.Context) {
		start := time.Now()
		response := utils.JSONResponse{
			Errors: []string{},
		}

		// get query params
		target := c.Query("target")        // target recipe
		algorithm_mode := c.Query("algo")  // bfs dfs
		search_mode := c.Query("shortest") // multi or shortest
		max := c.Query("max")              // max recipe tree if using multi mode

		// validate query params
		if algorithm_mode == "" || (search_mode == "" && max == "") || target == "" {
			response.Errors = append(response.Errors, "Missing Query Parameters")
			c.JSON(http.StatusBadRequest, response)
			return
		}

		if search_mode != "true" && max == "" {
			response.Errors = append(response.Errors, "Missing Query Parameters")
			c.JSON(http.StatusBadRequest, response)
			return
		}

		// parse query params
		useBFS := false
		findShortest := false
		maxRecipes := 1
		if algorithm_mode == "BFS" {
			useBFS = true
		}
		if search_mode == "true" {
			findShortest = true
		}
		if max != "" {
			if val, err := strconv.Atoi(max); err == nil {
				maxRecipes = val
			} else {
				response.Errors = append(response.Errors, "Max recipes paramaeter must be a number")
				c.JSON(http.StatusBadRequest, response)
				return
			}
		}

		// search recipe
		data, nodeCount, recipeFound, err := utils.Search(target, findShortest, useBFS, maxRecipes)
		if err != nil {
			response.Errors = append(response.Errors, err.Error())
			c.JSON(http.StatusBadRequest, response)
			return
		}

		// send search result
		response.Data = utils.ConvertToJSONFormat(data)
		response.NodeCount = nodeCount
		response.RecipeFound = recipeFound
		response.Time = time.Since(start).Milliseconds()

		c.JSON(http.StatusOK, response)
	})

	// live search pake websocket
	router.GET("/liveSearch", func(c *gin.Context) {
		start := time.Now()
		algorithm_mode := c.Query("algo") // bfs or dfs
		search_mode := c.Query("mode")    // multi or shortest
		max_recipe := c.Query("max")      // max recipe tree if using multi mode

		if algorithm_mode == "" || search_mode == "" || max_recipe == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Missing query parameters"})
			return
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println("error upgrading http connection to a websocket.")
			return
		}

		// cara nge send data tinggal ganti yang WriteJSON ini jadi data lu,
		// jadi for loop yg dibawah ini di inkorporasikan ke fungsi bfs/dfs nya buat
		// ngesend node nya yang sekarang atau terserah gmn sih

		defer conn.Close()
		for i := 0; i < 10; i++ {
			// Send progress data
			err := conn.WriteJSON(map[string]interface{}{
				"type":             "progress",
				"progress_counter": i,
				"data":             fmt.Sprintf("data %d", i),
				"timestamp":        time.Now(),
			})

			if err != nil {
				fmt.Println("Write error:", err)
				return
			}

			time.Sleep(time.Second / 2)
		}

		// algo done so we send message telling its done to sever connection
		conn.WriteJSON(map[string]interface{}{
			"type":     "complete",
			"duration": time.Since(start).Seconds(),
		})
	})

	router.GET("/elements", func(c *gin.Context) {
		// Read file
		jsonData, err := os.ReadFile("elements.json")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to read JSON file",
			})
			return
		}

		// parse json
		var jsonObj interface{}
		if err := json.Unmarshal(jsonData, &jsonObj); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid JSON format",
			})
			return
		}

		// send json
		c.JSON(http.StatusOK, jsonObj)
	})

	router.Run(":8081")
}

func Handler(w http.ResponseWriter, r *http.Request) {
	router.ServeHTTP(w, r)
}
