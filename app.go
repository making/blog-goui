package main

import (
	"github.com/making/blog-goui/service"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

const (
	DEFAULT_PORT = "9000"
)

func main() {
	var port string
	if port = os.Getenv("PORT"); len(port) == 0 {
		log.Printf("Warning, PORT not set. Defaulting to %+v", DEFAULT_PORT)
		port = DEFAULT_PORT
	}

	s := service.NewService("https://blog-api.cfapps.io/api")
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		entry, _ := s.GetEntry(300)
		c.String(200, entry.GetFrontMatter().GetTitle())
	})
	r.Run(":" + port)
}
