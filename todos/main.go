package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/mc0239/kumuluzee-go-config/config"
	"github.com/mc0239/kumuluzee-go-discovery/discovery"

	"github.com/gin-gonic/gin"
)

var mockDB []Todo
var conf config.Util
var disc discovery.Util

func main() {
	// initialize functions
	initDB()
	initConfig()
	initDiscovery()

	// register service to service registry
	disc.RegisterService(discovery.RegisterOptions{})

	router := gin.Default()

	// Registers middleware function, which for each request checks our external configuration and
	// if 'maintenance' key is set to true, it will return error saying service is unavailable,
	// otherwise it will call next handler.
	// To test, while running go to http://localhost:8500 and change key
	// 'environments/dev/services/node-service/1.0.0/config/rest-config/maintenance' to 'true' and
	// then try to perform a request. To enable it again, just change the key to 'false'
	router.Use(func(c *gin.Context) {
		maintenanceMode, _ := conf.GetBool("rest-config.maintenance")
		if maintenanceMode {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, ErrorResponse{
				http.StatusServiceUnavailable,
				"Service is undergoing maintenance, check back in a minute.",
			})
		} else {
			c.Next()
		}
	})

	// prepare routes and map them to handlers
	v1c := router.Group("/v1/todos")
	{
		// GET /v1/todos
		v1c.GET("/", getTodos)
		// GET /v1/todos/:userId
		v1c.GET("/user/:userId", getUserTodos)
		// GET /v1/todos/:id
		v1c.GET("/todo/:id", getTodoByID)
		// POST /v1/todos
		v1c.POST("/", createTodo)
	}

	// run REST API server
	port, ok := conf.GetInt("kumuluzee.server.http.port")
	if !ok {
		port = 9000
	}
	router.Run(fmt.Sprintf(":%d", port))
}

func initDB() {
	mockDB = make([]Todo, 0)
	mockDB = append(mockDB,
		Todo{ID: 0, UserID: 100, Title: "First TODO", Content: "Thank you todos, very cool!", DateCreated: time.Now().Format(time.RFC1123)})
}

func initConfig() {
	conf = config.NewUtil(config.Options{
		Extension:  "consul",
		ConfigPath: "config.yaml",
	})
}

func initDiscovery() {
	disc = discovery.New(discovery.Options{
		Extension:  "consul",
		ConfigPath: "config.yaml",
	})
}
