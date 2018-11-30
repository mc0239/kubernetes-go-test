package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dghubble/sling"
	"github.com/gin-gonic/gin"
	"github.com/mc0239/kumuluzee-go-discovery/discovery"
)

// returns array of users with 200 OK code
func getUsers(c *gin.Context) {
	c.JSON(http.StatusOK, mockDB)
	return
}

// returns user object with 200 OK code if found
// and 404 NOT FOUND code if such user doesn't exists
func getUserByID(c *gin.Context) {
	sid := c.Param("id")
	id, err := strconv.ParseInt(sid, 0, 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			http.StatusBadRequest,
			fmt.Sprintf("ID conversion to integer failed with error: %s", err.Error()),
		})
		return
	}

	for _, e := range mockDB {
		if e.ID == id {
			c.JSON(http.StatusOK, e)
			return
		}
	}

	c.JSON(http.StatusNotFound, ErrorResponse{
		http.StatusNotFound,
		fmt.Sprintf("User with id %d not found.", id),
	})
	return
}

// Adds a new user with specified username to database and returns 201 CREATED with created user.
// Returns 400 BAD REQUEST if parsing failed or username is empty.
// Note that there's not check for duplicated usernames.
func createUser(c *gin.Context) {
	var body User
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			http.StatusBadRequest,
			fmt.Sprintf(err.Error()),
		})
		return
	}

	if body.Username == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			http.StatusBadRequest,
			fmt.Sprintf("No username specified"),
		})
		return
	}

	body.ID = int64(len(mockDB) + 1)
	mockDB = append(mockDB, body)

	c.JSON(http.StatusCreated, body)
	return
}

func getUserTodos(c *gin.Context) {
	sid := c.Param("id")
	id, err := strconv.ParseInt(sid, 0, 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			http.StatusBadRequest,
			fmt.Sprintf("ID conversion to integer failed with error: %s", err.Error()),
		})
		return
	}

	service, err := disc.DiscoverService(discovery.DiscoverOptions{
		Value: "todo-service",
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	todosAddress := fmt.Sprintf("%s/v1/todos/user/%d", service, id)
	var todosResponse []Todo

	_, err = sling.New().Get(todosAddress).ReceiveSuccess(&todosResponse)

	if err != nil {
		// Java service returned something other than code 2xx
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, todosResponse)
	return
}
