package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func GetEnvWithDefault(name, defaultValue string) string {
	res, exists := os.LookupEnv(name)
	if !exists {
		fmt.Sprintf("env lookup for value \"%s\" failed. Using default value: \"%s\"", name, defaultValue)
		return defaultValue
	}
	return res
}

func ConnectPostgres(host string, dbport string, user string, password string, dbname string) (*sql.DB, error) {
	intport, err := strconv.Atoi(dbport)
	postgresqlDbInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, intport, user, password, dbname)

	db, err := sql.Open("postgres", postgresqlDbInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db, err
}

func PostgresQuery(db *sql.DB, username string, userpassword string) (string, error) {
	var token string
	defer db.Close()
	row := db.QueryRow("SELECT token FROM person where username=$1 and password=$2", username, userpassword)
	if err := row.Scan(&token); err != nil {
		fmt.Println("error", err)
		return "", err
	} else {
		return token, err
	}

}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {

	router := gin.Default()
	router.POST("/login", func(c *gin.Context) {
		var json User
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusUnauthorized, "error in receiving user info and transfer to json")
		}

		host := GetEnvWithDefault("DB_HOST", "localhost")
		dbport := GetEnvWithDefault("DB_PORT", "5432")
		user := GetEnvWithDefault("POSTGRES_USER", "switchboard_admin")
		password := GetEnvWithDefault("POSTGRES_PASSWORD", "12345687")
		dbname := GetEnvWithDefault("DB_NAME", "go_test")

		db, err := ConnectPostgres(host, dbport, user, password, dbname)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "ConnectPostgres error")
		} // Question: how to put c.JSON to the outside of POST

		token, err := PostgresQuery(db, json.Username, json.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "PostgresQuery error")
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status": "ok",
				"token":  token,
			})
		}
	})
	router.Run(":" + GetEnvWithDefault("PORT", "8081"))
}
