package main 

import (
	"github.com/gin-gonic/gin"
    "url_shorter/handlers"
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
	)

const (
    host = "192.168.1.8"
    port = 5432
    user = "postgres"
    password = "123"
    dbname = "baza"
)


func main() {
	r:= gin.Default()

	r.GET("/status", handlers.StatusCheck)
	r.POST("/short", handlers.ShortUrl)
    
    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
        fmt.Println("=================================")
        fmt.Println(err)
        return
    }
    
    defer db.Close()
    err = db.Ping()
    if err != nil {
            fmt.Println("ERROR 2")
    }
    fmt.Println("Succes")
    r.Run(":8080")

}
