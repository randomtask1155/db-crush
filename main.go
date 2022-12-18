package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	//socket = "/var/vcap/sys/run/pxc-mysql/mysqld.sock"
	user      = "root"
	password  string
	host      = "mysql.service.cf.internal"
	port      = "3306"
	database  = "crushdb"
	table     = "crushit"
	db        *sql.DB
	maxConn   = 100
	interval  int64
	randMax   int64
	sqlString string
)

func init() {
	interval = 1
	randMax = 60000
	rand.Seed(time.Now().Unix())

	HOST := os.Getenv("MYSQL_HOST")
	PORT := os.Getenv("MYSQL_TCP_PORT")
	DATABASE := os.Getenv("MYSQL_DATABASE")
	TABLE := os.Getenv("MYSQL_TABLE")
	USER := os.Getenv("MYSQL_USER")
	PASS := os.Getenv("MYSQL_PWD")
	INTERVAL := os.Getenv("Q_INTERVAL")
	MAXCONN := os.Getenv("Q_MAXCONN")
	RANDMAX := os.Getenv("Q_RANDMAX")
	if HOST != "" {
		host = HOST
	}
	if PORT != "" {
		port = PORT
	}
	if DATABASE != "" {
		database = DATABASE
	}
	if TABLE != "" {
		table = TABLE
	}
	if USER != "" {
		user = USER
	}
	if PASS != "" {
		password = PASS
	}
	if INTERVAL != "" {
		in, err := strconv.Atoi(INTERVAL)
		if err != nil {
			log.Fatalf("invalid interval settings:%s\n", err)
		}
		interval = int64(in)
	}
	if MAXCONN != "" {
		in, err := strconv.Atoi(MAXCONN)
		if err != nil {
			log.Fatalf("invalid interval settings:%s\n", err)
		}
		maxConn = in
	}
	if RANDMAX != "" {
		in, err := strconv.Atoi(RANDMAX)
		if err != nil {
			log.Fatalf("invalid interval settings:%s\n", err)
		}
		randMax = int64(in)
	}

	sqlString = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, database)
}

/*
| Field            | Type             | Null | Key | Default           | Extra          |
+------------------+------------------+------+-----+-------------------+----------------+
| code             | varchar(255)     | YES  | UNI | NULL              |                |
| authentication   | blob             | YES  |     | NULL              |                |
| created          | timestamp        | NO   |     | CURRENT_TIMESTAMP |                |
| expiresat        | bigint(20)       | NO   | MUL | 0                 |                |
| user_id          | varchar(36)      | YES  |     | NULL              |                |
| client_id        | varchar(255)     | YES  |     | NULL              |                |
| id               | int(11) unsigned | NO   | PRI | NULL              | auto_increment |
| identity_zone_id | varchar(36)      | YES  |     | NULL              |                |
+------------------+------------------+------+-----+-------------------+----------------+
*/
func initDB() {
	var err error
	//fmt.Println(sqlString)
	initdb, err := sql.Open("mysql", sqlString)
	if err != nil {
		log.Panicln(err)
	}
	defer initdb.Close()

	// log.Println("creating resources")
	// _, err = initdb.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", database))
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// _, err = initdb.Exec(fmt.Sprintf("use %s;", database))
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// _, err = initdb.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (code varchar(255), authentication BLOB, created TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL, expiresat BIGINT DEFAULT 0 NOT NULL, user_id VARCHAR(36) NULL, client_id VARCHAR(36) NULL, `id` int(11) unsigned PRIMARY KEY AUTO_INCREMENT);", table))
	// if err != nil {
	// 	log.Fatalln(err)
	// }

}

func getFurtureTime() string {
	t := time.Now().UnixMilli() + rand.Int63n(randMax)
	return strconv.Itoa(int(t))
}

func getCurrentTime() string {
	t := time.Now().UnixMilli()
	return strconv.Itoa(int(t))
}

func insertQuery(ch chan bool) {
	current := getCurrentTime()
	future := getFurtureTime()
	_, err := db.Exec(fmt.Sprintf("insert into %s values (? , '1234', DEFAULT, ?, 'user1', 'client1', DEFAULT, 'test' ); ", table), current, future)
	if err != nil {
		log.Println(err)
	}
	ch <- true
}

func deleteQuery(ch chan bool) {
	current := getCurrentTime()
	_, err := db.Exec(fmt.Sprintf("delete from %s where expiresat > 0 AND expiresat < ? ;", table), current)
	if err != nil {
		log.Println(err)
	}
	ch <- true
}

func dbStatsOutput() {
	for {
		time.Sleep(60 * time.Second)
		log.Printf("%v\n", db.Stats())
	}
}

func main() {
	log.Println("initdb starting")
	initDB()
	log.Println("initdb completed")
	var err error
	db, err = sql.Open("mysql", sqlString)
	if err != nil {
		log.Panicln(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(maxConn)

	log.Println("starting workload")
	go dbStatsOutput()
	ch := make(chan bool, maxConn*2)
	for {
		select {
		case <-ch:
			r := rand.Intn(2)
			if r > 0 {
				go insertQuery(ch)
			} else {
				go deleteQuery(ch)
			}
		default:
			ch <- true // keep buffer full
		}
		time.Sleep(time.Duration(interval) * time.Millisecond)
	}
	//
	//log.Println(db.Stats)
}
