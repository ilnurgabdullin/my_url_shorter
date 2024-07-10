package main 

import (
	"github.com/gin-gonic/gin"
    "url_shorter/handlers"
    "url_shorter/storage"
    "log"
    // "database/sql"
    // "fmt"
    _ "github.com/lib/pq"
	)




func main() {
	r:= gin.Default()
    r.LoadHTMLGlob("templates/*")
	r.GET("/", handlers.StatusCheck)
    r.GET("/o/:hash", handlers.OpenUrl)
	r.POST("/short", handlers.ShortUrl)
    
    storage.InitDB()
    r.Run(":8080")

}
