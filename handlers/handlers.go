package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    //"fmt"
    "crypto/sha256"
    "encoding/hex"
    "url_shorter/storage"
)

func StatusCheck(c *gin.Context) {
    c.HTML(http.StatusOK, "index.html", nil)  
}


func OpenUrl(c *gin.Context) {
    hash := c.Params.ByName("hash")
    longUlr, _ := storage.GetRecordByHash(hash);
    c.Redirect(http.StatusMovedPermanently, longUlr.Long)
    //c.JSON(http.StatusOK, gin.H{"answer":longUlr,})   
}

//

func ShortUrl(c *gin.Context){
    var json struct {
        Message string `json:"url"`    
    }    

    if err := c.BindJSON(&json); err != nil {
        return    
    }
    newUrl := GetShortHash(json.Message,10)
    storage.InsertRecord(json.Message,newUrl)
    c.JSON(http.StatusOK, gin.H{
        "status": "http://127.0.0.1:8080/o/"+newUrl,    
    })
}



func GetShortHash(input string, length int) string {
    hash := sha256.New()
    hash.Write([]byte(input))
    hashBytes := hash.Sum(nil)

    hashString := hex.EncodeToString(hashBytes)

    if length > len(hashString) {
        length = len(hashString)
    }
    return hashString[:length]
}
