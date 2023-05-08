package main

import (
	"context"
	"log"

	"github.com/KyleJonesNV/go-service-notes/pkg/handlers"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	// swagger embed files
	// gin-swagger middleware
)

var (
	ginLambda *ginadapter.GinLambda
)

func init() {
	// stdout and stderr are sent to AWS CloudWatch Logs
	log.Printf("Gin cold start")
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.JSON(200, gin.H{
			"message": "healthy",
		})
	})	

	r.POST("/getAllForUser", func(c *gin.Context) {
		resp := handlers.GetAllForUser(c.Request)
		c.Header("Access-Control-Allow-Origin", "*")		     
		c.JSON(resp.StatusCode, gin.H{
			"body": resp.Body,
		})
	})	

	r.POST("/insertTopic", func(c *gin.Context) {
		resp := handlers.InsertTopic(c.Request)
		c.Header("Access-Control-Allow-Origin", "*")
		c.JSON(resp.StatusCode, gin.H{
			"body": resp.Body,
		})
	})	

	r.DELETE("/deleteTopic", func(c *gin.Context) {
		resp := handlers.DeleteTopic(c.Request)
		c.Header("Access-Control-Allow-Origin", "*")
		c.JSON(resp.StatusCode, gin.H{
			"body": resp.Body,
		})
	})	

	r.POST("/insertNote", func(c *gin.Context) {
		resp := handlers.InsertNote(c.Request)
		c.Header("Access-Control-Allow-Origin", "*")
		c.JSON(resp.StatusCode, gin.H{
			"body": resp.Body,
		})
	})

	r.POST("/getAllNotes", func(c *gin.Context) {
		resp := handlers.GetAllNotes(c.Request)
		c.Header("Access-Control-Allow-Origin", "*")
		c.JSON(resp.StatusCode, gin.H{
			"body": resp.Body,
		})
	})

	r.POST("/deleteNote", func(c *gin.Context) {
		resp := handlers.DeleteNote(c.Request)
		c.Header("Access-Control-Allow-Origin", "*")
		c.JSON(resp.StatusCode, gin.H{
			"body": resp.Body,
		})
	})	

	ginLambda = ginadapter.New(r)
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(handler)
}
