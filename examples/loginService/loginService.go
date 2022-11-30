package main
 
import (
	"database/sql"
	"fmt"
	"net/http"
	_ "github.com/lib/pq"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"os"
	"strconv"
  )


func InitializeEnv() error {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return err
	}
	return nil
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
	if err := row.Scan(&token); err != nil{
		fmt.Println("error", err)
		return "", err
	}else{
		return token, err
	}
	
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}


func main(){
	
	InitializeEnv()
	router := gin.Default()
	router.POST("/loginJSON", func(c *gin.Context) { //Question: change into /login?
		var json User
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusUnauthorized, "error in receiving user info and transfer to json")
		}

		// fmt.Println("get username as: ", json.Username)

		// fmt.Println("trying to make database connection")
		// fmt.Println(os.Getenv("host"))
		// fmt.Println(os.Getenv("dbport"))
		// fmt.Println(os.Getenv("user"))
		// fmt.Println(os.Getenv("password"))
		// fmt.Println(os.Getenv("dbname"))
		db, err := ConnectPostgres(os.Getenv("host"), os.Getenv("dbport"), os.Getenv("user"), os.Getenv("password"), os.Getenv("dbname"))
		if err != nil{
			c.JSON(http.StatusUnauthorized, "ConnectPostgres error")
		} // Question: how to put c.JSON to the outside of POST

		token, err := PostgresQuery(db, json.Username, json.Password)
		if err != nil{
			c.JSON(http.StatusUnauthorized, "PostgresQuery error")
		} else{
			c.JSON(http.StatusOK, gin.H{
				"status": "ok",
				"token": token,
			})
		}
	})
	router.Run("localhost:"+os.Getenv("localhostserviceport"))
}