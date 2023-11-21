// main.go
package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
)

var session *gocql.Session

type Task struct {
	ID    gocql.UUID `json:"id"`
	Title string     `json:"title"`
}

func main() {
	var err error
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "example_keyspace"
	session, err = cluster.CreateSession()
	if err != nil {
		fmt.Println("Error creating session:", err)
		return
	}
	defer session.Close()

	router := gin.Default()
	router.GET("/tasks", getTasks)

	router.Run(":8080")
}

func getTasks(c *gin.Context) {
	var tasks []Task
	iter := session.Query("SELECT id, title FROM tasks").Iter()
	var task Task
	for iter.Scan(&task.ID, &task.Title) {
		tasks = append(tasks, task)
	}
	c.JSON(http.StatusOK, tasks)
}
