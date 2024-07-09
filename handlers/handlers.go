package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func StatusCheck(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
    "status":"ok",
    })   
}

func ShortUrl(c *gin.Context){
    var json struct {
        Message string `json:"url"`    
    }    

    if err := c.BindJSON(&json); err != nil {
        return    
    }

    c.JSON(http.StatusOK, gin.H{
        "ans": "new_url",    
    })
}
