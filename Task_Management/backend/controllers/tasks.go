// tasks.go
package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go/task_management/backend/utils"
	"github.com/gocql/gocql"
)

// InitTaskRoutes initializes routes for tasks
func InitTaskRoutes(router *gin.RouterGroup) {
	router.GET("/tasks", getTasks)
	router.POST("/task", createTask)
	// Add routes for updateTask and deleteTask
}

func getTasks(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	var tasks []utils.Task

	userID, err := getCurrentUserId(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	iter := utils.Session.Query(`
	SELECT id, name
	FROM projects where owner_id = ?`,
		userID).Iter()

	for {
		var projectID gocql.UUID
		var projectName string
		if !iter.Scan(&projectID, &projectName) {
			if err := iter.Close(); err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			break
		}

		fmt.Printf("Project ID: %v, Project Name: %v\n", projectID, projectName)

		// Fetch tasks for each project
		taskIter := utils.Session.Query(`
			SELECT task_id, task_name, progress, status
			FROM tasks
			WHERE project_id = ?
			`, projectID).Iter()

		var taskID gocql.UUID
		var taskTitle, taskStatus string
		var taskProgress int

		for taskIter.Scan(&taskID, &taskTitle, &taskProgress, &taskStatus) {
			var task utils.Task

			task.ProjectName = projectName
			task.Title = taskTitle
			task.Progress = taskProgress
			task.Status = taskStatus

			tasks = append(tasks, task)

			fmt.Printf("Project name : %v,Task ID: %v, Task Title: %v, Progress: %v, Status: %v\n", projectName, taskID, taskTitle, taskProgress, taskStatus)
		}

		// Close the task iterator
		if err := taskIter.Close(); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}
	// Print the response before sending it
	fmt.Println("Response:", tasks)
	c.JSON(200, tasks)

	// iter := utils.Session.Query("SELECT task_id, task_name, progress, status FROM tasks").Iter()
	// for {
	// 	var task utils.Task

	// 	if !iter.Scan(&task.ID, &task.Title, &task.Progress, &task.status) {
	// 		break
	// 	}

	// 	tasks = append(tasks, task)
	// }

	// if err := iter.Close(); err != nil {
	// 	c.JSON(500, gin.H{"error": err.Error()})
	// 	return
	// }
	// // Print the response before sending it
	// fmt.Println("Response:", tasks)
	// c.JSON(200, tasks)
}

func createTask(c *gin.Context) {
	var task utils.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	task.ID = gocql.TimeUUID()

	if err := utils.Session.Query("INSERT INTO tasks (id, title,progress) VALUES (?, ?, ?)", task.ID, task.Title, task.Progress).Exec(); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, task)
}

var jwtSecretK = []byte("secure_secret_key")

func getCurrentUserId(tokenString string) (gocql.UUID, error) {

	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

	token, err := jwt.ParseWithClaims(tokenString, &utils.Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method")
		}
		return jwtSecretK, nil
	})

	// Check for errors
	if err != nil {
		return gocql.UUID{}, err
	}

	// Check if the token is valid
	if claims, ok := token.Claims.(*utils.Claims); ok && token.Valid {
		return claims.UserID, nil
	}

	return gocql.UUID{}, errors.New("Invalid token")
}
