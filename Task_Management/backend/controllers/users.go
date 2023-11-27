// controllers/users.go
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
func InitUserRoutes(router *gin.RouterGroup) {
	router.GET("/users-list", getUsers)
	router.POST("/add-user", addUser)
	// Add routes for updateTask and deleteTask
}

// CreateProject creates a new project
func addUser(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	var user utils.Users
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	adminId, err := getAdminId(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	user.UserId = gocql.TimeUUID()

	if err := utils.Session.Query("INSERT INTO users (user_id, first_name, last_name, user_role, user_email, user_password, admin_id) VALUES (?, ?, ?, ?, ?, ?, ?)",
		user.UserId, user.FirstName, user.LastName, user.UserRole, user.UserEmail, user.UserPassword, adminId).Exec(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// GetProjects gets all projects
func getUsers(c *gin.Context) {
	var users []utils.Users
	tokenString1 := c.GetHeader("Authorization")

	admin_id, err := getAdminId(tokenString1)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	iter := utils.Session.Query("SELECT user_id, first_name, last_name, user_role, user_email  FROM users where admin_id = ?", admin_id).Iter()
	for {
		var user utils.Users

		if !iter.Scan(&user.UserId, &user.FirstName, &user.LastName, &user.UserRole, &user.UserEmail) {
			break
		}

		users = append(users, user)
	}

	if err := iter.Close(); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

var jwtSecretKey = []byte("secure_secret_key")

func getAdminId(tokenString string) (gocql.UUID, error) {

	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

	token, err := jwt.ParseWithClaims(tokenString, &utils.Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method")
		}
		return jwtSecretKey, nil
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
