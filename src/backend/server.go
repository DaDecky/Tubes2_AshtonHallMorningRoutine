package main

import (
	"fmt"
	"net/http"
	"time"
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

func main() {
	router := gin.Default()
  
	// search biasa
	router.POST("/search", func(c *gin.Context) {
		start := time.Now()
		algorithm_mode := c.Query("algo") // bfs dfs
		search_mode := c.Query("mode") // multi or shortest
		max_recipe := c.Query("max") // max recipe tree if using multi mode

		if algorithm_mode == "" || search_mode == "" || max_recipe == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing query parameters"})
			return
		}

		// do algo here i guess


		// send search result
		c.JSON(http.StatusOK, gin.H{
			"status": "Success",
			"data": "data here",
			"duration": time.Since(start).Seconds(),
		})
	})


	// live search pake websocket
	router.GET("/liveSearch", func(c *gin.Context) {
		start := time.Now()
		algorithm_mode := c.Query("algo") // bfs or dfs
		search_mode := c.Query("mode") // multi or shortest
		max_recipe := c.Query("max") // max recipe tree if using multi mode

		if algorithm_mode == "" || search_mode == "" || max_recipe == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Missing query parameters"})
			return
		}
		
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil);
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
				"type": "progress",
				"progress_counter":  i,
				"data": fmt.Sprintf("data %d", i),
				"timestamp": time.Now(),
			})

			if err != nil {
				fmt.Println("Write error:", err)
				return
			}

			time.Sleep(time.Second / 2)
		}


		// algo done so we send message telling its done to sever connection
		conn.WriteJSON(map[string]interface{}{
			"type": "complete",
			"duration": time.Since(start).Seconds(),
		})
	})
  
	router.Run("localhost:8081")
}