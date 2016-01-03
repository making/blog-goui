package main

import (
	"github.com/making/blog-goui/service"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"strconv"
	"html/template"
)

const (
	DefaultPort = "9000"
	DefaultApiUrl = "https://blog-api.cfapps.io/api"
)

func main() {
	var port string
	if port = os.Getenv("PORT"); len(port) == 0 {
		log.Printf("Warning, PORT not set. Defaulting to %+v", DefaultPort)
		port = DefaultPort
	}
	var apiUrl string
	if apiUrl = os.Getenv("API_URL"); len(apiUrl) == 0 {
		log.Printf("Warning, API_URL not set. Defaulting to %+v", DefaultApiUrl)
		apiUrl = DefaultApiUrl
	}

	s := service.NewService(apiUrl)
	r := gin.Default()
	r.GET("/entries/:entryId", func(c *gin.Context) {
		entryId, err := strconv.ParseInt(c.Param("entryId"), 10, 32)
		if err != nil {
			log.Println("Error: ", err)
			c.String(400, "error!")
		} else {
			entry, _ := s.GetEntry(entryId)
			model := gin.H{
				"Title": entry.GetFrontMatter().GetTitle(),
				"Content": entry.GetContent(),
			}
			tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/entry.html"))
			tmpl.Execute(c.Writer, model)
		}
	})
	
	r.Run(":" + port)
}
