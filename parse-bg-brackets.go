package main

import (
	"github.com/jharkins/parse-bg-brackets/toons"
	"gopkg.in/gin-gonic/gin.v1"
	"log"
	"net/http"
)

func main() {
	bNetURI := "https://worldofwarcraft.com/en-us/game/pvp/leaderboards/battlegrounds?page="
	numPages := 10

	tr, err := toons.NewToonRepository(bNetURI, numPages)
	defer tr.Close()

	if err != nil {
		panic(err)
	}

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.GET("/get-toons", func(c *gin.Context) {
		ts := tr.GetToons()
		c.JSON(200, gin.H{
			"toons": ts,
		})
	})

	r.GET("/", func(c *gin.Context) {
		ts := tr.GetToons()
		c.HTML(http.StatusOK, "index.html", gin.H{
			"toons": ts,
		})
	})

	r.GET("/class-stats", func(c *gin.Context) {
		if err != nil {
			log.Println(err)
			c.Abort()
			return
		}

		c.HTML(http.StatusOK, "class-stats.html", gin.H{
			"classCount": tr.GetClassCount(),
		})
	})

	r.GET("/spec-stats", func(c *gin.Context) {
		if err != nil {
			log.Println(err)
			c.Abort()
			return
		}

		c.HTML(http.StatusOK, "spec-stats.html", gin.H{
			"specCount": tr.GetSpecCount(),
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
