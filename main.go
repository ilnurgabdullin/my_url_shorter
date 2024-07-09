package main 

import (
	"github.com/gin-gonic/gin"
    "url_shorter/handlers"
	)

func main() {
	r:= gin.Default()

	r.GET("/status", handlers.StatusCheck)
	r.POST("/short", handlers.ShortUrl)

    r.Run(":8080")

}
