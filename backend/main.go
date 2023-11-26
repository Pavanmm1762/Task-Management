// main.go
package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"golang.org/x/crypto/bcrypt"
)

var session *gocql.Session

func init() {
	var err error
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "example_keyspace"

	session, err = cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	log.Println("Connected to Cassandra successfully")

	// Creating "users" table if not exists
	if err := session.Query(`
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY,
			username TEXT,
			password TEXT
		)`).Exec(); err != nil {
		log.Fatal(err)
	}
}

type User struct {
	ID       gocql.UUID `json:"id"`
	Username string     `json:"username"`
	Password string     `json:"password"`
}

func main() {
	r := gin.Default()

	// Apply CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"} // Adjust the origin to match your React app's URL
	r.Use(cors.New(config))

	r.POST("/api/signup", func(c *gin.Context) {
		var user User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		// Hash password before storing in Cassandra
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		user.ID = gocql.TimeUUID()
		user.Password = string(hashedPassword)

		if session.Closed() {
			log.Println("Cassandra session is closed. Reconnecting...")
			// Reconnect or handle closed session (example: create a new session)
			var err error
			cluster := gocql.NewCluster("127.0.0.1")
			cluster.Keyspace = "example_keyspace"

			session, err = cluster.CreateSession()
			if err != nil {
				log.Fatal(err)
			}
			defer session.Close()
			if err != nil {
				log.Println("Error reconnecting to Cassandra:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
				return
			}
			defer session.Close()
		}

		if err := session.Query(`
			INSERT INTO users (id, username, password) VALUES (?, ?, ?)`,
			user.ID, user.Username, user.Password).Exec(); err != nil {
			log.Println("Error inserting user:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
	})

	r.POST("/api/login", func(c *gin.Context) {
		var loginData User
		if err := c.BindJSON(&loginData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		var storedPassword string
		var storedID gocql.UUID
		if session.Closed() {
			log.Println("Cassandra session is closed. Reconnecting...")
			// Reconnect or handle closed session (example: create a new session)
			var err error
			cluster := gocql.NewCluster("127.0.0.1")
			cluster.Keyspace = "example_keyspace"

			session, err = cluster.CreateSession()
			if err != nil {
				log.Fatal(err)
			}
			defer session.Close()
		}
		if err := session.Query(`
			SELECT id, password FROM users WHERE username = ? LIMIT 1 ALLOW FILTERING;`,
			loginData.Username).Consistency(gocql.One).Scan(&storedID, &storedPassword); err != nil {
			log.Println("Error during login:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(loginData.Password))
		if err != nil {
			log.Println("Error during compare:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
	})

	r.Run(":8080")
}
