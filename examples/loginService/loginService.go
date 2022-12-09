package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Token struct {
	Token string `json:"token"`
}

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

func QueryLogin(db *sql.DB, username string, userpassword string) (string, error) {
	var token string
	row := db.QueryRow("SELECT token FROM person where username=$1 and password=$2", username, userpassword)
	if err := row.Scan(&token); err != nil {
		fmt.Println("error", err)
		return "", err
	} else {
		return token, err
	}
}

func QueryValidate(db *sql.DB, token string) error {
	row := db.QueryRow("SELECT token FROM person where token=$1", token)
	if err := row.Scan(&token); err != nil {
		fmt.Println("error validating token", err)
		return err
	} else {
		return err
	}
}

func main() {

	router := gin.Default()

	// Get env values
	host := GetEnvWithDefault("DB_HOST", "localhost")
	dbport := GetEnvWithDefault("DB_PORT", "5432")
	user := GetEnvWithDefault("POSTGRES_USER", "switchboard_admin")
	password := GetEnvWithDefault("POSTGRES_PASSWORD", "12345687")
	dbname := GetEnvWithDefault("DB_NAME", "go_test")
	port := GetEnvWithDefault("PORT", "8081")

	// Connect to Postgres
	time.Sleep(10 * time.Second)
	db, err := ConnectPostgres(host, dbport, user, password, dbname)
	if err != nil {
		fmt.Errorf("Postgres connection failed", err)
	}
	defer db.Close()

	router.POST("/login", func(c *gin.Context) {
		var json User
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusInternalServerError, "error in receiving user info and transfer to json")
		}

		token, err := QueryLogin(db, json.Username, json.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "Invalid username/password")
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status": "ok",
				"token":  token,
			})
		}
	})

	router.POST("/validate", func(c *gin.Context) {
		var json Token
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusInternalServerError, "error in receiving token info and transfer to json")
		}

		err = QueryValidate(db, json.Token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "Invalid token")
		} else {
			c.JSON(http.StatusOK, "User is authorized")
		}
	})

	router.Run(":" + port)
}
