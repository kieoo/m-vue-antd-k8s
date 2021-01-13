package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
)

type Pb struct {
	Code string `json:"code"`
}

func main() {
	r := gin.Default()
	r.Use(CorsMiddleware())
	r.GET("version", func(c *gin.Context){
		c.String(http.StatusOK, "OK")
	})

	r.POST("return", func(c *gin.Context) {
		body := Pb{}
		err := c.BindJSON(&body)
		if err != nil{
			fmt.Println(err.Error())
			c.String(http.StatusBadRequest, "bad")
		} else {
			c.String(http.StatusOK, string(body.Code))
		}
	})
	r.Run(":7001")

}

func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		var filterHost = [...]string{"http://192.168.1.5:8080"}
		// filterHost 做过滤器，防止不合法的域名访问
		var isAccess = false
		for _, v := range(filterHost) {
			match, _ := regexp.MatchString(v, origin)
			if match {
				isAccess = true
			}
		}
		if isAccess {
			// 核心处理方式
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
			c.Header("Access-Control-Allow-Methods", "GET, OPTIONS, POST, PUT, DELETE")
			c.Set("content-type", "application/json")
		}
		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}

		c.Next()
	}
}