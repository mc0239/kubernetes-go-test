package main

import (
	"fmt"
	"net/http"

	"github.com/mc0239/kumuluzee-go-config/config"
	"github.com/mc0239/kumuluzee-go-discovery/discovery"

	"github.com/gin-gonic/gin"
)

var mockDB []User
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
	v1c := router.Group("/v1/users")
	{
		// GET /v1/users
		v1c.GET("/", getUsers)
		// GET /v1/users/:id
		v1c.GET("/:id", getUserByID)
		// GET /v1/users/:id/todos
		v1c.GET("/:id/todos", getUserTodos)
		// POST /v1/users
		v1c.POST("/", createUser)
	}

	// run REST API server
	port, ok := conf.GetInt("kumuluzee.server.http.port")
	if !ok {
		port = 9000
	}
	router.Run(fmt.Sprintf(":%d", port))
}

func initDB() {
	mockDB = make([]User, 0)
	mockDB = append(mockDB,
		User{100, "john"},
		User{101, "ann"},
		User{102, "elizabeth"},
		User{103, "isaac"},
		User{104, "barret"},
		User{105, "terry"},
	)
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
