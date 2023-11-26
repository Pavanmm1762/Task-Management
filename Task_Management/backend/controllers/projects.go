// controllers.go
package controllers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go/task_management/backend/utils"
	"github.com/gocql/gocql"
)

// InitTaskRoutes initializes routes for tasks
func InitProjectRoutes(router *gin.RouterGroup) {
	router.GET("/projects", GetProjects)
	router.POST("/project", CreateProject)
	// Add routes for updateTask and deleteTask
}

// CreateProject creates a new project
func CreateProject(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	var project utils.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := getUserId(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	project.ID = gocql.TimeUUID()

	if err := utils.Session.Query("INSERT INTO projects (id, name, description, start_date, due_date, owner_id) VALUES (?, ?, ?, ?, ?, ?)",
		project.ID, project.Name, project.Description, project.StartDate, project.DueDate, userID).Exec(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, project)
}

// GetProjects gets all projects
func GetProjects(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")

	var projects []utils.Project

	userID, err := getUserId(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	iter := utils.Session.Query("SELECT id, name, description, start_date, due_date,status FROM projects where owner_id = ? ", userID).Iter()
	for {
		var project utils.Project

		if !iter.Scan(&project.ID, &project.Name, &project.Description, &project.StartDate, &project.DueDate, &project.Status) {
			break
		}

		projects = append(projects, project)
	}

	if err := iter.Close(); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, projects)
}

var jwtSecret = []byte("secure_secret_key")

func getUserId(tokenString string) (gocql.UUID, error) {

	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

	token, err := jwt.ParseWithClaims(tokenString, &utils.Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method")
		}
		return jwtSecret, nil
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
