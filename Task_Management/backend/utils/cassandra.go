package utils

// utils/cassandra.go

import (
	"log"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
)

var Session *gocql.Session

func InitCassandra() {
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "task_management_team_collaboration"
	var err error
	Session, err = cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	log.Println("cassandra successfully connected...")
}

func GenerateUUID() string {
	return uuid.New().String()
}
