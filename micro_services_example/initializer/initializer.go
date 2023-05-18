package initializer

import (
	"crypto/rand"
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
)

var DB *sql.DB
var RDB *redis.Client
var JwtKey = []byte("shh!it's_secret_key")

func Initialize() {
	initJwtKey()
	newDB()
	newRedis()
}

func initJwtKey() {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}
	JwtKey = bytes
}

func newDB() {
	var err error
	DB, err = sql.Open("mysql", "micro_services_example:micro_services_example@tcp(127.0.0.1:3306)/micro_services_example")
	if err != nil {
		log.Fatal(err)
	}
}

func newRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}
