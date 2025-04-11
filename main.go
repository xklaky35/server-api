package main

/* -------------------------------------------------------------------
This component is a interface to other internal services so 
they dont have to be exposed. 
If you want to use their enpoints your request needs to be validated
first before the request is redirected internaly
----------------------------------------------------------------------*/

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/xklaky35/welcomePageAPI/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil{
		fmt.Println(err)
		return
	}
	
	r := gin.Default()
	r.Use(middleware.RateLimiter())
	r.Use(middleware.HeaderSetup())

	// Setup auth
	r.Use(middleware.LoadValidUsers())
	r.Use(middleware.AuthMiddleware())

	// Redirection of welcomePageAPI URLs
	r.GET("/wP/:endpoint", func(ctx *gin.Context) {
		redirectURL := fmt.Sprintf("http://welcomepageapi:3001/%s", strings.Replace(ctx.Request.URL.String(), "/wP/", "", -1))
		data, code := RedirectGetRequest(ctx, redirectURL)
		if err != nil {
			ctx.AbortWithStatus(code)
		} 
		ctx.Data(http.StatusOK, "application/json", data)
	} )
	r.POST("/wP/:endpoint", func(ctx *gin.Context) {
		redirectURL := fmt.Sprintf("http://welcomepageapi:3001/%s", strings.Replace(ctx.Request.URL.String(), "/wP/", "", -1))
		code := RedirectPostRequest(ctx, redirectURL)
		if code != http.StatusOK {
			ctx.AbortWithStatus(code)
		}
	} )

	// Redirection of speed URLs
	r.GET("/speed/:endpoint", func(ctx *gin.Context) {
		redirectURL := fmt.Sprintf("http://welcomepageapi:3002/%s", strings.Replace(ctx.Request.URL.String(), "/speed/", "", -1))
		data, code := RedirectGetRequest(ctx, redirectURL)
		if err != nil {
			ctx.AbortWithStatus(code)
		} 
		ctx.Data(http.StatusOK, "application/json", data)
	} )
	r.POST("/speed/:endpoint", func(ctx *gin.Context) {
		redirectURL := fmt.Sprintf("http://welcomepageapi:3002/%s", strings.Replace(ctx.Request.URL.String(), "/speed/", "", -1))
		code := RedirectPostRequest(ctx, redirectURL)
		if code != http.StatusOK {
			ctx.AbortWithStatus(code)
		}
	} )

	// Redirection of p2g URLs
	r.GET("/p2g/:endpoint", func(ctx *gin.Context) {
		redirectURL := fmt.Sprintf("http://welcomepageapi:3003/%s", strings.Replace(ctx.Request.URL.String(), "/p2g/", "", -1))
		data, code := RedirectGetRequest(ctx, redirectURL)
		if err != nil {
			ctx.AbortWithStatus(code)
		} 
		ctx.Data(http.StatusOK, "application/json", data)
	})
	r.POST("/p2g/:endpoint", func(ctx *gin.Context) {
		redirectURL := fmt.Sprintf("http://welcomepageapi:3003/%s", strings.Replace(ctx.Request.URL.String(), "/p2g/", "", -1))
		code := RedirectPostRequest(ctx, redirectURL)
		if code != http.StatusOK {
			ctx.AbortWithStatus(code)
		}
	})
	
	r.Run(os.Getenv("PORT"))
}


func RedirectGetRequest(ctx *gin.Context, url string) ([]byte,int) {
	var client http.Client
	var data []byte

	// make new request
	req, err := http.NewRequest(http.MethodGet,url , ctx.Request.Body)
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil,resp.StatusCode
	}

	// read response
	sc := bufio.NewScanner(resp.Body)	
	defer resp.Body.Close()
	for sc.Scan() {
		data = []byte(sc.Text())	
	}

	return data,resp.StatusCode
}

func RedirectPostRequest(ctx *gin.Context, url string) int {

	var client http.Client

	// build request
	req, err := http.NewRequest(http.MethodPost, url, ctx.Request.Body)
	req.Header.Add("Content-Type", ctx.Request.Header.Get("Content-Type"))
	if err != nil {
		fmt.Println(err)
		return http.StatusInternalServerError
	}

	// make request
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK{
		return resp.StatusCode
	}

	return http.StatusOK
}

