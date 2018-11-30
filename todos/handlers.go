package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// returns array of todos with 200 OK code
func getTodos(c *gin.Context) {
	c.JSON(http.StatusOK, mockDB)
	return
}

func getUserTodos(c *gin.Context) {
	suid := c.Param("userId")
	uid, err := strconv.ParseInt(suid, 0, 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			http.StatusBadRequest,
			fmt.Sprintf("ID conversion to integer failed with error: %s", err.Error()),
		})
		return
	}

	var userTodos []Todo
	userTodos = make([]Todo, 0)

	for _, e := range mockDB {
		if e.UserID == uid {
			userTodos = append(userTodos, e)
		}
	}

	c.JSON(http.StatusOK, userTodos)
	return
}

func getTodoByID(c *gin.Context) {
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
		fmt.Sprintf("Todo with id %d not found.", id),
	})
	return
}

func createTodo(c *gin.Context) {
	var body Todo
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			http.StatusBadRequest,
			fmt.Sprintf(err.Error()),
		})
		return
	}

	if body.UserID == 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			http.StatusBadRequest,
			fmt.Sprintf("No userID specified"),
		})
		return
	}

	body.ID = int64(len(mockDB) + 1)
	body.DateCreated = time.Now().Format(time.RFC1123)
	mockDB = append(mockDB, body)

	c.JSON(http.StatusCreated, body)
	return
}
