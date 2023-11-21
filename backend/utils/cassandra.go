// utils/cassandra.go
package utils

import (
	"log"

	"github.com/gocql/gocql"
)

var session *gocql.Session

func InitCassandra() {
	var err error
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "example_keyspace"

	session, err = cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to Cassandra successfully")
}
