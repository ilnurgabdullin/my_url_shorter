package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "fmt"
    "crypto/sha256"
    "encoding/hex"
)

func StatusCheck(c *gin.Context) {
    c.HTML(http.StatusOK, "index.html", nil)  
}

func OpenUrl(c *gin.Context) {
    hash := c.Params.ByName("hash")
    c.JSON(http.StatusOK, gin.H{
    "answer":hash,
    })   
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
    fmt.Println("I get data", json.Message)
    fmt.Println("I make data", newUrl)
    c.JSON(http.StatusOK, gin.H{
        "status": "http://127.0.0.1:8080/o/"+newUrl,    
    })
}



func GetShortHash(input string, length int) string {
    // Применение хеш-функции SHA256 к строке
    hash := sha256.New()
    hash.Write([]byte(input))
    hashBytes := hash.Sum(nil)

    // Преобразование байтового среза в строку в формате hex
    hashString := hex.EncodeToString(hashBytes)

    // Оставляем только первые 'length' символов хеша
    if length > len(hashString) {
        length = len(hashString)
    }
    return hashString[:length]
}
