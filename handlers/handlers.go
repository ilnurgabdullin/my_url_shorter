package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    //"fmt"
    //"crypto/sha256"
    //"encoding/hex"
    "url_shorter/storage"
)

func StatusCheck(c *gin.Context) {
    c.HTML(http.StatusOK, "index.html", nil)  
}


func OpenUrl(c *gin.Context) {
    hash := c.Params.ByName("hash")
    longUlr, _ := storage.GetRecordByHash(hash);
    c.Redirect(http.StatusMovedPermanently, longUlr.Long) 
}

//

func ShortUrl(c *gin.Context){
    var json struct {
        Message string `json:"url"`    
    }    

    if err := c.BindJSON(&json); err != nil {
        return    
    }
    newUrl, _ := storage.GenerateUniqueShortHash(json.Message,10)
    storage.InsertRecord(json.Message,newUrl)
    c.JSON(http.StatusOK, gin.H{
        "status": "http://"+storage.GetLocalIP().String()+"/"+newUrl,    
    })
}




