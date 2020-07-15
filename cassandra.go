package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"github.com/gocql/gocql"
)

func main() {

	var ACCOUNTNAME = ""
	var PASSWORD = ""
	var CONTACTPOINT = ""

	// connect to the cluster
	fmt.Println("Creating Azure Cosmos DB Cassandra API Connection...")
	cluster := gocql.NewCluster(CONTACTPOINT)
	cluster.Port = 10350
	var sslOptions = new(gocql.SslOptions)
	cluster.Timeout = 30 * time.Second
	cluster.ConnectTimeout = 30 * time.Second
	cluster.DisableInitialHostLookup = true

	// If you want to enable client-side SSL server cert verification do this:
	sslOptions.EnableHostVerification = true
	sslOptions.CaPath = "<path/to/cert.cer>"
	sslOptions.Config = &tls.Config{}
	sslOptions.ServerName = CONTACTPOINT
	cluster.SslOpts = sslOptions

	// If you do NOT want to enable client-side SSL, set sslOptions.EnableHostVerification to false and ignore the other options.
	// sslOptions.EnableHostVerification = false
	// cluster.SslOpts = sslOptions

	cluster.ProtoVersion = 4
	cluster.Authenticator = gocql.PasswordAuthenticator{Username: ACCOUNTNAME, Password: PASSWORD}
	session, _ := cluster.CreateSession()
	defer session.Close()

	fmt.Println("Drop keyspace if exists...")
	if err := session.Query(`DROP KEYSPACE IF EXISTS uprofile;`).Exec(); err != nil {
		log.Fatal(err)
	}
	time.Sleep(5 * time.Second)
	fmt.Println("Creating keyspace...")
	if err := session.Query(`CREATE KEYSPACE uprofile WITH replication = {\'class\': \'NetworkTopologyStrategy\', \'datacenter1\' : \'1\' };`).Exec(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Finished creating keyspace. Creating table...")
	time.Sleep(8 * time.Second)

	if err := session.Query(`CREATE TABLE uprofile.user (user_id int PRIMARY KEY, user_name text, user_bcity text)`).Exec(); err != nil {
		log.Fatal(err)
	}
	time.Sleep(8 * time.Second)
	fmt.Println("Finished creating table. Inserting rows...")
	queries := []string{
		"`INSERT INTO  uprofile.user (user_id, user_name , user_bcity) VALUES (1, 'AdrianaS', 'Seattle');`",
		"`INSERT INTO  uprofile.user (user_id, user_name , user_bcity) VALUES (2, 'JiriK', 'Toronto');`",
		"`INSERT INTO  uprofile.user (user_id, user_name , user_bcity) VALUES (3, 'IvanH', 'Mumbai');`",
		"`INSERT INTO  uprofile.user (user_id, user_name , user_bcity) VALUES (4, 'IvanH', 'Seattle');`",
		"`INSERT INTO  uprofile.user (user_id, user_name , user_bcity) VALUES (5, 'IvanaV', 'Belgaum');`",
		"`INSERT INTO  uprofile.user (user_id, user_name , user_bcity) VALUES (6, 'LiliyaB', 'Seattle');`",
		"`INSERT INTO  uprofile.user (user_id, user_name , user_bcity) VALUES (7, 'JindrichH', 'Buenos Aires');`",
		"`INSERT INTO  uprofile.user (user_id, user_name , user_bcity) VALUES (8, 'AdrianaS', 'Seattle');`",
		"`INSERT INTO  uprofile.user (user_id, user_name , user_bcity) VALUES (9, 'JozefM', 'Seattle');`"}
	for i, query := range queries {
		if err := session.Query(query).Exec(); err != nil {
			fmt.Println(i, query)
			log.Fatal(err)
		}
	}

	fmt.Println("Finished inserting rows. Select from table...")
	time.Sleep(5 * time.Second)

	var user_id int
	var user_name string
	var user_bcity string

	fmt.Println("user_id 	| user_name 	| user_bcity")
	fmt.Println("--------------------------------------------")
	iter := session.Query(`SELECT user_id, user_name, user_bcity FROM uprofile.user;`).Iter()
	for iter.Scan(&user_id, &user_name, &user_bcity) {
		fmt.Println(user_id, "		|", user_name, "	|", user_bcity)
	}

}
